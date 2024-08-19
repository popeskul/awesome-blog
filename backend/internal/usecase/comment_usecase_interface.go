package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

//go:generate mockgen -destination=mocks/mock_comment_usecase.go -package=mockusecase github.com/popeskul/awesome-blog/backend/internal/usecase UseCaseComment

type UseCaseComment interface {
	CreateComment(ctx context.Context, comment *entity.NewComment) (*entity.Comment, error)
	GetComments(ctx context.Context, postID uuid.UUID, pagination *entity.Pagination) (*entity.Response[entity.Comment], error)
	UpdateComment(ctx context.Context, comment *entity.UpdateComment) error
	DeleteComment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetCommentByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
}
