package usecase_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository/mocks"
)

func TestCreatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepo := mocksrepository.NewMockPostRepository(ctrl)
	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewPostUseCase(postRepo, userRepo, logger)

	newPost := &entity.NewPost{
		AuthorId: authorId1,
		Title:    "Test Title",
		Content:  "Test Content",
	}

	createdPost := &entity.Post{
		Id:      postId1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	userRepo.EXPECT().
		GetUserById(gomock.Any(), newPost.AuthorId).
		Return(&entity.User{Id: newPost.AuthorId}, nil).Times(1)

	postRepo.EXPECT().
		CreatePost(gomock.Any(), newPost).
		Return(createdPost, nil).Times(1)

	expectedPost := &entity.Post{
		Id:      createdPost.Id,
		Title:   createdPost.Title,
		Content: createdPost.Content,
	}

	createdPost, err := uc.CreatePost(context.Background(), newPost)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, createdPost)
}

func TestCreatePost_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository, postRepo *mocksrepository.MockPostRepository)
		newPost       *entity.NewPost
		expectedError string
		expectedPost  *entity.Post
	}{
		{
			name: "User not found",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, postRepo *mocksrepository.MockPostRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("user not found")).Times(1)
			},
			newPost: &entity.NewPost{
				AuthorId: authorId1,
				Title:    "Test Title",
				Content:  "Test Content",
			},
			expectedError: usecase.ErrUserNotFound.Error(),
			expectedPost:  nil,
		},
		{
			name: "Failed to create post",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, postRepo *mocksrepository.MockPostRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{Id: authorId1}, nil).Times(1)
				postRepo.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("failed to create post")).Times(1)
			},
			newPost: &entity.NewPost{
				AuthorId: authorId1,
				Title:    "Test Title",
				Content:  "Test Content",
			},
			expectedError: "failed to create post",
			expectedPost:  nil,
		},
		{
			name: "Invalid post",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, postRepo *mocksrepository.MockPostRepository) {
				// No setup required as the error is expected to come from validation
			},
			newPost: &entity.NewPost{
				AuthorId: authorId1,
				Title:    "", // Invalid title
				Content:  "Test Content",
			},
			expectedError: usecase.ErrEmptyTitleOrContent.Error(),
			expectedPost:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postRepo := mocksrepository.NewMockPostRepository(ctrl)
			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewPostUseCase(postRepo, userRepo, logger)

			tt.mockSetup(userRepo, postRepo)

			createdPost, err := uc.CreatePost(context.Background(), tt.newPost)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPost, createdPost)
		})
	}
}

func TestGetPost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepo := mocksrepository.NewMockPostRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewPostUseCase(postRepo, nil, logger)

	expectedPost := &entity.Post{
		Id:      postId1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	postRepo.EXPECT().
		GetPostById(gomock.Any(), postId1).
		Return(expectedPost, nil).Times(1)

	foundPost, err := uc.GetPost(context.Background(), postId1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, foundPost)
}

func TestGetPost_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(postRepo *mocksrepository.MockPostRepository)
		postID        uuid.UUID
		expectedError string
		expectedPost  *entity.Post
	}{
		{
			name: "Post not found",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository) {
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("post not found")).Times(1)
			},
			postID:        postId1,
			expectedError: usecase.ErrPostNotFound.Error(),
			expectedPost:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postRepo := mocksrepository.NewMockPostRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewPostUseCase(postRepo, nil, logger)

			tt.mockSetup(postRepo)

			foundPost, err := uc.GetPost(context.Background(), tt.postID)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPost, foundPost)
		})
	}
}

func TestGetAllPosts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepo := mocksrepository.NewMockPostRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewPostUseCase(postRepo, nil, logger)

	paginationParams := &entity.Pagination{
		Page:   1,
		Limit:  10,
		Offset: 0,
		Total:  100,
	}

	expectedPosts := []*entity.Post{
		{Id: postId1, Title: "Post 1", Content: "Content 1"},
		{Id: postId2, Title: "Post 2", Content: "Content 2"},
	}

	postRepo.EXPECT().
		GetAll(gomock.Any(), paginationParams).
		Return(expectedPosts, nil).Times(1)

	postRepo.EXPECT().
		GetTotalPosts(gomock.Any()).
		Return(paginationParams.Total, nil).Times(1)

	result, err := uc.GetAllPosts(context.Background(), paginationParams)

	assert.NoError(t, err)
	assert.Equal(t, &entity.PostList{
		Posts: expectedPosts,
		Pagination: entity.Pagination{
			Total:  paginationParams.Total,
			Page:   paginationParams.Page,
			Limit:  paginationParams.Limit,
			Offset: paginationParams.Offset,
		},
	}, result)
}

