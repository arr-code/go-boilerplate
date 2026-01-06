package handler

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Health checks the health of the application and its dependencies
// GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	// Check database connection
	if err := h.db.Ping(); err != nil {
		c.JSON(503, gin.H{
			"status":    "unhealthy",
			"database":  "disconnected",
			"error":     err.Error(),
			"timestamp": time.Now().Unix(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": time.Now().Unix(),
	})
}
