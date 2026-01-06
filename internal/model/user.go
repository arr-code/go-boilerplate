package model

import (
	"database/sql"
	"time"
)

// User represents the users table
type User struct {
	ID            int64        `json:"id"`
	Email         string       `json:"email"`
	Username      string       `json:"username"`
	PasswordHash  string       `json:"-"` // Never expose password hash in JSON
	IsActive      bool         `json:"is_active"`
	EmailVerified bool         `json:"email_verified"`
	LastLogin     sql.NullTime `json:"last_login,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// UserDetail represents the user_details table
type UserDetail struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	FullName    sql.NullString `json:"full_name,omitempty"`
	Phone       sql.NullString `json:"phone,omitempty"`
	Address     sql.NullString `json:"address,omitempty"`
	DateOfBirth sql.NullTime   `json:"date_of_birth,omitempty"`
	AvatarURL   sql.NullString `json:"avatar_url,omitempty"`
	Bio         sql.NullString `json:"bio,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Role represents the roles table
type Role struct {
	ID          int64     `json:"id"`
	RoleName    string    `json:"role_name"`
	Description string    `json:"description,omitempty"`
	Permissions []byte    `json:"permissions,omitempty"` // JSONB stored as bytes
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole represents the user_roles junction table
type UserRole struct {
	ID         int64        `json:"id"`
	UserID     int64        `json:"user_id"`
	RoleID     int64        `json:"role_id"`
	AssignedAt time.Time    `json:"assigned_at"`
	AssignedBy sql.NullInt64 `json:"assigned_by,omitempty"`
}

// UserWithRoles combines user information with their roles
type UserWithRoles struct {
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

// UserWithDetails combines user information with their profile details
type UserWithDetails struct {
	User    User       `json:"user"`
	Details UserDetail `json:"details"`
}

// UserProfile combines user, details, and roles
type UserProfile struct {
	User    User       `json:"user"`
	Details UserDetail `json:"details"`
	Roles   []Role     `json:"roles"`
}
