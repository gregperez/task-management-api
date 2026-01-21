package service

import (
	"context"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/service/dto"
)

// AuthServiceInterface define los métodos del servicio de autenticación.
type AuthServiceInterface interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error
}

// TaskServiceInterface define los métodos del servicio de tareas.
type TaskServiceInterface interface {
	CreateTask(ctx context.Context, req dto.CreateTaskRequest, creatorID string, creatorRole domain.UserRole) (*domain.Task, error)
	UpdateTask(ctx context.Context, taskID string, req dto.UpdateTaskRequest, updaterRole domain.UserRole) error
	DeleteTask(ctx context.Context, taskID string, deleterRole domain.UserRole) error
	UpdateTaskStatus(ctx context.Context, taskID string, req dto.UpdateTaskStatusRequest, userID string, userRole domain.UserRole) error
	AddComment(ctx context.Context, taskID string, content string, userID string, userRole domain.UserRole) error
	ListUserTasks(ctx context.Context, userID string) ([]*domain.Task, error)
	ListAllTasks(ctx context.Context, requesterRole domain.UserRole) ([]*domain.Task, error)
	GetTask(ctx context.Context, taskID string, userID string, userRole domain.UserRole) (*domain.Task, error)
}

// UserServiceInterface define los métodos del servicio de usuarios.
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest, creatorRole domain.UserRole) (*domain.User, string, error)
	UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest, updaterRole domain.UserRole) error
	DeleteUser(ctx context.Context, userID string, deleterRole domain.UserRole) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	ListUsers(ctx context.Context, requesterRole domain.UserRole) ([]*domain.User, error)
}
