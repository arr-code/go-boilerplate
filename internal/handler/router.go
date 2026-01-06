package handler

import (
	"database/sql"

	"xnoia-go-boilerplate/internal/config"
	"xnoia-go-boilerplate/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes and middleware
func SetupRouter(
	cfg *config.Config,
	db *sql.DB,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	adminHandler *AdminHandler,
	healthHandler *HealthHandler,
) *gin.Engine {
	// Set Gin mode based on configuration
	gin.SetMode(cfg.Server.GinMode)

	router := gin.New()

	// Global middleware (applied to all routes)
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint (no authentication required)
	router.GET("/health", healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public authentication routes (no JWT required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (JWT authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// Auth endpoints (requires valid JWT)
			protected.GET("/auth/me", authHandler.GetMe)

			// User profile endpoints (requires valid JWT)
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.POST("/profile/avatar", userHandler.UpdateAvatar)
		}

		// Admin routes (JWT + admin role required)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		admin.Use(middleware.RequireAdmin())
		{
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/roles", adminHandler.AssignRoles)
			admin.DELETE("/users/:id", adminHandler.DeactivateUser)
		}
	}

	return router
}
