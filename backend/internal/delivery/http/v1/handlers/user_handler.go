package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/gen/api"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
	"github.com/popeskul/awesome-blog/backend/internal/validator"
)

type UsrHandler struct {
	userUseCase usecase.UseCaseUser
	validator   validator.Validator
	logger      *logrus.Logger
}

func NewUserHandler(userUseCase usecase.UseCaseUser, logger *logrus.Logger, validator validator.Validator) *UsrHandler {
	return &UsrHandler{
		userUseCase: userUseCase,
		logger:      logger,
		validator:   validator,
	}
}

func (h *UsrHandler) GetApiV1Users(w http.ResponseWriter, r *http.Request, params api.GetApiV1UsersParams) {
	h.logger.Printf("GetApiV1Users: Starting to handle request")

	ctx := r.Context()

	pagination, err := entity.NewPaginationFromParams(entity.RemoteParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   (*string)(params.Sort),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.userUseCase.GetAllUsers(ctx, pagination)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get users")
		respondError(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	pagination.Total = users.Pagination.Total

	respondJSON(w, http.StatusOK, users)
	h.logger.Printf("GetApiV1Users: Finished handling request")
}

func (h *UsrHandler) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
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
	createdUser, err := h.userUseCase.CreateUser(ctx, &newUser)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create user")
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, createdUser)
}

func (h *UsrHandler) DeleteApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	ctx := r.Context()

	err := h.userUseCase.DeleteUserByID(ctx, userId)
	if err != nil {
		h.logger.WithError(err).WithField("userId", userId).Error("Failed to delete user")
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UsrHandler) GetApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	ctx := r.Context()

	foundUser, err := h.userUseCase.GetUserByID(ctx, userId)
	if err != nil {
		h.logger.WithError(err).WithField("userId", userId).Error("Failed to get user")
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, foundUser)
}

func (h *UsrHandler) PutApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var updateUser entity.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&updateUser); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	ctx := r.Context()
	updatedUser, err := h.userUseCase.UpdateUserByID(ctx, userId, &updateUser)
	if err != nil {
		h.logger.WithError(err).WithField("userId", userId).Error("Failed to update user")
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}
