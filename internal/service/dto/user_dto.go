package dto

import "gregperez/task-management-api/internal/domain"

// CreateUserRequest representa la solicitud para crear un nuevo usuario.
type CreateUserRequest struct {
	Username string          `json:"username" validate:"required,min=3"`
	Email    string          `json:"email" validate:"required,email"`
	Role     domain.UserRole `json:"role" validate:"required,oneof=Ejecutor Auditor"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario existente.
type UpdateUserRequest struct {
	Email string          `json:"email" validate:"omitempty,email"`
	Role  domain.UserRole `json:"role" validate:"omitempty,oneof=Ejecutor Auditor"`
}
