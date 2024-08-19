package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository"
	"github.com/popeskul/awesome-blog/backend/internal/hash"
	"github.com/sirupsen/logrus"
)

type useCase struct {
	userRepo repository.UserRepository
	hash     hash.HashService
	logger   *logrus.Logger
}

func NewUserUseCase(userRepo repository.UserRepository, logger *logrus.Logger, hash hash.HashService) UseCaseUser {
	return &useCase{
		userRepo: userRepo,
		logger:   logger,
		hash:     hash,
	}
}

func (uc *useCase) CreateUser(ctx context.Context, user *entity.NewUser) (*entity.User, error) {
	if user.Username == "" || user.PasswordHash == "" {
		return nil, ErrEmptyCredentials
	}

	existingUser, err := uc.userRepo.GetUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := uc.hash.HashPassword(user.PasswordHash)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = hashedPassword

	createdUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create user")
		return nil, ErrCreateUserFailed
	}

	return createdUser, nil
}

func (uc *useCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.GetUserById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("userID", id).Error("Failed to get user")
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (uc *useCase) GetAllUsers(ctx context.Context, pagination *entity.Pagination) (*entity.Response[entity.User], error) {
	users, err := uc.userRepo.GetAllUsers(ctx, pagination)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get all users")
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	total, err := uc.userRepo.GetTotalUsers(ctx)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get total users count")
		return nil, fmt.Errorf("failed to get total users count: %w", err)
	}

	return &entity.Response[entity.User]{
		Data: users,
		Pagination: &entity.Pagination{
			Total: total,
			Limit: pagination.Limit,
			Page:  pagination.Page,
		},
	}, nil
}

func (uc *useCase) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	user, err := uc.userRepo.GetUserById(ctx, id)
	if err != nil {
		uc.logger.WithError(err).WithField("userID", id).Error("Failed to get user for deletion")
		return ErrUserNotFound
	}

	if err = uc.userRepo.DeleteUserById(ctx, user.Id); err != nil {
		uc.logger.WithError(err).WithField("userID", id).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (uc *useCase) UpdateUserByID(ctx context.Context, userID uuid.UUID, updateUser *entity.UpdateUser) (*entity.User, error) {
	existingUser, err := uc.userRepo.GetUserById(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).WithField("userID", userID).Error("Failed to get user for update")
		return nil, ErrUserNotFound
	}

	if updateUser.Username != "" {
		existingUser.Username = updateUser.Username
	}

	if updateUser.Password != "" {
		hashedPassword, err := uc.hash.HashPassword(updateUser.Password)
		if err != nil {
			uc.logger.WithError(err).Error("Failed to hash password")
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		existingUser.PasswordHash = hashedPassword
	}

	if updateUser.Email != "" {
		existingUser.Email = updateUser.Email
	}

	if err = uc.userRepo.UpdateUser(ctx, &entity.UpdateUser{
		Id:       userID,
		Username: existingUser.Username,
		Password: existingUser.PasswordHash,
		Email:    existingUser.Email,
	}); err != nil {
		uc.logger.WithError(err).WithField("userID", userID).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}
