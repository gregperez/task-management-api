package handler

import (
	"net/http"

	"gregperez/task-management-api/internal/handler/request"
	"gregperez/task-management-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login autentica un usuario y devuelve un token JWT
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.LoginRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := h.authService.Login(ctx, serviceReq)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithData(w, http.StatusOK, resp)
}

// ChangePassword cambia la contraseña del usuario autenticado
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := GetUserID(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	var req request.ChangePasswordRequest
	if err := DecodeJSON(r, &req); err != nil {
		RespondWithError(w, err)
		return
	}

	// Convertir request DTO a service request
	serviceReq := service.ChangePasswordRequest{
		CurrentPassword: req.OldPassword,
		NewPassword:     req.NewPassword,
	}

	if err := h.authService.ChangePassword(ctx, userID, serviceReq); err != nil {
		RespondWithError(w, err)
		return
	}

	RespondWithMessage(w, http.StatusOK, "contraseña actualizada exitosamente")
}

// Logout cierra la sesión del usuario (en JWT stateless se maneja en el cliente)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// En JWT stateless, el logout se maneja en el cliente eliminando el token
	RespondWithMessage(w, http.StatusOK, "sesión cerrada exitosamente")
}
