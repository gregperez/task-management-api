package dto

import (
	"time"

	"gregperez/task-management-api/internal/domain"
)

// CreateTaskRequest representa la solicitud para crear una nueva tarea.
type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required,min=3"`
	Description string    `json:"description" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required"`
}

// UpdateTaskRequest representa la solicitud para actualizar una tarea existente.
type UpdateTaskRequest struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

// UpdateTaskStatusRequest representa la solicitud para actualizar el estado de una tarea.
type UpdateTaskStatusRequest struct {
	Status domain.TaskStatus `json:"status" validate:"required"`
}
