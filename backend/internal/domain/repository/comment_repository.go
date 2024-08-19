package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_comment_repository.go -package=mocksrepository github.com/popeskul/awesome-blog/backend/internal/domain/repository CommentRepository

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.NewComment) (*entity.Comment, error)
	GetCommentById(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	GetComments(ctx context.Context, postID uuid.UUID, pagination *entity.Pagination) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.UpdateComment) error
	DeleteCommentById(ctx context.Context, id uuid.UUID) error
	GetTotalCommentsByPostID(ctx context.Context, postID uuid.UUID) (int, error)
}
