package middleware

import (
	"context"
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/pkg/jwt"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// Auth es un middleware que valida el token JWT y agrega información del usuario al contexto
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer token del header Authorization
			authHeader := r.Header.Get(AuthHeaderKey)
			token, err := ExtractBearerToken(authHeader)
			if err != nil {
				if err == ErrMissingToken {
					RespondAuthError(w, http.StatusUnauthorized, MsgMissingToken, ErrCodeMissingToken)
				} else {
					RespondAuthError(w, http.StatusUnauthorized, MsgInvalidTokenFormat, ErrCodeInvalidTokenFormat)
				}
				return
			}

			// Validar token JWT
			claims, err := jwt.ValidateToken(token, jwtSecret)
			if err != nil {
				message, code := MapJWTError(err)
				RespondAuthError(w, http.StatusUnauthorized, message, code)
				return
			}

			// Agregar información del usuario al contexto
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, domain.UserRole(claims.Role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole es un middleware que valida que el usuario tenga uno de los roles permitidos
func RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer rol del contexto de forma segura
			userRole, ok := r.Context().Value(UserRoleKey).(domain.UserRole)
			if !ok {
				RespondAuthError(w, http.StatusUnauthorized, MsgMissingRole, ErrCodeMissingRole)
				return
			}

			// Verificar si el usuario tiene alguno de los roles permitidos
			if !HasRequiredRole(userRole, roles) {
				RespondAuthError(w, http.StatusForbidden, MsgInsufficientPermissions, ErrCodeInsufficientPermissions)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
