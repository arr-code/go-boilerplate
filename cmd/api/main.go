package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xnoia-go-boilerplate/internal/config"
	"xnoia-go-boilerplate/internal/handler"
	"xnoia-go-boilerplate/internal/repository"
	"xnoia-go-boilerplate/internal/service"
	"xnoia-go-boilerplate/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// Connect to PostgreSQL database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	log.Println("Database connected successfully")

	// Initialize repository layer
	userRepo := repository.NewUserRepository(db)

	// Initialize service layer
	authService := service.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.JWT.RefreshExpiration,
	)
	userService := service.NewUserService(userRepo)
	adminService := service.NewAdminService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminService)
	healthHandler := handler.NewHealthHandler(db)

	// Setup router with all handlers
	router := handler.SetupRouter(
		cfg,
		db,
		authHandler,
		userHandler,
		adminHandler,
		healthHandler,
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s (mode: %s)", cfg.Server.Port, cfg.Server.GinMode)
		log.Printf("Server running at http://localhost:%s", cfg.Server.Port)
		log.Printf("Health check: http://localhost:%s/health", cfg.Server.Port)
		log.Printf("API endpoints: http://localhost:%s/api/v1", cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// SIGINT handles Ctrl+C, SIGTERM handles docker stop
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
