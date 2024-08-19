package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/domain/repository"
	"github.com/popeskul/awesome-blog/backend/internal/hash"
)

type authUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	logger      *logrus.Logger
	jwtSecret   []byte
	hash        hash.HashService
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	logger *logrus.Logger,
	cfg *config.Config,
	hash hash.HashService,
) UseCaseAuth {
	return &authUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
		jwtSecret:   []byte(cfg.JWT.SecretKey),
		hash:        hash,
	}
}

func (uc *authUseCase) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		uc.logger.WithError(err).WithField("username", username).Error("Failed to get user")
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	if err := uc.hash.ComparePassword(user.PasswordHash, password); err != nil {
		uc.logger.WithError(err).WithField("username", username).Error("Invalid password")
		return "", fmt.Errorf("invalid credentials")
	}

	sessionId := uuid.New()

	token, err := uc.generateToken(user.Id, sessionId)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to generate token")
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	session := &entity.Session{
		SessionID: sessionId,
		UserID:    user.Id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Token:     token,
	}

	if _, err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
		uc.logger.WithError(err).Error("Failed to create session")
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Register(ctx context.Context, newUser entity.NewUser) (*entity.User, error) {
	if _, err := uc.userRepo.GetUserByUsername(ctx, newUser.Username); err == nil {
		uc.logger.WithField("username", newUser.Username).Error("User already exists")
		return nil, fmt.Errorf("user already exists")
	}

	passwordHash, err := uc.hash.HashPassword(newUser.PasswordHash)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.NewUser{
		Username:     newUser.Username,
		Email:        newUser.Email,
		PasswordHash: passwordHash,
	}

	createdUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"userID":   createdUser.Id,
		"username": createdUser.Username,
		"email":    createdUser.Email,
	}).Info("User registered successfully")

	return createdUser, nil
}

func (uc *authUseCase) Logout(ctx context.Context, sessionId uuid.UUID) error {
	uc.logger.WithField("logout_session_id", sessionId).Info("Attempting to logout in usecase")

	session, err := uc.sessionRepo.GetSessionByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.logger.WithField("sessionID", sessionId).Warn("Session not found, considering it already logged out")
			return nil
		}

		uc.logger.WithError(err).Error("Failed to get session")
		return fmt.Errorf("failed to get session: %w", err)
	}

	if err := uc.sessionRepo.DeleteSession(ctx, session.SessionID); err != nil {
		uc.logger.WithError(err).WithField("logout_session_id", sessionId).Error("Failed to delete session")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	uc.logger.WithField("userID", session.UserID).Info("User logged out successfully")

	return nil
}

func (uc *authUseCase) generateToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}
