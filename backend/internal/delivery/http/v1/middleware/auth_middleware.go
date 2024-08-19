package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/config"
	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

// AuthMiddleware creates a middleware handler for authentication
func AuthMiddleware(cfg *config.Config, logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			publicPaths := []string{
				"/api/v1/posts",
				"/api/v1/posts/*/comments",
				"/api/v1/users",
			}

			for _, path := range publicPaths {
				if r.Method == http.MethodGet && (r.URL.Path == path || (strings.HasSuffix(path, "*") && strings.HasPrefix(r.URL.Path, strings.TrimSuffix(path, "*")))) {
					logger.WithField("path", r.URL.Path).Info("AuthMiddleware: Public path, skipping authentication")
					next.ServeHTTP(w, r)
					return
				}
			}

			tokenStr := extractTokenFromHeader(r)
			if tokenStr == "" {
				logger.Error("AuthMiddleware: No token provided")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(tokenStr, cfg)
			if err != nil {
				logger.WithError(err).Error("AuthMiddleware: Invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "session_id", claims.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func parseToken(tokenStr string, cfg *config.Config) (*entity.Session, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	session := &entity.Session{}

	if sessionIDStr, ok := claims["session_id"].(string); ok {
		if sessionID, err := uuid.Parse(sessionIDStr); err == nil {
			session.SessionID = sessionID
		}
	}

	if userIDStr, ok := claims["user_id"].(string); ok {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			session.UserID = userID
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		session.ExpiresAt = time.Unix(int64(exp), 0)
	}

	if session.UserID == uuid.Nil {
		return nil, fmt.Errorf("invalid token: missing user_id")
	}

	return session, nil
}
