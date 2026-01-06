package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"xnoia-go-boilerplate/internal/model"
	"xnoia-go-boilerplate/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int64) (*model.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID int64, req *model.UpdateProfileRequest) error
	UpdateAvatar(ctx context.Context, userID int64, fileHeader *multipart.FileHeader) (string, error)
}

type userService struct {
	repo repository.UserRepository
	// s3Client for future implementation
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// GetProfile retrieves the user profile with details and roles
func (s *userService) GetProfile(ctx context.Context, userID int64) (*model.ProfileResponse, error) {
	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user details
	details, err := s.repo.GetUserDetailsByUserID(ctx, userID)
	if err != nil && err != model.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// Get user roles
	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build role responses
	roleResponses := make([]model.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = model.RoleResponse{
			ID:          int64(role.ID),
			RoleName:    role.RoleName,
			Description: role.Description.String,
		}
	}

	// Build user response
	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	userResponse := model.UserResponse{
		ID:            int64(user.ID),
		Email:         user.Email,
		Username:      user.Username,
		EmailVerified: user.EmailVerified.Bool,
		IsActive:      user.IsActive.Bool,
		LastLogin:     lastLogin,
		Roles:         roleResponses,
		CreatedAt:     user.CreatedAt.Time,
	}

	// Build details response
	detailsResponse := model.UserDetailResponse{}
	if details != nil {
		if details.FullName.Valid {
			detailsResponse.FullName = &details.FullName.String
		}
		if details.Phone.Valid {
			detailsResponse.Phone = &details.Phone.String
		}
		if details.Address.Valid {
			detailsResponse.Address = &details.Address.String
		}
		if details.DateOfBirth.Valid {
			dobStr := details.DateOfBirth.Time.Format("2006-01-02")
			detailsResponse.DateOfBirth = &dobStr
		}
		if details.AvatarUrl.Valid {
			detailsResponse.AvatarURL = &details.AvatarUrl.String
		}
		if details.Bio.Valid {
			detailsResponse.Bio = &details.Bio.String
		}
	}

	return &model.ProfileResponse{
		User:    userResponse,
		Details: detailsResponse,
	}, nil
}

// UpdateProfile updates the user profile details
func (s *userService) UpdateProfile(ctx context.Context, userID int64, req *model.UpdateProfileRequest) error {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user details exist
	_, err = s.repo.GetUserDetailsByUserID(ctx, userID)
	if err == model.ErrRecordNotFound {
		// Create user details if they don't exist
		_, err = s.repo.CreateUserDetails(ctx, userID, "", "", "", nil, "", "")
		if err != nil {
			return fmt.Errorf("failed to create user details: %w", err)
		}
	} else if err != nil {
		return err
	}

	// Prepare values
	fullName := ""
	phone := ""
	address := ""
	bio := ""
	var dateOfBirth *time.Time

	if req.FullName != nil {
		fullName = *req.FullName
	}
	if req.Phone != nil {
		phone = *req.Phone
	}
	if req.Address != nil {
		address = *req.Address
	}
	if req.Bio != nil {
		bio = *req.Bio
	}
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		dateOfBirth = &dob
	}

	// Update user details
	err = s.repo.UpdateUserDetails(ctx, userID, fullName, phone, address, dateOfBirth, bio)
	if err != nil {
		return fmt.Errorf("failed to update user details: %w", err)
	}

	return nil
}

// UpdateAvatar updates the user avatar
func (s *userService) UpdateAvatar(ctx context.Context, userID int64, fileHeader *multipart.FileHeader) (string, error) {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		return "", model.ErrInvalidFileType
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		return "", model.ErrFileTooLarge
	}

	// TODO: Upload to S3/Cloudflare R2
	// For now, just return a placeholder URL
	// In production, you would:
	// 1. Open the file
	// 2. Upload to S3/R2
	// 3. Get the public URL
	// 4. Save to database

	avatarURL := fmt.Sprintf("/avatars/user-%d-%d.jpg", userID, time.Now().Unix())

	// Update avatar in database
	err = s.repo.UpdateUserAvatar(ctx, userID, avatarURL)
	if err != nil {
		return "", fmt.Errorf("failed to update avatar: %w", err)
	}

	return avatarURL, nil
}
