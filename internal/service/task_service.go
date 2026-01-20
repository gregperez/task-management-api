package service

import (
	"context"
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/repository"

	"github.com/google/uuid"
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

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required,min=3"`
	Description string    `json:"description" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required"`
}

func (s *TaskService) CreateTask(ctx context.Context, req CreateTaskRequest, creatorID string, creatorRole domain.UserRole) (*domain.Task, error) {
	// Solo administradores pueden crear tareas
	if creatorRole != domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}

	// Verificar que el usuario asignado existe y es ejecutor
	assignee, err := s.userRepo.GetByID(ctx, req.AssignedTo)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if !assignee.IsExecutor() {
		return nil, domain.ErrInvalidInput
	}

	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusAssigned,
		DueDate:     req.DueDate,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

type UpdateTaskRequest struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID string, req UpdateTaskRequest, updaterRole domain.UserRole) error {
	if updaterRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Solo se puede actualizar si está en estado "Asignado"
	if !task.CanBeModified() {
		return domain.ErrTaskCannotBeModified
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}

	task.UpdatedAt = time.Now()
	return s.taskRepo.Update(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string, deleterRole domain.UserRole) error {
	if deleterRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !task.CanBeModified() {
		return domain.ErrTaskCannotBeModified
	}

	return s.taskRepo.Delete(ctx, taskID)
}

type UpdateTaskStatusRequest struct {
	Status domain.TaskStatus `json:"status" validate:"required"`
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID string, req UpdateTaskStatusRequest, userID string, userRole domain.UserRole) error {
	if userRole != domain.RoleExecutor {
		return domain.ErrForbidden
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Verificar que la tarea esté asignada al usuario
	if task.AssignedTo != userID {
		return domain.ErrNotAssignedToUser
	}

	// No se puede actualizar una tarea vencida
	if task.IsExpired() {
		return domain.ErrTaskExpired
	}

	// Validar transición de estado
	if err := task.CanTransitionTo(req.Status); err != nil {
		return domain.ErrInvalidStatusTransition
	}

	task.Status = req.Status
	task.UpdatedAt = time.Now()
	return s.taskRepo.Update(ctx, task)
}

func (s *TaskService) AddComment(ctx context.Context, taskID string, content string, userID string, userRole domain.UserRole) error {
	if userRole != domain.RoleExecutor {
		return domain.ErrForbidden
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Solo se pueden agregar comentarios a tareas vencidas
	if !task.IsExpired() {
		return domain.ErrInvalidInput
	}

	comment := &domain.Comment{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	return s.taskRepo.AddComment(ctx, comment)
}

func (s *TaskService) ListUserTasks(ctx context.Context, userID string) ([]*domain.Task, error) {
	return s.taskRepo.ListByAssignee(ctx, userID)
}

func (s *TaskService) ListAllTasks(ctx context.Context, requesterRole domain.UserRole) ([]*domain.Task, error) {
	if requesterRole != domain.RoleAuditor {
		return nil, domain.ErrForbidden
	}

	return s.taskRepo.ListAll(ctx)
}

func (s *TaskService) GetTask(ctx context.Context, taskID string, userID string, userRole domain.UserRole) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Ejecutores solo pueden ver sus tareas, auditores pueden ver todas
	if userRole == domain.RoleExecutor && task.AssignedTo != userID {
		return nil, domain.ErrForbidden
	}

	return task, nil
}
