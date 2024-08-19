package usecase

import (
	"context"
	"github.com/google/uuid"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_auth_usecase.go -package=mockusecase github.com/popeskul/awesome-blog/backend/internal/usecase UseCaseAuth

type UseCaseAuth interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, newUser entity.NewUser) (*entity.User, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
}
