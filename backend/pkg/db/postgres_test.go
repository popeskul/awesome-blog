package db_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB(t *testing.T) {
	logger := logrus.New()
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	// Создаем мок SQL с опцией мониторинга пингов
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Ожидаем вызов PingContext
	mock.ExpectPing()

	// Заменяем реальную функцию sql.Open на нашу тестовую версию
	oldSqlOpen := db.SqlOpen
	defer func() { db.SqlOpen = oldSqlOpen }()
	db.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}

	// Вызываем тестируемую функцию
	dbForMock, err := db.NewPostgresDB(cfg, logger)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, dbForMock)

	// Проверяем, что все ожидания мока были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDB_Close(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	logger := logrus.New()
	pgDB := &db.PostgresDB{DB: mockDB, Logger: logger}

	mock.ExpectClose()

	err = pgDB.Close()
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDB_HealthCheck(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	logger := logrus.New()
	pgDB := &db.PostgresDB{DB: mockDB, Logger: logger}

	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	err = pgDB.HealthCheck(context.Background())
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDB_BeginTx(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	logger := logrus.New()
	pgDB := &db.PostgresDB{DB: mockDB, Logger: logger}

	mock.ExpectBegin()

	tx, err := pgDB.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// sqlOpen is a variable that holds the original sql.Open function
type sqlOpenFunc func(driverName, dataSourceName string) (*sql.DB, error)

// sqlOpen is a variable that holds the original sql.Open function
var sqlOpen sqlOpenFunc = sql.Open

func init() {
	sqlOpen = sql.Open
}
