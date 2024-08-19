package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/sirupsen/logrus"
)

// SqlOpen is a variable that holds the sql.Open function
var SqlOpen = sql.Open

type PostgresDB struct {
	*sql.DB
	Logger *logrus.Logger
}

func NewPostgresDB(cfg config.DatabaseConfig, logger *logrus.Logger) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := SqlOpen("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Не используем Exec для установки параметров
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	logger.Info("Successfully connected to the database")

	return &PostgresDB{DB: db, Logger: logger}, nil
}

func (pdb *PostgresDB) Close() error {
	if err := pdb.DB.Close(); err != nil {
		pdb.Logger.WithError(err).Error("Failed to close database connection")
		return fmt.Errorf("error closing database connection: %w", err)
	}

	pdb.Logger.Info("Database connection closed successfully")

	return nil
}

func (pdb *PostgresDB) HealthCheck(ctx context.Context) error {
	var result int
	if err := pdb.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		pdb.Logger.WithError(err).Error("Database health check failed")
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (pdb *PostgresDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := pdb.DB.BeginTx(ctx, nil)
	if err != nil {
		pdb.Logger.WithError(err).Error("Failed to start transaction")
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return tx, nil
}
