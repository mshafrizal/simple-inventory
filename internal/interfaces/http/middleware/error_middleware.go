package middleware

import (
	"log"
	"net/http"
	"simple-inventory/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Err)

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_server_error",
				Message: err.Error(),
			})
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "internal_server_error",
					Message: "An unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
