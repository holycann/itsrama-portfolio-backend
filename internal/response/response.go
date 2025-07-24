package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response structure
type Response struct {
	// Indicates whether the request was successful
	Success bool `json:"success"`

	// A human-readable message describing the result
	Message string `json:"message,omitempty"`

	// Metadata or additional information about the response
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// The actual data payload (can be an object, array, or null)
	Data interface{} `json:"data"`
}

// ErrorResponse represents a standard error response structure
type ErrorResponse struct {
	// Indicates whether the request was successful (always false for errors)
	Success bool `json:"success"`

	// The primary error message
	Error string `json:"error"`

	// Detailed error information
	Details interface{} `json:"details,omitempty"`
}

// Success creates a successful response with optional data and message
func Success(c *gin.Context, statusCode int, data interface{}, message string, metadata ...map[string]interface{}) {
	resp := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	// Add metadata if provided
	if len(metadata) > 0 {
		resp.Metadata = metadata[0]
	}

	c.JSON(statusCode, resp)
}

// SuccessCreated is a shorthand for successful creation responses
func SuccessCreated(c *gin.Context, data interface{}, message string, metadata ...map[string]interface{}) {
	Success(c, http.StatusCreated, data, message, metadata...)
}

// SuccessOK is a shorthand for successful OK responses
func SuccessOK(c *gin.Context, data interface{}, message string, metadata ...map[string]interface{}) {
	Success(c, http.StatusOK, data, message, metadata...)
}

// Error creates a standard error response
func Error(c *gin.Context, statusCode int, errorMessage string, details interface{}) {
	resp := ErrorResponse{
		Success: false,
		Error:   errorMessage,
		Details: details,
	}

	c.JSON(statusCode, resp)
}

// BadRequest generates a 400 Bad Request error response
func BadRequest(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusBadRequest, errorMessage, details)
}

// Unauthorized generates a 401 Unauthorized error response
func Unauthorized(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusUnauthorized, errorMessage, details)
}

// Forbidden generates a 403 Forbidden error response
func Forbidden(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusForbidden, errorMessage, details)
}

// NotFound generates a 404 Not Found error response
func NotFound(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusNotFound, errorMessage, details)
}

// Conflict generates a 409 Conflict error response
func Conflict(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusConflict, errorMessage, details)
}

// InternalServerError generates a 500 Internal Server Error response
func InternalServerError(c *gin.Context, errorMessage string, details interface{}) {
	Error(c, http.StatusInternalServerError, errorMessage, details)
}

// Pagination represents standard pagination metadata
type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
}

// WithPagination adds pagination metadata to a successful response
func WithPagination(c *gin.Context, data interface{}, total, page, perPage int) {
	// Calculate total pages
	totalPages := (total + perPage - 1) / perPage
	hasNextPage := page < totalPages

	metadata := map[string]interface{}{
		"pagination": Pagination{
			Total:       total,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			HasNextPage: hasNextPage,
		},
	}

	Success(c, http.StatusOK, data, "Retrieved successfully", metadata)
}