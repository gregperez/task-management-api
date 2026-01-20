package service

import "gregperez/task-management-api/internal/domain"

// validateTaskModification verifica si una tarea puede ser modificada.
// Una tarea solo puede ser modificada si está en estado "Asignado".
func validateTaskModification(task *domain.Task) error {
	if !task.CanBeModified() {
		return domain.ErrTaskCannotBeModified
	}
	return nil
}

// validateStatusTransition valida que la transición de estado sea permitida.
// Utiliza la lógica de negocio del dominio para verificar la transición.
func validateStatusTransition(task *domain.Task, newStatus domain.TaskStatus) error {
	if err := task.CanTransitionTo(newStatus); err != nil {
		return domain.ErrInvalidStatusTransition
	}
	return nil
}

// validateTaskNotExpired verifica que la tarea no esté vencida.
// Retorna ErrTaskExpired si la tarea ha pasado su fecha de vencimiento.
func validateTaskNotExpired(task *domain.Task) error {
	if task.IsExpired() {
		return domain.ErrTaskExpired
	}
	return nil
}

// validateCommentOnExpiredTask verifica que solo se puedan agregar comentarios
// a tareas vencidas (requisito de negocio).
func validateCommentOnExpiredTask(task *domain.Task) error {
	if !task.IsExpired() {
		return domain.ErrInvalidInput
	}
	return nil
}

// validateAssigneeIsExecutor verifica que el usuario asignado a una tarea
// tenga el rol de ejecutor.
func validateAssigneeIsExecutor(user *domain.User) error {
	if !user.IsExecutor() {
		return domain.ErrInvalidInput
	}
	return nil
}
