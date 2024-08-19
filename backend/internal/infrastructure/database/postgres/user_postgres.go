package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

type UserRepository struct {
	db     *db.PostgresDB
	logger *logrus.Logger
}

func NewUserRepository(db *db.PostgresDB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.NewUser) (*entity.User, error) {
	query := `
        INSERT INTO users (id, username, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, username, email, created_at, updated_at
    `

	id := uuid.New()
	var createdUser entity.User
	err := r.db.QueryRowContext(ctx, query, id, user.Username, user.Email, user.PasswordHash).Scan(
		&createdUser.Id, &createdUser.Username, &createdUser.Email, &createdUser.CreatedAt, &createdUser.UpdatedAt,
	)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &createdUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
        SELECT id, username, email, password, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
        SELECT id, username, email, password, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, params *entity.Pagination) ([]*entity.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users`

	if params.Sort != "" {
		sortField, sortOrder := parseSortParam(params.Sort)
		allowedSortFields := map[string]bool{
			"created_at": true,
			"username":   true,
			"email":      true,
		}
		if !allowedSortFields[sortField] {
			return nil, fmt.Errorf("invalid sort field: %s", sortField)
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " LIMIT $1 OFFSET $2"

	rows, err := r.db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			r.logger.WithError(err).Error("Failed to scan user")
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error occurred during row iteration")
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.UpdateUser) error {
	if user.Username == "" && user.Email == "" && user.Password == "" {
		return errors.New("no fields to update")
	}

	query := `UPDATE users SET updated_at = NOW()`
	var args []interface{}
	argIndex := 1

	if user.Username != "" {
		query += fmt.Sprintf(", username = $%d", argIndex)
		args = append(args, user.Username)
		argIndex++
	}
	if user.Email != "" {
		query += fmt.Sprintf(", email = $%d", argIndex)
		args = append(args, user.Email)
		argIndex++
	}
	if user.Password != "" {
		query += fmt.Sprintf(", password = $%d", argIndex)
		args = append(args, user.Password)
		argIndex++
	}

	query += fmt.Sprintf(` WHERE id = $%d`, argIndex)
	args = append(args, user.Id)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *UserRepository) GetTotalUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get total users count")
		return 0, fmt.Errorf("failed to get total users count: %w", err)
	}

	return total, nil
}
