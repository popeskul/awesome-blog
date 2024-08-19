package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository/mocks"
	"github.com/popeskul/awesome-blog/backend/internal/hash/mocks"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
)

var (
	userId1 = uuid.New()
	userId2 = uuid.New()
)

func TestCreateUser_Success(t *testing.T) {
	mockSetup := func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
		hashSvc.EXPECT().
			HashPassword("password123").
			Return("hashedPassword", nil).Times(1)
		userRepo.EXPECT().
			GetUserByUsername(gomock.Any(), "newuser").
			Return(nil, nil).Times(1)
		userRepo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(&entity.User{
				Id:           userId1,
				Username:     "newuser",
				Email:        "newuser@example.com",
				PasswordHash: "hashedPassword",
			}, nil).Times(1)
	}
	inputUser := &entity.NewUser{
		Username:     "newuser",
		PasswordHash: "password123",
		Email:        "newuser@example.com",
	}
	expectedUser := &entity.User{
		Id:           userId1,
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: "hashedPassword",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	hashSvc := mockshash.NewMockHashService(ctrl)
	logger := logrus.New()
	uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

	mockSetup(userRepo, hashSvc)

	createdUser, err := uc.CreateUser(context.Background(), inputUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, createdUser)
}

func TestCreateUser_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService)
		inputUser     *entity.NewUser
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name: "User already exists",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
				userRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "newuser").
					Return(&entity.User{}, nil).Times(1)
			},
			inputUser: &entity.NewUser{
				Username:     "newuser",
				PasswordHash: "password123",
				Email:        "newuser@example.com",
			},
			expectedUser:  nil,
			expectedError: usecase.ErrUserExists,
		},
		{
			name: "Failed to check if user exists",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
				userRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "newuser").
					Return(&entity.User{
						Username:     "newuser",
						PasswordHash: "password123",
						Email:        "newuser@example.com",
					}, nil).Times(1)
			},
			inputUser: &entity.NewUser{
				Username:     "newuser",
				PasswordHash: "password123",
				Email:        "newuser@example.com",
			},
			expectedUser:  nil,
			expectedError: usecase.ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			hashSvc := mockshash.NewMockHashService(ctrl)
			logger := logrus.New()
			uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

			tt.mockSetup(userRepo, hashSvc)

			createdUser, err := uc.CreateUser(context.Background(), tt.inputUser)

			assert.Equal(t, tt.expectedUser, createdUser)
			assert.ErrorContains(t, tt.expectedError, err.Error())
		})
	}
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	hashSvc := mockshash.NewMockHashService(ctrl)
	logger := logrus.New()
	uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

	expectedUser := &entity.User{
		Id:           userId1,
		Username:     "existinguser",
		Email:        "existinguser@example.com",
		PasswordHash: "hashedPassword",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.EXPECT().
		GetUserById(gomock.Any(), userId1).
		Return(expectedUser, nil).Times(1)

	foundUser, err := uc.GetUserByID(context.Background(), userId1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, foundUser)
}

func TestGetUserByID_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService)
		inputID       uuid.UUID
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name: "Failed to get user",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), userId1).
					Return(nil, errors.New("db error")).Times(1)
			},
			inputID:       userId1,
			expectedUser:  nil,
			expectedError: usecase.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			hashSvc := mockshash.NewMockHashService(ctrl)
			logger := logrus.New()
			uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

			tt.mockSetup(userRepo, hashSvc)

			foundUser, err := uc.GetUserByID(context.Background(), tt.inputID)

			assert.Equal(t, tt.expectedUser, foundUser)
			assert.ErrorContains(t, err, tt.expectedError.Error())
		})
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewUserUseCase(userRepo, logger, nil) // No need for hash service in this test

	pagination := &entity.Pagination{
		Page:  1,
		Limit: 10,
	}

	expectedUsers := []*entity.User{
		{Id: userId1, Username: "user1", Email: "user1@example.com"},
		{Id: userId2, Username: "user2", Email: "user2@example.com"},
	}

	userRepo.EXPECT().
		GetAllUsers(gomock.Any(), pagination).
		Return(expectedUsers, nil).Times(1)

	userRepo.EXPECT().
		GetTotalUsers(gomock.Any()).
		Return(2, nil).Times(1)

	expectedUserList := &entity.UserList{
		Users: expectedUsers,
		Pagination: &entity.Pagination{
			Total: int64(2),
			Limit: pagination.Limit,
			Page:  pagination.Page,
		},
	}

	userList, err := uc.GetAllUsers(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserList, userList)
}

