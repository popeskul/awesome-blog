package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository"
	"github.com/sirupsen/logrus"
)

type postUseCase struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	logger   *logrus.Logger
}

func NewPostUseCase(postRepo repository.PostRepository, userRepo repository.UserRepository, logger *logrus.Logger) UseCasePost {
	return &postUseCase{
		postRepo: postRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *postUseCase) CreatePost(ctx context.Context, post *entity.NewPost) (*entity.Post, error) {
	if err := ValidatePost(post); err != nil {
		return nil, err
	}

	if _, err := uc.userRepo.GetUserById(ctx, post.AuthorId); err != nil {
		uc.logger.WithError(err).WithField("userID", post.AuthorId).Error("Failed to get user")
		return nil, ErrUserNotFound
	}

	createdPost, err := uc.postRepo.CreatePost(ctx, post)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create post")
		return nil, err
	}

	return &entity.Post{
		Id:        createdPost.Id,
		Title:     createdPost.Title,
		Content:   createdPost.Content,
		AuthorId:  createdPost.AuthorId,
		CreatedAt: createdPost.CreatedAt,
		UpdatedAt: createdPost.UpdatedAt,
	}, nil
}

func (uc *postUseCase) GetPost(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	post, err := uc.postRepo.GetPostById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("postID", id).Error("Failed to get post")
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (uc *postUseCase) GetAllPosts(ctx context.Context, params *entity.Pagination) (*entity.Response[entity.Post], error) {
	if err := entity.ValidatePagination(params); err != nil {
		return nil, err
	}

	posts, err := uc.postRepo.GetAll(ctx, params)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get posts")
		return nil, err
	}

	total, err := uc.postRepo.GetTotalPosts(ctx)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get total posts")
		return nil, err
	}

	return &entity.Response[entity.Post]{
		Data: posts,
		Pagination: &entity.Pagination{
			Total:  int(total),
			Page:   params.Page,
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	}, nil
}

func (uc *postUseCase) UpdatePost(ctx context.Context, post *entity.Post, userID uuid.UUID) error {
	if err := ValidatePost(post); err != nil {
		return fmt.Errorf("invalid post data: %w", err)
	}

	existingPost, err := uc.postRepo.GetPostById(ctx, post.Id)
	if err != nil {
		uc.logger.WithError(err).WithField("postID", post.Id).Error("Failed to get post")
		return ErrPostNotFound
	}

	user, err := uc.userRepo.GetUserById(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).WithField("userID", post.AuthorId).Error("Failed to get user")
		return ErrUserNotFound
	}

	if existingPost.AuthorId != user.Id {
		return ErrUnauthorized
	}

	post.UpdatedAt = time.Now()

	if err := uc.postRepo.Update(ctx, post); err != nil {
		uc.logger.WithError(err).WithField("postID", post.Id).Error("Failed to update post")
		return err
	}

	return nil
}

func (uc *postUseCase) DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if _, err := uc.postRepo.GetPostById(ctx, id); err != nil {
		uc.logger.WithError(err).WithField("postID", id).Error("Failed to get post")
		return ErrPostNotFound
	}

	if _, err := uc.userRepo.GetUserById(ctx, userID); err != nil {
		uc.logger.WithError(err).WithField("userID", userID).Error("Failed to get user")
		return ErrUserNotFound
	}

	if err := uc.postRepo.Delete(ctx, id); err != nil {
		uc.logger.WithError(err).WithField("postID", id).Error("Failed to delete post")
		return err
	}

	return nil
}
