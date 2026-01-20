package app

import (
	"testing"
	"context"
	"os"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	
	"gregperez/task-management-api/internal/config"
)

func TestServer_Start(t *testing.T) {
	// Configuración mínima para el servidor
	cfg := &config.Config{
		Port:      "8081", // Puerto diferente para pruebas
		JWTSecret: "testsecret",
	}

	// Crear pool de base de datos simulado (mock)
	pool, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/testdb")
	if err != nil {
		t.Fatalf("Error creando pool de base de datos: %v", err)
	}
	defer pool.Close()

	app := New(cfg, pool)
	server := app.NewServer()
	server.Start()

	// Esperar un momento para asegurar que el servidor esté corriendo
	time.Sleep(1 * time.Second)

	// Enviar señal de shutdown
	go func() {
		time.Sleep(1 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()
}