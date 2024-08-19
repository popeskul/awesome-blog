package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/gen/api"
	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/popeskul/awesome-blog/backend/internal/delivery/http/v1/handlers"
	"github.com/popeskul/awesome-blog/backend/internal/delivery/http/v1/middleware"
)

type Handler interface {
	handlers.PostHandlers
	handlers.CommentHandlers
	handlers.UserHandlers
	handlers.AuthHandlers
}

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	logger     *logrus.Logger
	handler    Handler
	staticPath string
}

func NewServer(cfg *config.Config, logger *logrus.Logger, handler Handler, staticPath string) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		handler:    handler,
		staticPath: staticPath,
	}
}

func (s *Server) Run() error {
	router := s.setupRouter()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler: router,
	}

	s.logger.Infof("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Server is shutting down...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(s.cfg, s.logger))
		r.Get("/auth/logout", s.handler.PostAuthLogout)
		api.HandlerFromMuxWithBaseURL(s.handler, r, "")
	})

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", s.handler.PostAuthLogin)
		r.Post("/auth/register", s.handler.PostAuthRegister)

		r.Get("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
			queryParams := r.URL.Query()

			pageStr := queryParams.Get("page")
			limitStr := queryParams.Get("limit")
			offsetStr := queryParams.Get("offset")
			sortStr := queryParams.Get("sort")

			var page, limit, offset int
			var sort api.GetApiV1PostsParamsSort

			if pageStr != "" {
				page, _ = strconv.Atoi(pageStr)
			}
			if limitStr != "" {
				limit, _ = strconv.Atoi(limitStr)
			}
			if offsetStr != "" {
				offset, _ = strconv.Atoi(offsetStr)
			}
			if sortStr != "" {
				sort = api.GetApiV1PostsParamsSort(sortStr)
			}

			s.handler.GetApiV1Posts(w, r, api.GetApiV1PostsParams{
				Page:   &page,
				Limit:  &limit,
				Offset: &offset,
				Sort:   &sort,
			})
		})

		r.Get("/api/v1/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
			postId, err := uuid.Parse(chi.URLParam(r, "postId"))
			if err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			s.handler.GetApiV1PostsPostId(w, r, postId)
		})

		r.Get("/api/v1/posts/{postId}/comments", func(w http.ResponseWriter, r *http.Request) {
			postId, err := uuid.Parse(chi.URLParam(r, "postId"))
			if err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}

			queryParams := r.URL.Query()

			pageStr := queryParams.Get("page")
			limitStr := queryParams.Get("limit")
			offsetStr := queryParams.Get("offset")
			sortStr := queryParams.Get("sort")

			var page, limit, offset int
			var sort api.GetApiV1PostsPostIdCommentsParamsSort

			if pageStr != "" {
				page, _ = strconv.Atoi(pageStr)
			}
			if limitStr != "" {
				limit, _ = strconv.Atoi(limitStr)
			}
			if offsetStr != "" {
				offset, _ = strconv.Atoi(offsetStr)
			}
			if sortStr != "" {
				sort = api.GetApiV1PostsPostIdCommentsParamsSort(sortStr)
			}

			s.handler.GetApiV1PostsPostIdComments(w, r, postId, api.GetApiV1PostsPostIdCommentsParams{
				Page:   &page,
				Limit:  &limit,
				Offset: &offset,
				Sort:   &sort,
			})
		})
		r.Handle("/swagger/*", handlers.SwaggerHandler(s.staticPath))
	})

	return r
}
