package middleware

import (
	"net/http"
	"simple-inventory/internal/interfaces/http/dto"
	"simple-inventory/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "missing_authorization_header"})
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		session, err := authUseCase.ValidateSession(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid_token", Message: err.Error()})
			c.Abort()
			return
		}

		userDTO := dto.UserDTO{
			ID:       session.User.ID,
			Username: session.User.Username,
			Email:    session.User.Email,
			Role:     session.User.Role,
			IsActive: session.User.IsActive,
		}

		c.Set("user", userDTO)
		c.Set("user_id", session.User.ID)
		c.Next()
	}
}
