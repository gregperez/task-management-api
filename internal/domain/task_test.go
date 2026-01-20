package domain

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
)

func TestTask_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus TaskStatus
		newStatus     TaskStatus
		shouldError   bool
	}{
		{
			name:          "De Asignado a Iniciado - válido",
			currentStatus: StatusAssigned,
			newStatus:     StatusInProgress,
			shouldError:   false,
		},
		{
			name:          "De Asignado a Éxito - inválido",
			currentStatus: StatusAssigned,
			newStatus:     StatusSuccess,
			shouldError:   true,
		},
		{
			name:          "De Iniciado a Éxito - válido",
			currentStatus: StatusInProgress,
			newStatus:     StatusSuccess,
			shouldError:   false,
		},
		{
			name:          "De Iniciado a Error - válido",
			currentStatus: StatusInProgress,
			newStatus:     StatusError,
			shouldError:   false,
		},
		{
			name:          "De Iniciado a En espera - válido",
			currentStatus: StatusInProgress,
			newStatus:     StatusPending,
			shouldError:   false,
		},
		{
			name:          "De En espera a Iniciado - válido",
			currentStatus: StatusPending,
			newStatus:     StatusInProgress,
			shouldError:   false,
		},
		{
			name:          "De Éxito a cualquier estado - inválido",
			currentStatus: StatusSuccess,
			newStatus:     StatusInProgress,
			shouldError:   true,
		},
		{
			name:		   "Estado no válido",
			currentStatus: "EstadoDesconocido",
			newStatus:     "",
			shouldError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.currentStatus}
			err := task.CanTransitionTo(tt.newStatus)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTask_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		dueDate  time.Time
		expected bool
	}{
		{
			name:     "Tarea vencida",
			dueDate:  time.Now().Add(-2 * time.Hour),
			expected: true,
		},
		{
			name:     "Tarea no vencida",
			dueDate:  time.Now().Add(2 * time.Hour),
			expected: false,
		},
		{
			name:     "Tarea con vencimiento ahora",
			dueDate:  time.Now().Add(-1 * time.Second),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{DueDate: tt.dueDate}
			result := task.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTask_CanBeModified(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{
			name:     "Tarea Asignada puede ser modificada",
			status:   StatusAssigned,
			expected: true,
		},
		{
			name:     "Tarea Iniciada no puede ser modificada",
			status:   StatusInProgress,
			expected: false,
		},
		{
			name:     "Tarea Finalizada no puede ser modificada",
			status:   StatusSuccess,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.status}
			result := task.CanBeModified()
			assert.Equal(t, tt.expected, result)
		})
	}
}