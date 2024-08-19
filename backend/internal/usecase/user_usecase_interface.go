package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_user_usecase.go -package=mockusecase github.com/popeskul/awesome-blog/backend/internal/usecase UseCaseUser

type UseCaseUser interface {
	CreateUser(ctx context.Context, user *entity.NewUser) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetAllUsers(ctx context.Context, params *entity.Pagination) (*entity.Response[entity.User], error)
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	UpdateUserByID(ctx context.Context, userID uuid.UUID, updateUser *entity.UpdateUser) (*entity.User, error)
}
