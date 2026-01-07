package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"xnoia-go-boilerplate/db/sqlc"
	"xnoia-go-boilerplate/internal/model"

	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewUserRepository creates a new PostgresUserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// WithTx returns a new repository with a transaction
func (r *PostgresUserRepository) WithTx(tx *sql.Tx) UserRepository {
	return &PostgresUserRepository{
		db:      r.db,
		queries: r.queries.WithTx(tx),
	}
}

// User operations

func (r *PostgresUserRepository) CreateUser(ctx context.Context, email, username, passwordHash string, isActive, emailVerified bool) (*sqlc.User, error) {
	user, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		IsActive:      sql.NullBool{Bool: isActive, Valid: true},
		EmailVerified: sql.NullBool{Bool: emailVerified, Valid: true},
	})
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByEmailOrUsername(ctx context.Context, login string) (*sqlc.User, error) {
	user, err := r.queries.GetUserByEmailOrUsername(ctx, sqlc.GetUserByEmailOrUsernameParams{
		Email:    login,
		Username: login,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) UpdateUserLastLogin(ctx context.Context, id string, loginTime time.Time) error {
	err := r.queries.UpdateUserLastLogin(ctx, sqlc.UpdateUserLastLoginParams{
		ID:        id,
		LastLogin: sql.NullTime{Time: loginTime, Valid: true},
	})
	return err
}

func (r *PostgresUserRepository) UpdateUserPassword(ctx context.Context, id string, passwordHash string) error {
	err := r.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
	return err
}

func (r *PostgresUserRepository) DeactivateUser(ctx context.Context, id string) error {
	return r.queries.DeactivateUser(ctx, id)
}

func (r *PostgresUserRepository) ActivateUser(ctx context.Context, id string) error {
	return r.queries.ActivateUser(ctx, id)
}

func (r *PostgresUserRepository) VerifyUserEmail(ctx context.Context, id string) error {
	return r.queries.VerifyUserEmail(ctx, id)
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]sqlc.User, error) {
	users, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// User details operations

func (r *PostgresUserRepository) CreateUserDetails(ctx context.Context, userID string, fullName, phone, address string, dateOfBirth *time.Time, avatarURL, bio string) (*sqlc.UserDetail, error) {
	var dob sql.NullTime
	if dateOfBirth != nil {
		dob = sql.NullTime{Time: *dateOfBirth, Valid: true}
	}

	details, err := r.queries.CreateUserDetails(ctx, sqlc.CreateUserDetailsParams{
		UserID:      userID,
		FullName:    sql.NullString{String: fullName, Valid: fullName != ""},
		Phone:       sql.NullString{String: phone, Valid: phone != ""},
		Address:     sql.NullString{String: address, Valid: address != ""},
		DateOfBirth: dob,
		AvatarUrl:   sql.NullString{String: avatarURL, Valid: avatarURL != ""},
		Bio:         sql.NullString{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return nil, err
	}
	return &details, nil
}

func (r *PostgresUserRepository) GetUserDetailsByUserID(ctx context.Context, userID string) (*sqlc.UserDetail, error) {
	details, err := r.queries.GetUserDetailsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrRecordNotFound
		}
		return nil, err
	}
	return &details, nil
}

func (r *PostgresUserRepository) UpdateUserDetails(ctx context.Context, userID string, fullName, phone, address string, dateOfBirth *time.Time, bio string) error {
	var dob sql.NullTime
	if dateOfBirth != nil {
		dob = sql.NullTime{Time: *dateOfBirth, Valid: true}
	}

	err := r.queries.UpdateUserDetails(ctx, sqlc.UpdateUserDetailsParams{
		UserID:      userID,
		FullName:    sql.NullString{String: fullName, Valid: true},
		Phone:       sql.NullString{String: phone, Valid: true},
		Address:     sql.NullString{String: address, Valid: true},
		DateOfBirth: dob,
		Bio:         sql.NullString{String: bio, Valid: true},
	})
	return err
}

func (r *PostgresUserRepository) UpdateUserAvatar(ctx context.Context, userID string, avatarURL string) error {
	err := r.queries.UpdateUserAvatar(ctx, sqlc.UpdateUserAvatarParams{
		UserID:    userID,
		AvatarUrl: sql.NullString{String: avatarURL, Valid: true},
	})
	return err
}

// Role operations

func (r *PostgresUserRepository) GetRoleByID(ctx context.Context, id string) (*sqlc.Role, error) {
	role, err := r.queries.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *PostgresUserRepository) GetRoleByName(ctx context.Context, name string) (*sqlc.Role, error) {
	role, err := r.queries.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *PostgresUserRepository) ListRoles(ctx context.Context) ([]sqlc.Role, error) {
	roles, err := r.queries.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *PostgresUserRepository) CreateRole(ctx context.Context, roleName, description string, permissions []byte) (*sqlc.Role, error) {
	role, err := r.queries.CreateRole(ctx, sqlc.CreateRoleParams{
		RoleName:    roleName,
		Description: sql.NullString{String: description, Valid: description != ""},
		Permissions: permissions,
	})
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &role, nil
}

// User-role operations

func (r *PostgresUserRepository) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string) error {
	_, err := r.queries.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: sql.NullString{String: assignedBy, Valid: assignedBy != ""},
	})
	if err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *PostgresUserRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	err := r.queries.RemoveRoleFromUser(ctx, sqlc.RemoveRoleFromUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	return err
}

func (r *PostgresUserRepository) RemoveAllUserRoles(ctx context.Context, userID string) error {
	err := r.queries.RemoveAllUserRoles(ctx, userID)
	return err
}

func (r *PostgresUserRepository) GetUserRoles(ctx context.Context, userID string) ([]sqlc.Role, error) {
	roles, err := r.queries.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *PostgresUserRepository) CheckUserHasRole(ctx context.Context, userID string, roleName string) (bool, error) {
	hasRole, err := r.queries.CheckUserHasRole(ctx, sqlc.CheckUserHasRoleParams{
		UserID:   userID,
		RoleName: roleName,
	})
	if err != nil {
		return false, err
	}
	return hasRole, nil
}

func (r *PostgresUserRepository) CheckUserHasAnyRole(ctx context.Context, userID string, roleNames []string) (bool, error) {
	hasAnyRole, err := r.queries.CheckUserHasAnyRole(ctx, sqlc.CheckUserHasAnyRoleParams{
		UserID:  userID,
		Column2: roleNames,
	})
	if err != nil {
		return false, err
	}
	return hasAnyRole, nil
}

// Helper function to map database errors to custom errors
func mapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Check for PostgreSQL unique violation
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505": // unique_violation
			if pqErr.Constraint == "users_email_key" {
				return model.ErrEmailAlreadyExists
			}
			if pqErr.Constraint == "users_username_key" {
				return model.ErrUsernameAlreadyExists
			}
			return model.ErrUserAlreadyExists
		case "23503": // foreign_key_violation
			return errors.New("foreign key constraint violation")
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrRecordNotFound
	}

	return err
}