func TestGetAllUsers_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository)
		pagination    *entity.Pagination
		expectedList  *entity.UserList
		expectedError error
	}{
		{
			name: "Error fetching users",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetAllUsers(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			pagination:    &entity.Pagination{Page: 1, Limit: 10},
			expectedList:  nil,
			expectedError: errors.New("failed to get all users: db error"),
		},
		{
			name: "Error fetching total users count",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetAllUsers(gomock.Any(), gomock.Any()).
					Return([]*entity.User{}, nil).Times(1)
				userRepo.EXPECT().
					GetTotalUsers(gomock.Any()).
					Return(0, errors.New("db error")).Times(1)
			},
			pagination:    &entity.Pagination{Page: 1, Limit: 10},
			expectedList:  nil,
			expectedError: errors.New("failed to get total users count: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewUserUseCase(userRepo, logger, nil) // No need for hash service in this test

			tt.mockSetup(userRepo)

			userList, err := uc.GetAllUsers(context.Background(), tt.pagination)

			assert.Equal(t, tt.expectedList, userList)
			assert.ErrorContains(t, err, tt.expectedError.Error())
		})
	}
}

func TestDeleteUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	logger := logrus.New()
	uc := usecase.NewUserUseCase(userRepo, logger, nil) // No need for hash service in this test

	getUser := &entity.User{
		Id:       userId1,
		Username: "user1",
		Email:    "user1@example.com",
	}

	userRepo.EXPECT().
		GetUserById(gomock.Any(), userId1).
		Return(getUser, nil).Times(1)

	userRepo.EXPECT().
		DeleteUserById(gomock.Any(), userId1).
		Return(nil).Times(1)

	err := uc.DeleteUserByID(context.Background(), userId1)

	assert.NoError(t, err)
}

func TestDeleteUserByID_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository)
		userID        uuid.UUID
		expectedError error
	}{
		{
			name: "User not found",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			userID:        userId1,
			expectedError: usecase.ErrUserNotFound,
		},
		{
			name: "Failed to delete user",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{Id: userId1}, nil).Times(1)
				userRepo.EXPECT().
					DeleteUserById(gomock.Any(), gomock.Any()).
					Return(errors.New("db error")).Times(1)
			},
			userID:        userId1,
			expectedError: fmt.Errorf("failed to delete user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			logger := logrus.New()
			uc := usecase.NewUserUseCase(userRepo, logger, nil) // No need for hash service in this test

			tt.mockSetup(userRepo)

			err := uc.DeleteUserByID(context.Background(), tt.userID)

			assert.ErrorContains(t, err, tt.expectedError.Error())
		})
	}
}

func TestUpdateUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocksrepository.NewMockUserRepository(ctrl)
	hashSvc := mockshash.NewMockHashService(ctrl)
	logger := logrus.New()
	uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

	existingUser := &entity.User{
		Id:           userId1,
		Username:     "olduser",
		Email:        "olduser@example.com",
		PasswordHash: "oldpassword",
	}

	updateUser := &entity.UpdateUser{
		Username: "newuser",
		Password: "newpassword",
		Email:    "newuser@example.com",
	}

	hashedPassword := "hashedNewPassword"

	userRepo.EXPECT().
		GetUserById(gomock.Any(), userId1).
		Return(existingUser, nil).Times(1)

	hashSvc.EXPECT().
		HashPassword("newpassword").
		Return(hashedPassword, nil).Times(1)

	userRepo.EXPECT().
		UpdateUser(gomock.Any(), gomock.Any()).
		Return(nil).Times(1)

	expectedUser := &entity.User{
		Id:           userId1,
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: hashedPassword,
	}

	updatedUser, err := uc.UpdateUserByID(context.Background(), userId1, updateUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, updatedUser)
}

func TestUpdateUserByID_Fail(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService)
		userID        uuid.UUID
		updateUser    *entity.UpdateUser
		expectedError string
		expectedUser  *entity.User
	}{
		{
			name: "User not found",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error")).Times(1)
			},
			userID:        userId1,
			updateUser:    &entity.UpdateUser{},
			expectedError: usecase.ErrUserNotFound.Error(),
			expectedUser:  nil,
		},
		{
			name: "Update failure",
			mockSetup: func(userRepo *mocksrepository.MockUserRepository, hashSvc *mockshash.MockHashService) {
				userRepo.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(&entity.User{Id: userId1}, nil).Times(1)
				hashSvc.EXPECT().
					HashPassword("newpassword").
					Return("hashedNewPassword", nil).Times(1)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("db error")).Times(1)
			},
			userID:        userId1,
			updateUser:    &entity.UpdateUser{Password: "newpassword"},
			expectedError: "failed to update user: db error",
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocksrepository.NewMockUserRepository(ctrl)
			hashSvc := mockshash.NewMockHashService(ctrl)
			logger := logrus.New()
			uc := usecase.NewUserUseCase(userRepo, logger, hashSvc)

			tt.mockSetup(userRepo, hashSvc)

			updatedUser, err := uc.UpdateUserByID(context.Background(), tt.userID, tt.updateUser)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUser, updatedUser)
		})
	}
}
