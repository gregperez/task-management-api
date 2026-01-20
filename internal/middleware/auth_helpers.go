package middleware

import (
	"slices"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gregperez/task-management-api/internal/domain"
)

var (
	// ErrMissingToken indica que no se proporcionó un token de autorización
	ErrMissingToken = errors.New("token de autorización faltante")
	
	// ErrInvalidTokenFormat indica que el formato del token es inválido
	ErrInvalidTokenFormat = errors.New("formato de token inválido")
)

// AuthErrorResponse representa la estructura de error para respuestas de autenticación
type AuthErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ExtractBearerToken extrae el token JWT del header Authorization
// Espera el formato: "Bearer <token>"
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", ErrInvalidTokenFormat
	}

	if parts[0] != BearerScheme {
		return "", ErrInvalidTokenFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrInvalidTokenFormat
	}

	return token, nil
}

// HasRequiredRole verifica si el usuario tiene uno de los roles permitidos
func HasRequiredRole(userRole domain.UserRole, allowedRoles []domain.UserRole) bool {
	return slices.Contains(allowedRoles, userRole)
}

// RespondAuthError envía una respuesta JSON de error de autenticación
func RespondAuthError(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := AuthErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    code,
	}
	
	json.NewEncoder(w).Encode(response)
}

// MapJWTError mapea errores de JWT a códigos y mensajes específicos
func MapJWTError(err error) (message, code string) {
	errMsg := err.Error()
	
	// Detectar tipos comunes de errores JWT
	switch {
	case strings.Contains(errMsg, "expired"):
		return MsgExpiredToken, ErrCodeExpiredToken
	case strings.Contains(errMsg, "malformed"):
		return MsgInvalidTokenFormat, ErrCodeInvalidTokenFormat
	default:
		return MsgInvalidToken, ErrCodeInvalidToken
	}
}
