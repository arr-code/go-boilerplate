package middleware

import (
	"errors"
	"log"
	"runtime/debug"

	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", err, debug.Stack())

				// Return 500 error to client
				utils.SendError(c, 500, "Internal Server Error", errors.New("an unexpected error occurred"))

				c.Abort()
			}
		}()

		c.Next()
	}
}
