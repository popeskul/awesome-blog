package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/infrastructure/database/postgres"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

func TestUserRepository_CreateUser_Success(t *testing.T) {
	tests := []struct {
		name         string
		user         *entity.NewUser
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedUser *entity.User
	}{
		{
			name: "Create user successfully",
			user: &entity.NewUser{
				Username:     "testuser",
				Email:        "testuser@example.com",
				PasswordHash: "securepassword",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sqlmock.AnyArg(), "testuser", "testuser@example.com", "securepassword").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
						AddRow(uuid.New(), "testuser", "testuser@example.com", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Username: "testuser",
				Email:    "testuser@example.com",
			},
		},
		{
			name: "Create user with different data",
			user: &entity.NewUser{
				Username:     "anotheruser",
				Email:        "anotheruser@example.com",
				PasswordHash: "anotherpassword",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sqlmock.AnyArg(), "anotheruser", "anotheruser@example.com", "anotherpassword").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
						AddRow(uuid.New(), "anotheruser", "anotheruser@example.com", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Username: "anotheruser",
				Email:    "anotheruser@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.CreateUser(context.Background(), tt.user)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedUser.Username, user.Username)
			assert.Equal(t, tt.expectedUser.Email, user.Email)
			assert.NotEqual(t, uuid.Nil, user.Id)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_CreateUser_Fail(t *testing.T) {
	tests := []struct {
		name        string
		user        *entity.NewUser
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "Failed to create user - SQL error",
			user:        &entity.NewUser{Username: "testuser", Email: "testuser@example.com", PasswordHash: "securepassword"},
			mockError:   errors.New("failed to create user"),
			expectedErr: "failed to create user",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sqlmock.AnyArg(), "testuser", "testuser@example.com", "securepassword").
					WillReturnError(errors.New("failed to create user"))
			},
		},
		{
			name:        "Failed to create user - Duplicate entry",
			user:        &entity.NewUser{Username: "duplicateuser", Email: "duplicateuser@example.com", PasswordHash: "password"},
			mockError:   sql.ErrNoRows,
			expectedErr: "failed to create user",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sqlmock.AnyArg(), "duplicateuser", "duplicateuser@example.com", "password").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.CreateUser(context.Background(), tt.user)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
			assert.Nil(t, user)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserById_Success(t *testing.T) {
	tests := []struct {
		name         string
		id           uuid.UUID
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedUser *entity.User
	}{
		{
			name: "Get user by ID successfully",
			id:   userId1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userId1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
						AddRow(userId1, "testuser", "testuser@example.com", "securepassword", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Id:           userId1,
				Username:     "testuser",
				Email:        "testuser@example.com",
				PasswordHash: "securepassword",
			},
		},
		{
			name: "Get user by different ID",
			id:   userId2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userId2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
						AddRow(userId2, "anotheruser", "anotheruser@example.com", "anotherpassword", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Id:           userId2,
				Username:     "anotheruser",
				Email:        "anotheruser@example.com",
				PasswordHash: "anotherpassword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.GetUserById(context.Background(), tt.id)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedUser.Id, user.Id)
			assert.Equal(t, tt.expectedUser.Username, user.Username)
			assert.Equal(t, tt.expectedUser.Email, user.Email)
			assert.Equal(t, tt.expectedUser.PasswordHash, user.PasswordHash)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserById_Fail(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "Failed to get user by ID - SQL error",
			id:          userId1,
			mockError:   errors.New("failed to get user"),
			expectedErr: "failed to get user",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userId1).
					WillReturnError(errors.New("failed to get user"))
			},
		},
		{
			name:        "Failed to get user by ID - No rows found",
			id:          userId2,
			mockError:   sql.ErrNoRows,
			expectedErr: "user not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userId2).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.GetUserById(context.Background(), tt.id)

			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserByUsername_Success(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedUser *entity.User
	}{
		{
			name:     "Get user by username successfully",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = \$1`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
						AddRow(userId1, "testuser", "testuser@example.com", "securepassword", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Id:           userId1,
				Username:     "testuser",
				Email:        "testuser@example.com",
				PasswordHash: "securepassword",
			},
		},
		{
			name:     "Get user by another username",
			username: "anotheruser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = \$1`).
					WithArgs("anotheruser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
						AddRow(userId2, "anotheruser", "anotheruser@example.com", "anotherpassword", time.Now(), time.Now()))
			},
			expectedUser: &entity.User{
				Id:           userId2,
				Username:     "anotheruser",
				Email:        "anotheruser@example.com",
				PasswordHash: "anotherpassword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.GetUserByUsername(context.Background(), tt.username)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedUser.Id, user.Id)
			assert.Equal(t, tt.expectedUser.Username, user.Username)
			assert.Equal(t, tt.expectedUser.Email, user.Email)
			assert.Equal(t, tt.expectedUser.PasswordHash, user.PasswordHash)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserByUsername_Fail(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "Failed to get user by username - SQL error",
			username:    "testuser",
			mockError:   errors.New("failed to get user"),
			expectedErr: "failed to get user",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = \$1`).
					WithArgs("testuser").
					WillReturnError(errors.New("failed to get user"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			user, err := repo.GetUserByUsername(context.Background(), tt.username)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
			assert.Nil(t, user)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetAllUsers_Fail(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "Failed to get all users - SQL error",
			expectedErr: "failed to get all users",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, username, email, created_at, updated_at FROM users`).
					WithArgs(10, 0).
					WillReturnError(errors.New("failed to get users"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			pagination := &entity.Pagination{Limit: 10, Offset: 0}
			users, err := repo.GetAllUsers(context.Background(), pagination)

			if tt.expectedErr == "failed to scan user" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, users)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		user      *entity.UpdateUser
	}{
		{
			name: "Successfully update user with all fields",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET").
					WithArgs("newUsername", "newEmail@example.com", "password", userId1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user: &entity.UpdateUser{
				Id:       userId1,
				Username: "newUsername",
				Email:    "newEmail@example.com",
				Password: "password",
			},
		},
		{
			name: "Successfully update user with only password",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET").
					WithArgs("newPassword", userId2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user: &entity.UpdateUser{
				Id:       userId2,
				Password: "newPassword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			err = repo.UpdateUser(context.Background(), tt.user)

			assert.NoError(t, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UpdateUser_Fail(t *testing.T) {
	tests := []struct {
		name        string
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
		user        *entity.UpdateUser
	}{
		{
			name:        "Failed to update user - No fields to update",
			expectedErr: "no fields to update",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			user: &entity.UpdateUser{
				Id:       userId1,
				Username: "",
				Email:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			if tt.mockError != nil {
				tt.mockSetup(mock)
			}

			err = repo.UpdateUser(context.Background(), tt.user)

			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_DeleteUserById_Success(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		id        uuid.UUID
	}{
		{
			name: "Successfully delete user by ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
					WithArgs(userId1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: userId1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			err = repo.DeleteUserById(context.Background(), tt.id)

			assert.NoError(t, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_DeleteUserById_Fail(t *testing.T) {
	tests := []struct {
		name        string
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
		id          uuid.UUID
	}{
		{
			name:        "Failed to delete user by ID - SQL error",
			mockError:   errors.New("failed to delete user"),
			expectedErr: "failed to delete user",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
					WithArgs(userId1).
					WillReturnError(errors.New("failed to delete user"))
			},
			id: userId1,
		},
		{
			name:        "Failed to delete user by ID - No rows affected",
			mockError:   errors.New("no rows affected"),
			expectedErr: "no rows affected",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
					WithArgs(userId2).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			id: userId2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			err = repo.DeleteUserById(context.Background(), tt.id)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetTotalUsers_Success(t *testing.T) {
	tests := []struct {
		name          string
		expectedTotal int
		mockSetup     func(mock sqlmock.Sqlmock)
	}{
		{
			name:          "Successfully get total users",
			expectedTotal: 42,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewUserRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			total, err := repo.GetTotalUsers(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotal, total)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
