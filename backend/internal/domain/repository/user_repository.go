package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_user_repository.go -package=mocksrepository github.com/popeskul/awesome-blog/backend/internal/domain/repository UserRepository

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.NewUser) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetAllUsers(ctx context.Context, params *entity.Pagination) ([]*entity.User, error)
	GetTotalUsers(ctx context.Context) (int, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, user *entity.UpdateUser) error
}
