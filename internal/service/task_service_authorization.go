package service

import "gregperez/task-management-api/internal/domain"

// requireRole verifica que el usuario tenga el rol requerido.
// Retorna ErrForbidden si el rol no coincide.
func requireRole(userRole domain.UserRole, requiredRole domain.UserRole) error {
	if userRole != requiredRole {
		return domain.ErrForbidden
	}
	return nil
}

// canViewTask verifica si un usuario puede ver una tarea específica.
// Los auditores pueden ver todas las tareas.
// Los ejecutores solo pueden ver tareas asignadas a ellos.
// Los administradores pueden ver todas las tareas.
func canViewTask(task *domain.Task, userID string, userRole domain.UserRole) error {
	// Auditores y administradores pueden ver todas las tareas
	if userRole == domain.RoleAuditor || userRole == domain.RoleAdmin {
		return nil
	}

	// Ejecutores solo pueden ver sus tareas asignadas
	if userRole == domain.RoleExecutor && task.AssignedTo == userID {
		return nil
	}

	return domain.ErrForbidden
}

// canModifyTask verifica si un usuario con el rol dado puede modificar tareas.
// Solo los administradores pueden modificar tareas.
func canModifyTask(userRole domain.UserRole) error {
	return requireRole(userRole, domain.RoleAdmin)
}

// isTaskAssignedToUser verifica si la tarea está asignada al usuario especificado.
// Retorna ErrNotAssignedToUser si la tarea no está asignada al usuario.
func isTaskAssignedToUser(task *domain.Task, userID string) error {
	if task.AssignedTo != userID {
		return domain.ErrNotAssignedToUser
	}
	return nil
}

// canUpdateTaskStatus verifica si un ejecutor puede actualizar el estado de una tarea.
// Solo ejecutores pueden actualizar estados y solo de tareas asignadas a ellos.
func canUpdateTaskStatus(task *domain.Task, userID string, userRole domain.UserRole) error {
	if err := requireRole(userRole, domain.RoleExecutor); err != nil {
		return err
	}
	return isTaskAssignedToUser(task, userID)
}

// canAddComment verifica si un usuario puede agregar comentarios a una tarea.
// Solo ejecutores pueden agregar comentarios.
func canAddComment(userRole domain.UserRole) error {
	return requireRole(userRole, domain.RoleExecutor)
}
