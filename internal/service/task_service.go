package service

import (
	"context"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"
	"gregperez/task-management-api/internal/service/dto"
)

type TaskService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
}

func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// CreateTask crea una nueva tarea asignada a un ejecutor.
// Solo los administradores pueden crear tareas.
func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest, creatorID string, creatorRole domain.UserRole) (*domain.Task, error) {
	if err := canModifyTask(creatorRole); err != nil {
		return nil, err
	}

	assignee, err := s.userRepo.GetByID(ctx, req.AssignedTo)
	if err != nil {
		return nil, err
	}

	if err := validateAssigneeIsExecutor(assignee); err != nil {
		return nil, err
	}

	task := buildNewTask(req, creatorID)

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask actualiza los campos de una tarea existente.
// Solo los administradores pueden actualizar tareas y solo si están en estado "Asignado".
func (s *TaskService) UpdateTask(ctx context.Context, taskID string, req dto.UpdateTaskRequest, updaterRole domain.UserRole) error {
	if err := canModifyTask(updaterRole); err != nil {
		return err
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if err := validateTaskModification(task); err != nil {
		return err
	}

	updateTaskFields(task, req)
	markTaskAsUpdated(task)

	return s.taskRepo.Update(ctx, task)
}

// DeleteTask elimina una tarea existente.
// Solo los administradores pueden eliminar tareas y solo si están en estado "Asignado".
func (s *TaskService) DeleteTask(ctx context.Context, taskID string, deleterRole domain.UserRole) error {
	if err := canModifyTask(deleterRole); err != nil {
		return err
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if err := validateTaskModification(task); err != nil {
		return err
	}

	return s.taskRepo.Delete(ctx, taskID)
}

// UpdateTaskStatus actualiza el estado de una tarea.
// Solo los ejecutores pueden actualizar el estado de sus tareas asignadas.
// La tarea no debe estar vencida y la transición de estado debe ser válida.
func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID string, req dto.UpdateTaskStatusRequest, userID string, userRole domain.UserRole) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if err := canUpdateTaskStatus(task, userID, userRole); err != nil {
		return err
	}

	if err := validateTaskNotExpired(task); err != nil {
		return err
	}

	if err := validateStatusTransition(task, req.Status); err != nil {
		return err
	}

	updateTaskStatus(task, req.Status)

	return s.taskRepo.Update(ctx, task)
}

// AddComment agrega un comentario a una tarea.
// Solo los ejecutores pueden agregar comentarios y solo a tareas vencidas.
func (s *TaskService) AddComment(ctx context.Context, taskID string, content string, userID string, userRole domain.UserRole) error {
	if err := canAddComment(userRole); err != nil {
		return err
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if err := validateCommentOnExpiredTask(task); err != nil {
		return err
	}

	comment := buildNewComment(taskID, userID, content)

	return s.taskRepo.AddComment(ctx, comment)
}

// ListUserTasks retorna todas las tareas asignadas a un usuario específico.
func (s *TaskService) ListUserTasks(ctx context.Context, userID string) ([]*domain.Task, error) {
	return s.taskRepo.ListByAssignee(ctx, userID)
}

// ListAllTasks retorna todas las tareas del sistema.
// Solo los auditores pueden listar todas las tareas.
func (s *TaskService) ListAllTasks(ctx context.Context, requesterRole domain.UserRole) ([]*domain.Task, error) {
	if err := requireRole(requesterRole, domain.RoleAuditor); err != nil {
		return nil, err
	}

	return s.taskRepo.ListAll(ctx)
}

// GetTask obtiene una tarea por su ID.
// Los ejecutores solo pueden ver tareas asignadas a ellos.
// Los auditores y administradores pueden ver todas las tareas.
func (s *TaskService) GetTask(ctx context.Context, taskID string, userID string, userRole domain.UserRole) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if err := canViewTask(task, userID, userRole); err != nil {
		return nil, err
	}

	return task, nil
}
