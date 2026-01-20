package handler

import (
	"encoding/json"
	"net/http"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/middleware"
	"gregperez/task-management-api/internal/service"
	"gregperez/task-management-api/pkg/helper"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	user, tempPassword, err := h.userService.CreateUser(r.Context(), req, userRole)
	if err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"user":               user,
		"temporary_password": tempPassword,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		helper.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	var req service.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	if err := h.userService.UpdateUser(r.Context(), userID, req, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "usuario actualizado"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	if err := h.userService.DeleteUser(r.Context(), userID, userRole); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "usuario eliminado"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value(middleware.UserRoleKey).(domain.UserRole)

	users, err := h.userService.ListUsers(r.Context(), userRole)
	if err != nil {
		helper.RespondError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, users)
}
