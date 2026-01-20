package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"gregperez/task-management-api/internal/domain"
	"gregperez/task-management-api/internal/handler"
	"gregperez/task-management-api/internal/middleware"
)

// Config contiene la configuración necesaria para el router
type Config struct {
	JWTSecret   string
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
}

// New crea y configura un nuevo router con todas las rutas de la aplicación
func New(cfg *Config) *chi.Mux {
	r := chi.NewRouter()

	// Configurar middlewares globales
	setupMiddlewares(r)

	// Configurar rutas
	setupPublicRoutes(r, cfg)
	setupProtectedRoutes(r, cfg)
	setupHealthCheck(r)

	return r
}

// setupMiddlewares configura todos los middlewares globales
func setupMiddlewares(r *chi.Mux) {
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
}

// setupPublicRoutes configura las rutas públicas (sin autenticación)
func setupPublicRoutes(r *chi.Mux, cfg *Config) {
	r.Post("/api/v1/login", cfg.AuthHandler.Login)
}

// setupProtectedRoutes configura todas las rutas protegidas por autenticación
func setupProtectedRoutes(r *chi.Mux, cfg *Config) {
	r.Route("/api/v1", func(r chi.Router) {
		// Middleware de autenticación para todas las rutas protegidas
		r.Use(middleware.Auth(cfg.JWTSecret))

		// Rutas de autenticación
		setupAuthRoutes(r, cfg.AuthHandler)

		// Rutas de usuarios (solo Admin)
		setupUserRoutes(r, cfg.UserHandler)

		// Rutas de tareas (según rol)
		setupTaskRoutes(r, cfg.TaskHandler)
	})
}

// setupAuthRoutes configura las rutas de autenticación
func setupAuthRoutes(r chi.Router, authHandler *handler.AuthHandler) {
	r.Post("/logout", authHandler.Logout)
	r.Post("/change-password", authHandler.ChangePassword)
}

// setupUserRoutes configura las rutas de gestión de usuarios (solo Admin)
func setupUserRoutes(r chi.Router, userHandler *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.RequireRole(domain.RoleAdmin))
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.ListUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}

// setupTaskRoutes configura las rutas de tareas según los roles
func setupTaskRoutes(r chi.Router, taskHandler *handler.TaskHandler) {
	r.Route("/tasks", func(r chi.Router) {
		// Rutas Admin: crear, actualizar, eliminar tareas
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAdmin))
			r.Post("/", taskHandler.CreateTask)
			r.Put("/{id}", taskHandler.UpdateTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
		})

		// Rutas Ejecutor: listar sus tareas, actualizar estado, comentar
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleExecutor))
			r.Get("/my-tasks", taskHandler.ListMyTasks)
			r.Patch("/{id}/status", taskHandler.UpdateTaskStatus)
			r.Post("/{id}/comments", taskHandler.AddComment)
		})

		// Rutas Auditor: ver todas las tareas
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(domain.RoleAuditor))
			r.Get("/all", taskHandler.ListAllTasks)
		})

		// Ruta común: ver detalle de tarea (todos los roles autenticados)
		r.Get("/{id}", taskHandler.GetTask)
	})
}

// setupHealthCheck configura el endpoint de health check
func setupHealthCheck(r *chi.Mux) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