func TestGetAllPosts_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(postRepo *mocksrepository.MockPostRepository)
		params        *entity.Pagination
		expectedError string
		expectedPosts *entity.PostList
	}{
		{
			name: "Invalid pagination parameters",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository) {
				// No interactions expected
			},
			params:        nil,
			expectedError: "pagination params cannot be nil",
			expectedPosts: nil,
		},
		{
			name: "Failed to get posts",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository) {
				postRepo.EXPECT().
					GetAll(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("failed to get posts: db error")).Times(1)
			},
			params:        &entity.Pagination{Page: 1, Limit: 10},
			expectedError: "failed to get posts: db error",
			expectedPosts: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postRepo := mocksrepository.NewMockPostRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewPostUseCase(postRepo, nil, logger)

			tt.mockSetup(postRepo)

			result, err := uc.GetAllPosts(context.Background(), tt.params)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPosts, result)
		})
	}
}

func TestUpdatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepo := mocksrepository.NewMockPostRepository(ctrl)
	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewPostUseCase(postRepo, userRepo, logger)

	updatedPost := &entity.Post{
		Id:       postId1,
		Title:    "Updated Title",
		Content:  "Updated Content",
		AuthorId: authorId1,
	}

	existingPost := &entity.Post{
		Id:       postId1,
		AuthorId: authorId1,
	}

	user := &entity.User{
		Id: authorId1,
	}

	postRepo.EXPECT().
		GetPostById(gomock.Any(), updatedPost.Id).
		Return(existingPost, nil).Times(1)

	userRepo.EXPECT().
		GetUserById(gomock.Any(), user.Id).
		Return(user, nil).Times(1)

	postRepo.EXPECT().
		Update(gomock.Any(), updatedPost).
		Return(nil).Times(1)

	err := uc.UpdatePost(context.Background(), updatedPost, user.Id)
	assert.NoError(t, err)
}

func TestUpdatePost_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository)
		post          *entity.Post
		userID        uuid.UUID
		expectedError string
	}{
		{
			name: "Invalid post",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				// No interactions expected
			},
			post:          &entity.Post{Title: ""}, // Invalid post
			userID:        authorId1,
			expectedError: "invalid post data: title and content are required",
		},
		{
			name: "Post not found",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("post not found")).Times(1)
			},
			post: &entity.Post{
				Id:      postId1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			userID:        authorId1,
			expectedError: usecase.ErrPostNotFound.Error(), // Adjusted to match the updated error
		},
		{
			name: "User not found",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(&entity.Post{Id: postId1, AuthorId: authorId1}, nil).Times(1)
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			post: &entity.Post{
				Id:      postId1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			userID:        authorId1,
			expectedError: usecase.ErrUserNotFound.Error(), // Adjusted to match the updated error
		},
		{
			name: "Failed to update post",
			mockSetup: func(postRepo *mocksrepository.MockPostRepository, userRepo *mocksrepository.MockUserRepository) {
				postRepo.EXPECT().
					GetPostById(gomock.Any(), gomock.Any()).
					Return(&entity.Post{Id: postId1, AuthorId: authorId1}, nil).Times(1)
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{Id: authorId1}, nil).Times(1)
				postRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("failed to update post: db error")).Times(1)
			},
			post: &entity.Post{
				Id:      postId1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			userID:        authorId1,
			expectedError: "failed to update post: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postRepo := mocksrepository.NewMockPostRepository(ctrl)
			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewPostUseCase(postRepo, userRepo, logger)

			tt.mockSetup(postRepo, userRepo)

			err := uc.UpdatePost(context.Background(), tt.post, tt.userID)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
