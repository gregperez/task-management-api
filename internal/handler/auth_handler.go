package handler

import (
	"encoding/json"
	"net/http"

	"gregperez/task-management-api/internal/middleware"
	"gregperez/task-management-api/internal/service"
	"gregperez/task-management-api/pkg/helper"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		helper.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req service.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, "entrada inválida")
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, req); err != nil {
		helper.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "contraseña actualizada exitosamente"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// En JWT stateless, el logout se maneja en el cliente eliminando el token
	helper.RespondJSON(w, http.StatusOK, map[string]string{"message": "sesión cerrada exitosamente"})
}
