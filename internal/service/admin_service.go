package service

import (
	"context"
	"fmt"
	"time"

	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/repository"
)

type AdminService interface {
	ListUsers(ctx context.Context, page, limit int) ([]model.UserListResponse, int64, error)
	AssignRoles(ctx context.Context, userID int64, roleIDs []int64, assignedBy int64) error
	DeactivateUser(ctx context.Context, userID int64) error
}

type adminService struct {
	repo repository.UserRepository
}

func NewAdminService(repo repository.UserRepository) AdminService {
	return &adminService{
		repo: repo,
	}
}

// ListUsers retrieves a paginated list of users with their roles
func (s *adminService) ListUsers(ctx context.Context, page, limit int) ([]model.UserListResponse, int64, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get users
	users, err := s.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Build response with roles for each user
	userResponses := make([]model.UserListResponse, len(users))
	for i, user := range users {
		// Get user roles
		roles, err := s.repo.GetUserRoles(ctx, int64(user.ID))
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get roles for user %d: %w", user.ID, err)
		}

		// Build role responses
		roleResponses := make([]model.RoleResponse, len(roles))
		for j, role := range roles {
			roleResponses[j] = model.RoleResponse{
				ID:          int64(role.ID),
				RoleName:    role.RoleName,
				Description: role.Description.String,
			}
		}

		var lastLogin *time.Time
		if user.LastLogin.Valid {
			lastLogin = &user.LastLogin.Time
		}

		userResponses[i] = model.UserListResponse{
			ID:            int64(user.ID),
			Email:         user.Email,
			Username:      user.Username,
			EmailVerified: user.EmailVerified.Bool,
			IsActive:      user.IsActive.Bool,
			LastLogin:     lastLogin,
			Roles:         roleResponses,
			CreatedAt:     user.CreatedAt.Time,
		}
	}

	return userResponses, total, nil
}

// AssignRoles assigns roles to a user (removes existing roles first)
func (s *adminService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64, assignedBy int64) error {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Validate that all roles exist
	for _, roleID := range roleIDs {
		_, err := s.repo.GetRoleByID(ctx, roleID)
		if err != nil {
			if err == model.ErrRoleNotFound {
				return fmt.Errorf("role with ID %d not found", roleID)
			}
			return err
		}
	}

	// Remove all existing roles
	err = s.repo.RemoveAllUserRoles(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to remove existing roles: %w", err)
	}

	// Assign new roles
	for _, roleID := range roleIDs {
		err := s.repo.AssignRoleToUser(ctx, userID, roleID, assignedBy)
		if err != nil {
			return fmt.Errorf("failed to assign role %d: %w", roleID, err)
		}
	}

	return nil
}

// DeactivateUser marks a user as inactive
func (s *adminService) DeactivateUser(ctx context.Context, userID int64) error {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Deactivate user
	err = s.repo.DeactivateUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}
