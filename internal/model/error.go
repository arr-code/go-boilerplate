package model

import "errors"

// Authentication and Authorization Errors
var (
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden: insufficient permissions")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrInvalidCredentials = errors.New("invalid email/username or password")
)

// User Errors
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrUserNotVerified     = errors.New("user email is not verified")
)

// Role Errors
var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrInvalidRole       = errors.New("invalid role")
)

// Validation Errors
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidUsername = errors.New("invalid username format")
	ErrWeakPassword    = errors.New("password does not meet strength requirements")
	ErrInvalidInput    = errors.New("invalid input data")
)

// Database Errors
var (
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
	ErrRecordNotFound     = errors.New("record not found")
)

// File Upload Errors
var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds limit")
	ErrFileUploadFailed = errors.New("file upload failed")
)
