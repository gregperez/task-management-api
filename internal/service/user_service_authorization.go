package service

import "gregperez/task-management-api/internal/domain"

// requireAdminRole verifica que el usuario tenga rol de administrador.
// Retorna ErrForbidden si el usuario no es administrador.
func requireAdminRole(userRole domain.UserRole) error {
	if userRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}
	return nil
}

// canCreateUser verifica si un usuario puede crear otros usuarios.
// Solo los administradores pueden crear usuarios.
func canCreateUser(userRole domain.UserRole) error {
	return requireAdminRole(userRole)
}

// canUpdateUser verifica si un usuario puede actualizar otros usuarios.
// Solo los administradores pueden actualizar usuarios.
func canUpdateUser(userRole domain.UserRole) error {
	return requireAdminRole(userRole)
}

// canDeleteUser verifica si un usuario puede eliminar otros usuarios.
// Solo los administradores pueden eliminar usuarios.
func canDeleteUser(userRole domain.UserRole) error {
	return requireAdminRole(userRole)
}

// canListUsers verifica si un usuario puede listar todos los usuarios.
// Solo los administradores pueden listar usuarios.
func canListUsers(userRole domain.UserRole) error {
	return requireAdminRole(userRole)
}
