package handler

import (
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/handler/request"
	"gregperez/task-management-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser crea un nuevo usuario (solo Admin)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.CreateUserRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := service.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Role:     domain.UserRole(req.Role),
	}

	user, tempPassword, err := h.userService.CreateUser(ctx, serviceReq, userRole)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusCreated, map[string]interface{}{
		"user":               user,
		"temporary_password": tempPassword,
	})
}

// GetUser obtiene un usuario por ID (solo Admin)
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "id")

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, user)
}

// UpdateUser actualiza un usuario existente (solo Admin)
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "id")
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.UpdateUserRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := service.UpdateUserRequest{
		Email:    req.Email,
		Role:     domain.UserRole(req.Role),
	}

	if err := h.userService.UpdateUser(ctx, userID, serviceReq, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "usuario actualizado")
}

// DeleteUser elimina un usuario (solo Admin)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "id")
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	if err := h.userService.DeleteUser(ctx, userID, userRole); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "usuario eliminado")
}

// ListUsers lista todos los usuarios (solo Admin)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRole, err := GetUserRole(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	users, err := h.userService.ListUsers(ctx, userRole)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, users)
}
