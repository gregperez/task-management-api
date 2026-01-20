package dto

import "gregperez/task-management-api/internal/domain"

// LoginRequest representa la solicitud de inicio de sesión.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa la respuesta de inicio de sesión exitoso.
type LoginResponse struct {
	Token              string       `json:"token"`
	User               *domain.User `json:"user"`
	MustChangePassword bool         `json:"must_change_password"`
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
