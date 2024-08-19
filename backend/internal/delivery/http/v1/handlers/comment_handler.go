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

type CommentHandler struct {
	commentUseCase usecase.UseCaseComment
	logger         *logrus.Logger
	validator      validator.Validator
}

func NewCommentHandler(commentUseCase usecase.UseCaseComment, logger *logrus.Logger, validator validator.Validator) *CommentHandler {
	return &CommentHandler{
		commentUseCase: commentUseCase,
		logger:         logger,
		validator:      validator,
	}
}

func (h *CommentHandler) GetApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID, params api.GetApiV1PostsPostIdCommentsParams) {
	ctx := r.Context()

	h.logger.WithField("params.Sort", params.Sort).Info("GetApiV1PostsPostIdComments")

	paginationFromParams, err := entity.NewPaginationFromParams(entity.RemoteParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   (*string)(params.Sort),
	})
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pagination from params")
		respondError(w, http.StatusBadRequest, "Invalid request params")
		return
	}

	result, err := h.commentUseCase.GetComments(ctx, postId, paginationFromParams)
	if err != nil {
		h.logger.WithError(err).WithField("postId", postId).Error("Failed to get comments")
		respondError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *CommentHandler) PostApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	var newComment entity.NewComment
	if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&newComment); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	newComment.PostId = postId

	ctx := r.Context()

	createdComment, err := h.commentUseCase.CreateComment(ctx, &newComment)
	if err != nil {
		h.logger.WithError(err).WithField("postId", postId).Error("Failed to create comment")
		respondError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	respondJSON(w, http.StatusCreated, createdComment)
}
