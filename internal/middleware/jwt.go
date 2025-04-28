package middleware

import (
	"context"
	"github.com/gekich/go-web/internal/auth"
	"github.com/gekich/go-web/internal/util/response"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey = contextKey("userID")

func JWTAuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, "missing auth header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.WriteError(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			claims, err := jwtManager.Verify(parts[1])
			if err != nil {
				response.WriteError(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
