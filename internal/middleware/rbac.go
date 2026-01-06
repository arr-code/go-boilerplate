package middleware

import (
	"errors"

	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has at least one of the specified roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get roles from context (set by AuthMiddleware)
		userRoles, exists := c.Get("roles")
		if !exists {
			utils.SendError(c, 403, "Forbidden", errors.New("no roles found in token"))
			c.Abort()
			return
		}

		// Convert to string slice
		rolesList, ok := userRoles.([]string)
		if !ok {
			utils.SendError(c, 403, "Forbidden", errors.New("invalid role format"))
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range rolesList {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			utils.SendError(c, 403, "Forbidden", errors.New("insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a convenience wrapper for requiring admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireModerator is a convenience wrapper for requiring moderator or admin role
func RequireModerator() gin.HandlerFunc {
	return RequireRole("moderator", "admin")
}
