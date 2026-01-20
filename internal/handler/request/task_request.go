package request

import "time"

// TaskRequests representa las estructuras de solicitud para operaciones de tareas

// CreateTaskRequest representa la solicitud para crear una nueva tarea
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     time.Time `json:"due_date"`
	AssignedTo  string `json:"assigned_to"`
}

// UpdateTaskRequest representa la solicitud para actualizar una tarea
type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

// UpdateTaskStatusRequest representa la solicitud para actualizar solo el estado
type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

// AddCommentRequest representa la solicitud para agregar un comentario
type AddCommentRequest struct {
	Content string `json:"content"`
}
