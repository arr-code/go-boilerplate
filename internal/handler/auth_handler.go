package handler

import (
	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/service"
	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		// Map specific errors to appropriate status codes
		switch err {
		case model.ErrEmailAlreadyExists, model.ErrUsernameAlreadyExists, model.ErrUserAlreadyExists:
			utils.SendError(c, 409, "Registration failed", err)
		case model.ErrInvalidEmail, model.ErrInvalidUsername, model.ErrWeakPassword:
			utils.SendError(c, 400, "Validation failed", err)
		default:
			utils.SendError(c, 500, "Registration failed", err)
		}
		return
	}

	utils.SendSuccess(c, 201, "User registered successfully", resp)
}

// Login handles user authentication
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		// Map specific errors to appropriate status codes
		switch err {
		case model.ErrInvalidCredentials:
			utils.SendError(c, 401, "Login failed", err)
		case model.ErrUserInactive:
			utils.SendError(c, 403, "Login failed", err)
		default:
			utils.SendError(c, 500, "Login failed", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Login successful", resp)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		// Map specific errors to appropriate status codes
		switch err {
		case model.ErrInvalidToken:
			utils.SendError(c, 401, "Token refresh failed", err)
		case model.ErrUserInactive:
			utils.SendError(c, 403, "Token refresh failed", err)
		default:
			utils.SendError(c, 500, "Token refresh failed", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Token refreshed successfully", resp)
}

// GetMe retrieves the current authenticated user's information
// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "Unauthorized", model.ErrUnauthorized)
		return
	}

	resp, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		if err == model.ErrUserNotFound {
			utils.SendError(c, 404, "User not found", err)
		} else {
			utils.SendError(c, 500, "Failed to retrieve user", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "User retrieved successfully", resp)
}
