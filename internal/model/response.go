package model

import "time"

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"` // Seconds until access token expires
	User         UserResponse `json:"user"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID            int64          `json:"id"`
	Email         string         `json:"email"`
	Username      string         `json:"username"`
	EmailVerified bool           `json:"email_verified"`
	IsActive      bool           `json:"is_active"`
	LastLogin     *time.Time     `json:"last_login,omitempty"`
	Roles         []RoleResponse `json:"roles"`
	CreatedAt     time.Time      `json:"created_at"`
}

// RoleResponse represents the role data in responses
type RoleResponse struct {
	ID          int64  `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description,omitempty"`
}

// ProfileResponse represents the complete user profile
type ProfileResponse struct {
	User    UserResponse       `json:"user"`
	Details UserDetailResponse `json:"details"`
}

// UserDetailResponse represents the user profile details
type UserDetailResponse struct {
	FullName    *string `json:"full_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Address     *string `json:"address,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

// UserListResponse represents a user in the admin user list
type UserListResponse struct {
	ID            int64          `json:"id"`
	Email         string         `json:"email"`
	Username      string         `json:"username"`
	EmailVerified bool           `json:"email_verified"`
	IsActive      bool           `json:"is_active"`
	LastLogin     *time.Time     `json:"last_login,omitempty"`
	Roles         []RoleResponse `json:"roles"`
	CreatedAt     time.Time      `json:"created_at"`
}
