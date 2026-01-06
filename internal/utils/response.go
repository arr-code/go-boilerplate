package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PaginationResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SendSuccess sends a successful JSON response
func SendSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error JSON response
func SendError(c *gin.Context, code int, message string, err error) {
	response := ErrorResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(code, response)
}

// SendValidationError sends a validation error response with detailed field errors
func SendValidationError(c *gin.Context, err error) {
	var errorMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, formatValidationError(e))
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	c.JSON(400, gin.H{
		"success": false,
		"message": "Validation failed",
		"errors":  errorMessages,
	})
}

// SendPaginated sends a paginated response
func SendPaginated(c *gin.Context, data interface{}, page, limit int, total int64) {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(200, PaginationResponse{
		Success: true,
		Message: "Data retrieved successfully",
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// formatValidationError formats a validation error into a readable message
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email address"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "datetime":
		return e.Field() + " must be a valid date in format " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
