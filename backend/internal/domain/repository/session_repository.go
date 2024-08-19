package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) (*entity.Session, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error)
	UpdateSession(ctx context.Context, session *entity.Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
}
