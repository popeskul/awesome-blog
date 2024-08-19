package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository"
)

type commentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	logger      *logrus.Logger
}

func NewCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	logger *logrus.Logger,
) UseCaseComment {
	return &commentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (uc *commentUseCase) CreateComment(ctx context.Context, comment *entity.NewComment) (*entity.Comment, error) {
	if comment == nil || comment.AuthorId.ID() == 0 || comment.PostId.ID() == 0 || comment.Content == "" {
		return nil, ErrInvalidComment
	}

	if _, err := uc.userRepo.GetUserById(ctx, comment.AuthorId); err != nil {
		uc.logger.WithError(err).WithField("userID", comment.AuthorId).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if _, err := uc.postRepo.GetPostById(ctx, comment.PostId); err != nil {
		uc.logger.WithError(err).WithField("postID", comment.PostId).Error("Failed to get post")
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	createdComment, err := uc.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create comment")
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"commentID": createdComment,
		"postID":    createdComment.PostId,
		"authorID":  createdComment.AuthorId,
		"content":   createdComment.Content,
		"createdAt": createdComment.CreatedAt,
		"updatedAt": createdComment.UpdatedAt,
		"author_id": createdComment.AuthorId,
	}).Info("Comment created successfully")

	return createdComment, nil
}

func (uc *commentUseCase) GetCommentByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	comment, err := uc.commentRepo.GetCommentById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("commentID", id).Error("Failed to get comment")
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

func (uc *commentUseCase) GetComments(ctx context.Context, postID uuid.UUID, pagination *entity.Pagination) (*entity.Response[entity.Comment], error) {
	if err := uc.validatePagination(pagination); err != nil {
		return nil, fmt.Errorf("invalid pagination: %w", err)
	}

	comments, err := uc.commentRepo.GetComments(ctx, postID, pagination)
	if err != nil {
		uc.logger.WithError(err).WithField("postID", postID).Error("Failed to get comments")
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	total, err := uc.commentRepo.GetTotalCommentsByPostID(ctx, postID)
	if err != nil {
		uc.logger.WithError(err).WithField("postID", postID).Error("Failed to get total comments")
		return nil, fmt.Errorf("failed to get total comments: %w", err)
	}

	return &entity.Response[entity.Comment]{
		Data: comments,
		Pagination: &entity.Pagination{
			Total:  total,
			Page:   pagination.Page,
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		},
	}, nil
}

func (uc *commentUseCase) UpdateComment(ctx context.Context, comment *entity.UpdateComment) error {
	// TODO: Add validation for comment content
	if comment == nil || comment.Id.ID() == 0 || comment.Content == "" {
		return ErrInvalidComment
	}

	existingComment, err := uc.commentRepo.GetCommentById(ctx, comment.Id)
	if err != nil {
		uc.logger.WithError(err).WithField("commentID", comment.Id).Error("Failed to get comment")
		return ErrCommentNotFound
	}

	if existingComment.AuthorId != comment.AuthorId {
		return ErrUnauthorized
	}

	existingComment.Content = comment.Content
	existingComment.UpdatedAt = time.Now()

	err = uc.commentRepo.UpdateComment(ctx, &entity.UpdateComment{
		Id:      existingComment.Id,
		Content: existingComment.Content,
	})
	if err != nil {
		uc.logger.WithError(err).WithField("commentID", comment.Id).Error("Failed to update comment")
		return fmt.Errorf("failed to update comment: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"commentID": comment.Id,
		"content":   comment.Content,
		"updatedAt": existingComment.UpdatedAt,
		"author_id": existingComment.AuthorId,
	}).Info("Comment updated successfully")

	return nil
}

func (uc *commentUseCase) DeleteComment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	comment, err := uc.commentRepo.GetCommentById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("commentID", id).Error("Failed to get comment")
		return ErrCommentNotFound
	}

	if comment.AuthorId != userID {
		return ErrUnauthorized
	}

	err = uc.commentRepo.DeleteCommentById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("commentID", id).Error("Failed to delete comment")
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	uc.logger.WithField("commentID", id).Info("Comment deleted successfully")

	return nil
}

func (uc *commentUseCase) validatePagination(pagination *entity.Pagination) error {
	if pagination.Page <= 0 {
		return ErrInvalidPage
	}

	if pagination.Limit <= 0 {
		return ErrInvalidLimit
	}

	return nil
}
