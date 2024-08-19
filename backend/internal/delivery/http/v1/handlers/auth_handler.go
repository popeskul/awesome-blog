package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
	"github.com/popeskul/awesome-blog/backend/internal/validator"
)

type AuthHandler struct {
	authUseCase usecase.UseCaseAuth
	userUseCase usecase.UseCaseUser
	logger      *logrus.Logger
	validator   validator.Validator
}

func NewAuthHandler(
	authUseCase usecase.UseCaseAuth,
	userUseCase usecase.UseCaseUser,
	logger *logrus.Logger,
	validator validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		logger:      logger,
		validator:   validator,
	}
}

func (h *AuthHandler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var credentials entity.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&credentials); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	ctx := r.Context()
	token, err := h.authUseCase.Authenticate(ctx, credentials.Username, credentials.Password)
	if err != nil {
		h.logger.WithError(err).Error("Authentication failed")
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&newUser); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	ctx := r.Context()
	createdUser, err := h.authUseCase.Register(ctx, newUser)
	if err != nil {
		h.logger.WithError(err).Error("Failed to register user")
		respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err).Error())
		return
	}

	respondJSON(w, http.StatusCreated, createdUser)
}

func (h *AuthHandler) GetAuthMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDRaw := ctx.Value("user_id")

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		h.logger.WithField("user_id_raw", userIDRaw).Error("Failed to get user ID from context")
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userUseCase.GetUserByID(ctx, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user information")
		respondError(w, http.StatusInternalServerError, "Failed to get user information")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) PostAuthLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID, ok := ctx.Value("session_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get session ID from context")
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	h.logger.WithField("logout_session_id", sessionID).Info("Attempting to logout")

	if err := h.authUseCase.Logout(ctx, sessionID); err != nil {
		h.logger.WithError(err).Error("Failed to logout")
		respondError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	h.logger.WithField("logout_session_id", sessionID).Info("Logout successful")
	respondJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
