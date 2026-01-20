package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/middleware"
	"gregperez/task-management-api/pkg/helper"
)

// Errores personalizados para mejor manejo
var (
	ErrInvalidJSON     = errors.New("entrada JSON inválida")
	ErrMissingUserID   = errors.New("ID de usuario no encontrado en contexto")
	ErrMissingUserRole = errors.New("rol de usuario no encontrado en contexto")
	ErrInvalidUserID   = errors.New("ID de usuario inválido en contexto")
	ErrInvalidUserRole = errors.New("rol de usuario inválido en contexto")
)

// GetUserID extrae de forma segura el ID del usuario del contexto
func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return "", ErrMissingUserID
	}
	return userID, nil
}

// GetUserRole extrae de forma segura el rol del usuario del contexto
func GetUserRole(ctx context.Context) (domain.UserRole, error) {
	userRole, ok := ctx.Value(middleware.UserRoleKey).(domain.UserRole)
	if !ok {
		return "", ErrMissingUserRole
	}
	return userRole, nil
}

// GetUserContext extrae de forma segura ambos valores (ID y rol) del contexto
func GetUserContext(ctx context.Context) (string, domain.UserRole, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return "", "", err
	}

	userRole, err := GetUserRole(ctx)
	if err != nil {
		return "", "", err
	}

	return userID, userRole, nil
}

// DecodeJSON decodifica el body del request a la estructura proporcionada
func DecodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return ErrInvalidJSON
	}
	return nil
}

// MapErrorToStatus mapea errores del dominio a códigos de estado HTTP
func MapErrorToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// Errores de dominio específicos
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict
	}

	// Errores de handler
	switch err {
	case ErrInvalidJSON, ErrInvalidUserID, ErrInvalidUserRole:
		return http.StatusBadRequest
	case ErrMissingUserID, ErrMissingUserRole:
		return http.StatusUnauthorized
	}

	// Error genérico
	return http.StatusInternalServerError
}

// RespondWithData envía una respuesta JSON exitosa con datos
func RespondWithData(w http.ResponseWriter, statusCode int, data interface{}) {
	helper.RespondJSON(w, statusCode, data)
}

// RespondWithMessage envía una respuesta JSON exitosa con un mensaje simple
func RespondWithMessage(w http.ResponseWriter, statusCode int, message string) {
	helper.RespondJSON(w, statusCode, map[string]string{"message": message})
}

// RespondWithError envía una respuesta de error con el código HTTP apropiado
func RespondWithError(w http.ResponseWriter, err error) {
	statusCode := MapErrorToStatus(err)
	helper.RespondError(w, statusCode, err.Error())
}
