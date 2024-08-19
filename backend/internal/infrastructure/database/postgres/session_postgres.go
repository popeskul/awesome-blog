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

type SessionRepository struct {
	db     *db.PostgresDB
	logger *logrus.Logger
}

func NewSessionRepository(db *db.PostgresDB, logger *logrus.Logger) *SessionRepository {
	return &SessionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	query := `INSERT INTO sessions (id, session_id, user_id, created_at, updated_at, expires_at, token) 
              VALUES ($1, $2, $3, NOW(), NOW(), $4, $5) 
              RETURNING id, session_id, user_id, created_at, updated_at, expires_at, token`

	var createdSession entity.Session
	err := r.db.QueryRowContext(ctx, query,
		uuid.New(),        // id
		session.SessionID, // session_id
		session.UserID,
		session.ExpiresAt,
		session.Token).Scan(
		&createdSession.ID,
		&createdSession.SessionID,
		&createdSession.UserID,
		&createdSession.CreatedAt,
		&createdSession.UpdatedAt,
		&createdSession.ExpiresAt,
		&createdSession.Token,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &createdSession, nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {
	query := `SELECT id, session_id, user_id, created_at, updated_at, expires_at, token
              FROM sessions WHERE session_id = $1`

	var session entity.Session

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.SessionID,
		&session.UserID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.ExpiresAt,
		&session.Token,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session not found")
		}
		r.logger.WithError(err).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *entity.Session) error {
	query := `UPDATE sessions SET token = $1, expires_at = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, session.Token, session.ExpiresAt, session.ID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update session")
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	r.logger.WithField("delete_session_id", sessionID).Info("Deleting session")
	result, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE session_id = $1`, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.WithError(err).Error("Failed to check rows affected")
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.WithField("delete_session_id", sessionID).Warn("No session found to delete")
		return nil
	}

	r.logger.WithField("delete_session_id", sessionID).Info("Session successfully deleted")

	return nil
}
