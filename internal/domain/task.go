package domain

import (
	"errors"
	"time"
)

type TaskStatus string

const (
	StatusAssigned   TaskStatus = "Asignado"
	StatusInProgress TaskStatus = "Iniciado"
	StatusPending    TaskStatus = "En espera"
	StatusSuccess    TaskStatus = "Finalizado Éxito"
	StatusError      TaskStatus = "Finalizado Error"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     time.Time  `json:"due_date"`
	AssignedTo  string     `json:"assigned_to"`
	CreatedBy   string     `json:"created_by"`
	Comments    []Comment  `json:"comments,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Comment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired verifica si la tarea está vencida
func (t *Task) IsExpired() bool {
	return time.Now().After(t.DueDate)
}

// CanTransitionTo valida las transiciones de estado según el diagrama
func (t *Task) CanTransitionTo(newStatus TaskStatus) error {
	validTransitions := map[TaskStatus][]TaskStatus{
		StatusAssigned:   {StatusInProgress},
		StatusInProgress: {StatusSuccess, StatusError, StatusPending},
		StatusPending:    {StatusInProgress},
		StatusSuccess:    {},
		StatusError:      {},
	}

	allowed, exists := validTransitions[t.Status]
	if !exists {
		return errors.New("estado actual no válido")
	}

	for _, valid := range allowed {
		if valid == newStatus {
			return nil
		}
	}

	return errors.New("transición de estado no permitida")
}

// CanBeModified verifica si la tarea puede ser modificada o eliminada
func (t *Task) CanBeModified() bool {
	return t.Status == StatusAssigned
}
