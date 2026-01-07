package service

import (
	"context"
	"fmt"
	"time"

	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/repository"
	"xnoia-go-boilerplate/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	GetCurrentUser(ctx context.Context, userID string) (*model.UserResponse, error)
}

type authService struct {
	repo       repository.UserRepository
	jwtSecret  string
	jwtExp     time.Duration
	refreshExp time.Duration
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, jwtExp, refreshExp time.Duration) AuthService {
	return &authService{
		repo:       repo,
		jwtSecret:  jwtSecret,
		jwtExp:     jwtExp,
		refreshExp: refreshExp,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Validate username format
	if !utils.IsValidUsername(req.Username) {
		return nil, model.ErrInvalidUsername
	}

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return nil, model.ErrInvalidEmail
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, model.ErrEmailAlreadyExists
	} else if err != model.ErrUserNotFound {
		return nil, err
	}

	// Check if username already exists
	_, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, model.ErrUsernameAlreadyExists
	} else if err != model.ErrUserNotFound {
		return nil, err
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, req.Email, req.Username, passwordHash, true, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create user details
	_, err = s.repo.CreateUserDetails(ctx, user.ID, req.FullName, "", "", nil, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create user details: %w", err)
	}

	// Assign default 'user' role
	defaultRole, err := s.repo.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to get default role: %w", err)
	}

	err = s.repo.AssignRoleToUser(ctx, user.ID, defaultRole.ID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	// Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build role names for JWT
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.RoleName
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Username, roleNames, s.jwtSecret, s.jwtExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret, s.refreshExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Build response
	roleResponses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = model.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description.String,
		}
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExp.Seconds()),
		User: model.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Username:      user.Username,
			EmailVerified: user.EmailVerified.Bool,
			IsActive:      user.IsActive.Bool,
			Roles:         roleResponses,
			CreatedAt:     user.CreatedAt.Time,
		},
	}, nil
}

// Login authenticates a user with email/username and password
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Find user by email or username
	user, err := s.repo.GetUserByEmailOrUsername(ctx, req.Login)
	if err != nil {
		if err == model.ErrUserNotFound {
			return nil, model.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive.Bool {
		return nil, model.ErrUserInactive
	}

	// Update last login
	err = s.repo.UpdateUserLastLogin(ctx, user.ID, time.Now())
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build role names for JWT
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.RoleName
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Username, roleNames, s.jwtSecret, s.jwtExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret, s.refreshExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Build response
	roleResponses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = model.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description.String,
		}
	}

	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExp.Seconds()),
		User: model.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Username:      user.Username,
			EmailVerified: user.EmailVerified.Bool,
			IsActive:      user.IsActive.Bool,
			LastLogin:     lastLogin,
			Roles:         roleResponses,
			CreatedAt:     user.CreatedAt.Time,
		},
	}, nil
}

// RefreshToken generates new access and refresh tokens
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	// Validate refresh token
	userID, err := utils.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if err == model.ErrUserNotFound {
			return nil, model.ErrInvalidToken
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive.Bool {
		return nil, model.ErrUserInactive
	}

	// Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build role names for JWT
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.RoleName
	}

	// Generate new tokens
	newAccessToken, err := utils.GenerateToken(user.ID, user.Email, user.Username, roleNames, s.jwtSecret, s.jwtExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret, s.refreshExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Build response
	roleResponses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = model.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description.String,
		}
	}

	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	return &model.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExp.Seconds()),
		User: model.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Username:      user.Username,
			EmailVerified: user.EmailVerified.Bool,
			IsActive:      user.IsActive.Bool,
			LastLogin:     lastLogin,
			Roles:         roleResponses,
			CreatedAt:     user.CreatedAt.Time,
		},
	}, nil
}

// GetCurrentUser retrieves the current user information
func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*model.UserResponse, error) {
	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build role responses
	roleResponses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = model.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description.String,
		}
	}

	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	return &model.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		EmailVerified: user.EmailVerified.Bool,
		IsActive:      user.IsActive.Bool,
		LastLogin:     lastLogin,
		Roles:         roleResponses,
		CreatedAt:     user.CreatedAt.Time,
	}, nil
}
