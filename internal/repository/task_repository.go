package repository

import (
	"context"
	"gregperez/task-management-api/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
	ListByAssignee(ctx context.Context, userID string) ([]*domain.Task, error)
	ListAll(ctx context.Context) ([]*domain.Task, error)
	AddComment(ctx context.Context, comment *domain.Comment) error
}
