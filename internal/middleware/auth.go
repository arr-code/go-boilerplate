package middleware

import (
	"errors"

	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and injects user claims into context
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		tokenString, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			utils.SendError(c, 401, "Unauthorized", errors.New("missing or invalid authorization header"))
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.SendError(c, 401, "Unauthorized", errors.New("invalid or expired token"))
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}
