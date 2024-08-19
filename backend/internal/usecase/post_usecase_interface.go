package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_post_usecase.go -package=mockusecase github.com/popeskul/awesome-blog/backend/internal/usecase UseCasePost

type UseCasePost interface {
	CreatePost(ctx context.Context, post *entity.NewPost) (*entity.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	GetAllPosts(ctx context.Context, params *entity.Pagination) (*entity.Response[entity.Post], error)
	UpdatePost(ctx context.Context, post *entity.Post, userID uuid.UUID) error
	DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
