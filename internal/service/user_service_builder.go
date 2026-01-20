package service

import (
	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/service/dto"
	"gregperez/task-management-api/pkg/password"
)

// buildNewUser construye un nuevo usuario con contraseña temporal.
// Retorna el usuario creado y la contraseña temporal en texto plano.
func buildNewUser(req dto.CreateUserRequest) (*domain.User, string, error) {
	tempPassword := password.GenerateTemporaryPassword()
	hashedPassword, err := password.HashPassword(tempPassword)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		ID:                  generateID(),
		Username:            req.Username,
		Email:               req.Email,
		Password:            hashedPassword,
		Role:                req.Role,
		IsTemporaryPassword: true,
		CreatedAt:           now(),
		UpdatedAt:           now(),
	}

	return user, tempPassword, nil
}

// updateUserFields actualiza los campos de un usuario desde el request.
// Solo actualiza los campos que no están vacíos.
func updateUserFields(user *domain.User, req dto.UpdateUserRequest) {
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
}

// markUserAsUpdated actualiza el timestamp UpdatedAt del usuario.
func markUserAsUpdated(user *domain.User) {
	user.UpdatedAt = now()
}
