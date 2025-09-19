package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sudo-init-do/okies_core/pkg/response"
	"github.com/sudo-init-do/okies_core/pkg/utils"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextRole   contextKey = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Write(w, http.StatusUnauthorized, "missing authorization header", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Write(w, http.StatusUnauthorized, "invalid authorization header", nil)
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			response.Write(w, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		// Add user_id and role to context
		ctx := context.WithValue(r.Context(), ContextUserID, claims["user_id"].(string))
		ctx = context.WithValue(ctx, ContextRole, claims["role"].(string))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
