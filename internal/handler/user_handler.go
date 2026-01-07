package handler

import (
	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/service"
	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves the current user's profile
// GET /api/v1/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "Unauthorized", model.ErrUnauthorized)
		return
	}

	resp, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		if err == model.ErrUserNotFound {
			utils.SendError(c, 404, "Profile not found", err)
		} else {
			utils.SendError(c, 500, "Failed to retrieve profile", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Profile retrieved successfully", resp)
}

// UpdateProfile updates the current user's profile
// PUT /api/v1/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "Unauthorized", model.ErrUnauthorized)
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	err := h.userService.UpdateProfile(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if err == model.ErrUserNotFound {
			utils.SendError(c, 404, "User not found", err)
		} else {
			utils.SendError(c, 500, "Failed to update profile", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Profile updated successfully", nil)
}

// UpdateAvatar uploads and updates the user's avatar
// POST /api/v1/profile/avatar
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "Unauthorized", model.ErrUnauthorized)
		return
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.SendError(c, 400, "File upload failed", err)
		return
	}

	// Update avatar
	avatarURL, err := h.userService.UpdateAvatar(c.Request.Context(), userID.(string), file)
	if err != nil {
		switch err {
		case model.ErrInvalidFileType:
			utils.SendError(c, 400, "Invalid file type", err)
		case model.ErrFileTooLarge:
			utils.SendError(c, 400, "File too large", err)
		case model.ErrUserNotFound:
			utils.SendError(c, 404, "User not found", err)
		default:
			utils.SendError(c, 500, "Failed to update avatar", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Avatar updated successfully", gin.H{
		"avatar_url": avatarURL,
	})
}
