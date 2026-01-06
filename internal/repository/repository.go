package repository

import (
	"context"
	"database/sql"
	"time"

	"xnoia-go-boilerplate/db/sqlc"
)

// UserRepository defines all database operations for users
type UserRepository interface {
	// User operations
	CreateUser(ctx context.Context, email, username, passwordHash string, isActive, emailVerified bool) (*sqlc.User, error)
	GetUserByID(ctx context.Context, id int64) (*sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error)
	GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error)
	GetUserByEmailOrUsername(ctx context.Context, login string) (*sqlc.User, error)
	UpdateUserLastLogin(ctx context.Context, id int64, loginTime time.Time) error
	UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error
	DeactivateUser(ctx context.Context, id int64) error
	ActivateUser(ctx context.Context, id int64) error
	VerifyUserEmail(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int) ([]sqlc.User, error)
	CountUsers(ctx context.Context) (int64, error)

	// User details operations
	CreateUserDetails(ctx context.Context, userID int64, fullName, phone, address string, dateOfBirth *time.Time, avatarURL, bio string) (*sqlc.UserDetail, error)
	GetUserDetailsByUserID(ctx context.Context, userID int64) (*sqlc.UserDetail, error)
	UpdateUserDetails(ctx context.Context, userID int64, fullName, phone, address string, dateOfBirth *time.Time, bio string) error
	UpdateUserAvatar(ctx context.Context, userID int64, avatarURL string) error

	// Role operations
	GetRoleByID(ctx context.Context, id int64) (*sqlc.Role, error)
	GetRoleByName(ctx context.Context, name string) (*sqlc.Role, error)
	ListRoles(ctx context.Context) ([]sqlc.Role, error)
	CreateRole(ctx context.Context, roleName, description string, permissions []byte) (*sqlc.Role, error)

	// User-role operations
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy int64) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error
	RemoveAllUserRoles(ctx context.Context, userID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]sqlc.Role, error)
	CheckUserHasRole(ctx context.Context, userID int64, roleName string) (bool, error)
	CheckUserHasAnyRole(ctx context.Context, userID int64, roleNames []string) (bool, error)

	// Transaction support
	WithTx(tx *sql.Tx) UserRepository
}
