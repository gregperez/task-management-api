package app

import (
	"gregperez/task-management-api/internal/handler"
	"gregperez/task-management-api/internal/repository/postgres"
	"gregperez/task-management-api/internal/service"
)

// initDependencies inicializa todos los repositorios, servicios y handlers
func (a *App) initDependencies() {
	// Inicializar repositorios
	userRepo := postgres.NewUserRepository(a.DB)
	taskRepo := postgres.NewTaskRepository(a.DB)

	// Inicializar servicios
	authService := service.NewAuthService(userRepo, a.Config.JWTSecret, a.Config.TokenExpiry)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo)

	// Inicializar handlers
	a.AuthHandler = handler.NewAuthHandler(authService)
	a.UserHandler = handler.NewUserHandler(userService)
	a.TaskHandler = handler.NewTaskHandler(taskService)
}
