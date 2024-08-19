package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_post_repository.go -package=mocksrepository github.com/popeskul/awesome-blog/backend/internal/domain/repository PostRepository

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.NewPost) (*entity.Post, error)
	GetPostById(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	GetAll(ctx context.Context, pagination *entity.Pagination) ([]*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalPosts(ctx context.Context) (int64, error)
}
