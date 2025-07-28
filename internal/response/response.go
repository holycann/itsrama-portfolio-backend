package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// APIResponse is the standard structure for all API responses
type APIResponse struct {
	// Status of the response (success/error)
	Success bool `json:"success"`

	// Unique request identifier for tracing
	RequestID uuid.UUID `json:"request_id"`

	// Human-readable message
	Message string `json:"message,omitempty"`

	// Detailed error information (only populated for error responses)
	Error *ErrorDetails `json:"error,omitempty"`

	// Pagination information (optional)
	Pagination *Pagination `json:"pagination,omitempty"`

	// Actual response data
	Data interface{} `json:"data,omitempty"`
}

// ErrorDetails provides structured error information
type ErrorDetails struct {
	// Machine-readable error code
	Code string `json:"code,omitempty"`

	// Detailed error description
	Details string `json:"details,omitempty"`
}

// Pagination represents standard pagination metadata
type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
}

// Success creates a successful API response
func Success(c *gin.Context, statusCode int, data interface{}, message string, pagination ...*Pagination) {
	resp := APIResponse{
		Success:   true,
		RequestID: uuid.New(),
		Message:   message,
		Data:      data,
	}

	// Add pagination if provided
	if len(pagination) > 0 && pagination[0] != nil {
		resp.Pagination = pagination[0]
	}

	c.JSON(statusCode, resp)
}

// Error creates a standardized error response
func Error(c *gin.Context, statusCode int, errorCode string, message string, details string) {
	resp := APIResponse{
		Success:   false,
		RequestID: uuid.New(),
		Error: &ErrorDetails{
			Code:    errorCode,
			Details: details,
		},
		Message: message,
	}

	c.JSON(statusCode, resp)
}

// SuccessCreated is a shorthand for successful creation responses
func SuccessCreated(c *gin.Context, data interface{}, message string) {
	Success(c, http.StatusCreated, data, message)
}

// SuccessOK is a shorthand for successful OK responses
func SuccessOK(c *gin.Context, data interface{}, message string, pagination ...*Pagination) {
	Success(c, http.StatusOK, data, message, pagination...)
}

// BadRequest generates a 400 Bad Request error response
func BadRequest(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusBadRequest, errorCode, message, details)
}

// Unauthorized generates a 401 Unauthorized error response
func Unauthorized(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusUnauthorized, errorCode, message, details)
}

// Forbidden generates a 403 Forbidden error response
func Forbidden(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusForbidden, errorCode, message, details)
}

// NotFound generates a 404 Not Found error response
func NotFound(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusNotFound, errorCode, message, details)
}

// Conflict generates a 409 Conflict error response
func Conflict(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusConflict, errorCode, message, details)
}

// InternalServerError generates a 500 Internal Server Error response
func InternalServerError(c *gin.Context, errorCode string, message string, details string) {
	Error(c, http.StatusInternalServerError, errorCode, message, details)
}
