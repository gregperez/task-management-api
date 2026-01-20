package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gregperez/task-management-api/internal/router"
)

// Server representa el servidor HTTP
type Server struct {
	httpServer *http.Server
}

// NewServer crea y configura un nuevo servidor HTTP
func (a *App) NewServer() *Server {
	// Configurar router con todas las rutas
	routerConfig := &router.Config{
		JWTSecret:   a.Config.JWTSecret,
		AuthHandler: a.AuthHandler,
		UserHandler: a.UserHandler,
		TaskHandler: a.TaskHandler,
	}

	r := router.New(routerConfig)

	// Crear servidor HTTP con configuración de timeouts
	httpServer := &http.Server{
		Addr:         ":" + a.Config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start inicia el servidor HTTP de forma asíncrona
func (s *Server) Start() {
	go func() {
		log.Printf("Servidor iniciado en el puerto %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()
}

// WaitForShutdown espera señales de interrupción y ejecuta graceful shutdown
func (s *Server) WaitForShutdown() {
	// Canal para recibir señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")

	// Contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown graceful del servidor
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Error apagando servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}
