package middleware

import (
	"context"
	"net/http"
	"strings"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/pkg/jwt"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "token no proporcionado", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "formato de token inválido", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ValidateToken(parts[1], jwtSecret)
			if err != nil {
				http.Error(w, "token inválido", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, domain.UserRole(claims.Role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(UserRoleKey).(domain.UserRole)
			if !ok {
				http.Error(w, "rol no encontrado", http.StatusUnauthorized)
				return
			}

			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "acceso denegado", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
