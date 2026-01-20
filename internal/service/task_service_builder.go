package service

import (
	"time"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/service/dto"

	"github.com/google/uuid"
)

// generateID genera un nuevo UUID como string.
// Esta función facilita el mocking en tests.
func generateID() string {
	return uuid.New().String()
}

// now retorna el timestamp actual.
// Esta función facilita el mocking en tests.
func now() time.Time {
	return time.Now()
}

// buildNewTask construye una nueva tarea con los valores proporcionados
// y los valores por defecto apropiados.
func buildNewTask(req dto.CreateTaskRequest, creatorID string) *domain.Task {
	timestamp := now()
	return &domain.Task{
		ID:          generateID(),
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusAssigned,
		DueDate:     req.DueDate,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   creatorID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}
}

// buildNewComment construye un nuevo comentario con los valores proporcionados.
func buildNewComment(taskID, userID, content string) *domain.Comment {
	return &domain.Comment{
		ID:        generateID(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now(),
	}
}

// updateTaskFields actualiza los campos de una tarea desde el request.
// Solo actualiza los campos que no están vacíos.
func updateTaskFields(task *domain.Task, req dto.UpdateTaskRequest) {
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}
}

// markTaskAsUpdated actualiza el timestamp UpdatedAt de la tarea.
func markTaskAsUpdated(task *domain.Task) {
	task.UpdatedAt = now()
}

// updateTaskStatus actualiza el estado de una tarea y marca el timestamp.
func updateTaskStatus(task *domain.Task, status domain.TaskStatus) {
	task.Status = status
	markTaskAsUpdated(task)
}
