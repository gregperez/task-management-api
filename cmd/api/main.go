package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"gregperez/task-management-api/internal/app"
	"gregperez/task-management-api/internal/config"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Conectar a la base de datos
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer pool.Close()

	// Inicializar aplicación con todas sus dependencias
	application := app.New(cfg, pool)
	defer application.Close()

	// Crear y configurar servidor
	server := application.NewServer()

	// Iniciar servidor
	server.Start()

	// Esperar señal de apagado y ejecutar graceful shutdown
	server.WaitForShutdown()
}
