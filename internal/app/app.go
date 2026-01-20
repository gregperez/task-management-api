package app

import (
	"gregperez/task-management-api/internal/config"
	"gregperez/task-management-api/internal/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

// App encapsula todas las dependencias de la aplicación
type App struct {
	Config      *config.Config
	DB          *pgxpool.Pool
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
}

// New crea una nueva instancia de App con todas las dependencias inicializadas
func New(cfg *config.Config, pool *pgxpool.Pool) *App {
	app := &App{
		Config: cfg,
		DB:     pool,
	}

	// Inicializar dependencias
	app.initDependencies()

	return app
}

// Close cierra todas las conexiones de la aplicación
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
