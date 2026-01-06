package model

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"omitempty"`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`    // Can be email or username
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest represents the profile update request
type UpdateProfileRequest struct {
	FullName    *string `json:"full_name" binding:"omitempty"`
	Phone       *string `json:"phone" binding:"omitempty"`
	Address     *string `json:"address" binding:"omitempty"`
	DateOfBirth *string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	Bio         *string `json:"bio" binding:"omitempty,max=500"`
}

// AssignRolesRequest represents the role assignment request
type AssignRolesRequest struct {
	RoleIDs []int64 `json:"role_ids" binding:"required,min=1"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// GetDefaults returns default pagination values if not set
func (p *PaginationRequest) GetDefaults() (page, limit int) {
	page = p.Page
	if page == 0 {
		page = 1
	}

	limit = p.Limit
	if limit == 0 {
		limit = 10
	}

	return page, limit
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	page, limit := p.GetDefaults()
	return (page - 1) * limit
}
