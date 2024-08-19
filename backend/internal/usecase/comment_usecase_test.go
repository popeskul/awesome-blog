package usecase_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository/mocks"
)

func TestCreateComment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
	postRepo := mocksrepository.NewMockPostRepository(ctrl)
	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewCommentUseCase(commentRepo, postRepo, userRepo, logger)

	newComment := &entity.NewComment{
		AuthorId: authorId1,
		PostId:   postId1,
		Content:  "This is a comment",
	}

	createdComment := &entity.Comment{
		Id:        commentId1,
		AuthorId:  authorId1,
		PostId:    postId1,
		Content:   "This is a comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo.EXPECT().
		GetUserById(gomock.Any(), newComment.AuthorId).
		Return(&entity.User{Id: newComment.AuthorId}, nil).Times(1)

	postRepo.EXPECT().
		GetPostById(gomock.Any(), newComment.PostId).
		Return(&entity.Post{Id: newComment.PostId}, nil).Times(1)

	commentRepo.EXPECT().
		CreateComment(gomock.Any(), newComment).
		Return(createdComment, nil).Times(1)

	result, err := uc.CreateComment(context.Background(), newComment)
	assert.NoError(t, err)
	assert.Equal(t, createdComment, result)
}

func TestCreateComment_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(commentRepo *mocksrepository.MockCommentRepository, postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository)
		comment       *entity.NewComment
		expectedError string
	}{
		{
			name: "Invalid new comment",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository, postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Times(0)

				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Times(0)

				commentRepo.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			comment: &entity.NewComment{
				Content:  "",
				AuthorId: authorId1,
				PostId:   postId1,
			},
			expectedError: "invalid comment",
		},
		{
			name: "User not found",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository, postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(nil, usecase.ErrUserNotFound).Times(1)
			},
			comment: &entity.NewComment{
				AuthorId: authorId1,
				PostId:   postId1,
				Content:  "This is a comment",
			},
			expectedError: "failed to get user",
		},
		{
			name: "Post not found",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository, postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{
						Id: authorId1,
					}, nil).Times(1)
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			comment: &entity.NewComment{
				AuthorId: authorId1,
				PostId:   postId1,
				Content:  "asd",
			},
			expectedError: "failed to get post",
		},
		{
			name: "Failed to create comment",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository, postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{
						Id: authorId1,
					}, nil).Times(1)
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(&entity.Post{
						Id: postId1,
					}, nil).Times(1)
				commentRepo.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			comment: &entity.NewComment{
				AuthorId: authorId1,
				PostId:   postId1,
				Content:  "This is a comment",
			},
			expectedError: "failed to create comment: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
			postRepo := mocksrepository.NewMockPostRepository(ctrl)
			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewCommentUseCase(commentRepo, postRepo, userRepo, logger)

			tt.mockSetup(commentRepo, postRepo, userRepo)

			result, err := uc.CreateComment(context.Background(), tt.comment)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetCommentByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewCommentUseCase(commentRepo, nil, nil, logger)

	expectedComment := &entity.Comment{
		Id:        commentId1,
		PostId:    postId1,
		AuthorId:  authorId1,
		Content:   "This is a comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	commentRepo.EXPECT().
		GetCommentById(gomock.Any(), expectedComment.Id).
		Return(expectedComment, nil).Times(1)

	result, err := uc.GetCommentByID(context.Background(), expectedComment.Id)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, result)
}

func TestGetCommentByID_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(commentRepo *mocksrepository.MockCommentRepository)
		commentID     uuid.UUID
		expectedError string
	}{
		{
			name: "Comment not found",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository) {
				commentRepo.EXPECT().
					GetCommentById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			commentID:     commentId1,
			expectedError: "failed to get comment: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewCommentUseCase(commentRepo, nil, nil, logger)

			tt.mockSetup(commentRepo)

			result, err := uc.GetCommentByID(context.Background(), tt.commentID)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetComments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewCommentUseCase(commentRepo, nil, nil, logger)

	pagination := &entity.Pagination{Page: 1, Limit: 10}
	expectedComments := &entity.CommentList{
		Data: []*entity.Comment{
			{
				Id:       commentId1,
				PostId:   postId1,
				AuthorId: authorId1,
				Content:  "Comment 1",
			},
			{
				Id:       commentId2,
				PostId:   postId2,
				AuthorId: authorId2,
				Content:  "Comment 2",
			},
		},
	}

	commentRepo.EXPECT().
		GetComments(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedComments, nil).Times(1)

	result, err := uc.GetComments(context.Background(), postId1, pagination)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, result)
}

func TestGetComments_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(commentRepo *mocksrepository.MockCommentRepository)
		postID        uuid.UUID
		pagination    *entity.Pagination
		expectedError string
	}{
		{
			name: "Invalid pagination parameters",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository) {
			},
			postID: postId1,
			pagination: &entity.Pagination{
				Page:  -1,
				Limit: 10,
			},
			expectedError: "invalid pagination",
		},
		{
			name: "Failed to get comments",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository) {
				commentRepo.EXPECT().
					GetComments(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			postID:        postId1,
			pagination:    &entity.Pagination{Page: 1, Limit: 10},
			expectedError: "failed to get comments: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewCommentUseCase(commentRepo, nil, nil, logger)

			tt.mockSetup(commentRepo)

			result, err := uc.GetComments(context.Background(), tt.postID, tt.pagination)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUpdateComment_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(commentRepo *mocksrepository.MockCommentRepository)
		comment       *entity.UpdateComment
		expectedError string
	}{
		{
			name: "Invalid comment content",
			mockSetup: func(commentRepo *mocksrepository.MockCommentRepository) {
				// No interactions with repository expected
			},
			comment:       &entity.UpdateComment{Content: ""}, // Invalid content
			expectedError: usecase.ErrInvalidComment.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commentRepo := mocksrepository.NewMockCommentRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewCommentUseCase(commentRepo, nil, nil, logger)

			tt.mockSetup(commentRepo)

			err := uc.UpdateComment(context.Background(), tt.comment)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
