package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/popeskul/awesome-blog/backend/internal/delivery/http/v1/handlers"
	"github.com/popeskul/awesome-blog/backend/internal/hash"
	"github.com/popeskul/awesome-blog/backend/internal/infrastructure/database/postgres"
	"github.com/popeskul/awesome-blog/backend/internal/server"
	"github.com/popeskul/awesome-blog/backend/internal/usecase"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
	"github.com/popeskul/awesome-blog/backend/pkg/migrator"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := config.LoadConfig([]string{"/app/config/"})
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	migrationsPath := filepath.Join(".", "migrations")
	if err = migrator.MigrateUp(cfg.Database.DSN(), migrationsPath); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	database, err := db.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err = database.Close(); err != nil {
			logger.Errorf("Error closing database connection: %v", err)
		}
	}()

	postRepo := postgres.NewPostRepository(database, logger)
	commentRepo := postgres.NewCommentRepository(database, logger)
	userRepo := postgres.NewUserRepository(database, logger)
	sessionRepo := postgres.NewSessionRepository(database, logger)

	hashService := &hash.BcryptHashService{}
	validatorService := validator.New()

	postUseCase := usecase.NewPostUseCase(postRepo, userRepo, logger)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, postRepo, userRepo, logger)
	userUseCase := usecase.NewUserUseCase(userRepo, logger, hashService)
	authUseCase := usecase.NewAuthUseCase(userRepo, sessionRepo, logger, cfg, hashService)

	postHandler := handlers.NewPostHandler(postUseCase, logger, validatorService)
	commentHandler := handlers.NewCommentHandler(commentUseCase, logger, validatorService)
	userHandler := handlers.NewUserHandler(userUseCase, logger, validatorService)
	authHandler := handlers.NewAuthHandler(authUseCase, userUseCase, logger, validatorService)

	handler := handlers.NewHandler(postHandler, commentHandler, userHandler, authHandler)

	logger.Info("Starting server...")

	staticPath := filepath.Join("/app", "static")

	srv := server.NewServer(cfg, logger, handler, staticPath)

	go func() {
		if err = srv.Run(); err != nil {
			logger.Fatalf("Error occurred while running server: %v", err)
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Errorf("Error occurred while shutting down server: %v", err)
	}

	if err = database.HealthCheck(ctx); err != nil {
		logger.Errorf("Database health check failed: %v", err)
	} else {
		logger.Info("Database health check passed")
	}

	logger.Info("Server exited properly")
}
