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

type PostHandler struct {
	postUseCase usecase.UseCasePost
	validator   validator.Validator
	logger      *logrus.Logger
}

func NewPostHandler(postUseCase usecase.UseCasePost, logger *logrus.Logger, validator validator.Validator) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
		logger:      logger,
		validator:   validator,
	}
}

func (h *PostHandler) GetApiV1Posts(w http.ResponseWriter, r *http.Request, params api.GetApiV1PostsParams) {
	ctx := r.Context()

	paginationFromParams, err := entity.NewPaginationFromParams(entity.RemoteParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   (*string)(params.Sort),
	})
	if err != nil {
		h.logger.WithError(err).Error("Failed to get pagination from params")
		respondError(w, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	result, err := h.postUseCase.GetAllPosts(ctx, paginationFromParams)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get posts")
		respondError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *PostHandler) PostApiV1Posts(w http.ResponseWriter, r *http.Request) {
	var newPost entity.NewPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&newPost); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	ctx := r.Context()
	createdPost, err := h.postUseCase.CreatePost(ctx, &newPost)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create post")
		respondError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	respondJSON(w, http.StatusCreated, createdPost)
}

func (h *PostHandler) DeleteApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	ctx := r.Context()
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get user_id from context")
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.postUseCase.DeletePost(ctx, postId, userId)
	if err != nil {
		h.logger.WithError(err).WithField("postId", postId).Error("Failed to delete post")
		respondError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) GetApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	ctx := r.Context()
	foundPost, err := h.postUseCase.GetPost(ctx, postId)
	if err != nil {
		h.logger.WithError(err).WithField("postId", postId).Error("Failed to get post")
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondJSON(w, http.StatusOK, foundPost)
}

func (h *PostHandler) PutApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	ctx := r.Context()
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get user_id from context")
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var postPut entity.Post
	if err := json.NewDecoder(r.Body).Decode(&postPut); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(&postPut); err != nil {
		h.logger.WithError(err).Error("Failed to validate request body")
		respondError(w, http.StatusBadRequest, "Validation failed "+err.Error())
		return
	}

	postPut.Id = postId

	err := h.postUseCase.UpdatePost(ctx, &postPut, userId)
	if err != nil {
		h.logger.WithError(err).WithField("postId", postId).Error("Failed to update post")
		respondError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	respondJSON(w, http.StatusOK, postPut)
}
