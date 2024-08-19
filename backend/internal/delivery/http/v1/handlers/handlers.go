package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/popeskul/awesome-blog/backend/gen/api"
)

type PostHandlers interface {
	GetApiV1Posts(w http.ResponseWriter, r *http.Request, params api.GetApiV1PostsParams)
	PostApiV1Posts(w http.ResponseWriter, r *http.Request)
	DeleteApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID)
	GetApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID)
	PutApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID)
}

type CommentHandlers interface {
	GetApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID, params api.GetApiV1PostsPostIdCommentsParams)
	PostApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID)
}

type UserHandlers interface {
	GetApiV1Users(w http.ResponseWriter, r *http.Request, params api.GetApiV1UsersParams)
	PostApiV1Users(w http.ResponseWriter, r *http.Request)
	DeleteApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID)
	GetApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID)
	PutApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID)
}

type AuthHandlers interface {
	PostAuthRegister(w http.ResponseWriter, r *http.Request)
	PostAuthLogin(w http.ResponseWriter, r *http.Request)
	GetAuthMe(w http.ResponseWriter, r *http.Request)
	PostAuthLogout(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	api.Unimplemented

	postHandlers    PostHandlers
	commentHandlers CommentHandlers
	userHandlers    UserHandlers
	authHandlers    AuthHandlers
}

func NewHandler(
	postHandler PostHandlers,
	commentHandler CommentHandlers,
	userHandler UserHandlers,
	authHandler AuthHandlers,
) *Handler {
	return &Handler{
		postHandlers:    postHandler,
		commentHandlers: commentHandler,
		userHandlers:    userHandler,
		authHandlers:    authHandler,
	}
}

func (h *Handler) GetApiV1Posts(w http.ResponseWriter, r *http.Request, params api.GetApiV1PostsParams) {
	h.postHandlers.GetApiV1Posts(w, r, params)
}

func (h *Handler) PostApiV1Posts(w http.ResponseWriter, r *http.Request) {
	h.postHandlers.PostApiV1Posts(w, r)
}

func (h *Handler) DeleteApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	h.postHandlers.DeleteApiV1PostsPostId(w, r, postId)
}

func (h *Handler) GetApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	h.postHandlers.GetApiV1PostsPostId(w, r, postId)
}

func (h *Handler) PutApiV1PostsPostId(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	h.postHandlers.PutApiV1PostsPostId(w, r, postId)
}

func (h *Handler) GetApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID, params api.GetApiV1PostsPostIdCommentsParams) {
	h.commentHandlers.GetApiV1PostsPostIdComments(w, r, postId, params)
}

func (h *Handler) PostApiV1PostsPostIdComments(w http.ResponseWriter, r *http.Request, postId uuid.UUID) {
	h.commentHandlers.PostApiV1PostsPostIdComments(w, r, postId)
}

func (h *Handler) GetApiV1Users(w http.ResponseWriter, r *http.Request, params api.GetApiV1UsersParams) {
	h.userHandlers.GetApiV1Users(w, r, params)
}

func (h *Handler) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	h.userHandlers.PostApiV1Users(w, r)
}

func (h *Handler) DeleteApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	h.userHandlers.DeleteApiV1UsersUserId(w, r, userId)
}

func (h *Handler) GetApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	h.userHandlers.GetApiV1UsersUserId(w, r, userId)
}

func (h *Handler) PutApiV1UsersUserId(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	h.userHandlers.PutApiV1UsersUserId(w, r, userId)
}

func (h *Handler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	h.authHandlers.PostAuthLogin(w, r)
}

func (h *Handler) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	h.authHandlers.PostAuthRegister(w, r)
}

func (h *Handler) GetAuthMe(w http.ResponseWriter, r *http.Request) {
	h.authHandlers.GetAuthMe(w, r)
}

func (h *Handler) PostAuthLogout(w http.ResponseWriter, r *http.Request) {
	h.authHandlers.PostAuthLogout(w, r)
}
