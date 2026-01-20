package app

import (
	"context"
	"gregperez/task-management-api/internal/config"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestApp_New(t *testing.T) {
	// Configuración de prueba
	cfg := &config.Config{
		Port:        "8080",
		DatabaseURL: "postgres://user:password@localhost:5432/testdb",
		JWTSecret:   "testsecret",
		TokenExpiry: 3600,
	}
	
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Error creando pool de conexión: %v", err)
	}
	defer pool.Close()

	// Crear nueva instancia de App
	appInstance := New(cfg, pool)
	if appInstance == nil {
		t.Fatal("App instance is nil")
	}
	if appInstance.Config != cfg {
		t.Error("App Config not set correctly")
	}
	if appInstance.DB != pool {
		t.Error("App DB not set correctly")
	}
	if appInstance.AuthHandler == nil {
		t.Error("AuthHandler not initialized")
	}
	if appInstance.UserHandler == nil {
		t.Error("UserHandler not initialized")
	}
	if appInstance.TaskHandler == nil {
		t.Error("TaskHandler not initialized")
	}
}

func TestApp_Close(t *testing.T) {
	// Configuración de prueba
	cfg := &config.Config{
		Port:        "8080",
		DatabaseURL: "postgres://user:password@localhost:5432/testdb",
		JWTSecret:   "testsecret",
		TokenExpiry: 3600,
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Error creando pool de conexión: %v", err)
	}
	appInstance := New(cfg, pool)

	// Cerrar la aplicación
	appInstance.Close()
	// Verificar que la conexión a la base de datos esté cerrada
	if appInstance.DB.Stat().TotalConns() != 0 {
		t.Error("Database connection pool not closed")
	}
}