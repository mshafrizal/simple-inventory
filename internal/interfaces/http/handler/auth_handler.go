package handler

import (
	"net/http"
	"simple-inventory/internal/interfaces/http/dto"
	"simple-inventory/internal/interfaces/http/util"
	"simple-inventory/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "registration_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "User registered successfully",
		Data: dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	session, err := h.authUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "login_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token:     session.Token,
		ExpiresAt: util.FormatTimeUTC(session.ExpiresAt),
		User: dto.UserDTO{
			ID:       session.User.ID,
			Username: session.User.Username,
			Email:    session.User.Email,
			Role:     session.User.Role,
			IsActive: session.User.IsActive,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "missing_token"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.authUseCase.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "logout_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Logged out successfully"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}
