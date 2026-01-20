package service

import (
	"context"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
)

// validateNotAdminRole verifica que el rol no sea Admin.
// Los administradores no pueden ser creados a través de la API normal.
func validateNotAdminRole(role domain.UserRole) error {
	if role == domain.RoleAdmin {
		return domain.ErrCannotCreateAdmin
	}
	return nil
}

// validateUserNotExists verifica que un usuario con el username dado no exista.
// Retorna ErrUserAlreadyExists si el usuario ya existe.
func validateUserNotExists(ctx context.Context, userRepo repository.UserRepository, username string) error {
	existing, _ := userRepo.GetByUsername(ctx, username)
	if existing != nil {
		return domain.ErrUserAlreadyExists
	}
	return nil
}

// validateRoleChange valida que el cambio de rol sea permitido.
// No se permite cambiar un usuario a rol Admin.
func validateRoleChange(newRole domain.UserRole) error {
	if newRole == "" {
		return nil // Sin cambio de rol
	}
	return validateNotAdminRole(newRole)
}
