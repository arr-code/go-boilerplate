package handler

import (
	"strconv"

	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/service"
	"xnoia-go-boilerplate/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ListUsers retrieves a paginated list of all users
// GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var pagination model.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Get defaults
	page, limit := pagination.GetDefaults()

	users, total, err := h.adminService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		utils.SendError(c, 500, "Failed to list users", err)
		return
	}

	utils.SendPaginated(c, users, page, limit, total)
}

// AssignRoles assigns roles to a user
// PUT /api/v1/admin/users/:id/roles
func (h *AdminHandler) AssignRoles(c *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		utils.SendError(c, 400, "Invalid user ID", err)
		return
	}

	// Get admin user ID from context
	adminID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "Unauthorized", model.ErrUnauthorized)
		return
	}

	var req model.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	err = h.adminService.AssignRoles(c.Request.Context(), userID, req.RoleIDs, adminID.(int64))
	if err != nil {
		if err == model.ErrUserNotFound {
			utils.SendError(c, 404, "User not found", err)
		} else if err == model.ErrRoleNotFound {
			utils.SendError(c, 404, "Role not found", err)
		} else {
			utils.SendError(c, 500, "Failed to assign roles", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "Roles assigned successfully", nil)
}

// DeactivateUser deactivates a user account
// DELETE /api/v1/admin/users/:id
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		utils.SendError(c, 400, "Invalid user ID", err)
		return
	}

	err = h.adminService.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		if err == model.ErrUserNotFound {
			utils.SendError(c, 404, "User not found", err)
		} else {
			utils.SendError(c, 500, "Failed to deactivate user", err)
		}
		return
	}

	utils.SendSuccess(c, 200, "User deactivated successfully", nil)
}
