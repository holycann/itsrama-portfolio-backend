package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// UserHandler handles user-related HTTP requests
// @Description Manages user-related operations such as creation, retrieval, update, and deletion
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new instance of UserHandler
// @Description Initializes a new UserHandler with the provided UserService
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user body models.UserCreate true "User Registration Details"
// @Success 201 {object} response.APIResponse{data=models.User} "User created successfully"
// @Failure 400 {object} response.APIResponse "Invalid user creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Create a user model to bind request body
	var userCreate models.UserCreate

	// Bind and validate input
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error(), "")
		return
	}

	// Convert UserCreate to User model
	user := models.User{
		Email:    userCreate.Email,
		Password: userCreate.Password,
		Role:     userCreate.Role,
	}

	// Validate email
	if user.Email == "" || user.Password == "" || user.Role == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"email":    user.Email == "",
			"password": user.Password == "",
			"role":     user.Role == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	// Create user through service
	if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
		response.Conflict(c, "Failed to create user", err.Error(), "")
		return
	}

	// Respond with created user (excluding sensitive info)
	response.SuccessCreated(c, user, "User created successfully")
}

// ListUsers godoc
// @Summary List users
// @Description Retrieve a list of users with pagination and filtering
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of users to retrieve" default(10)
// @Param offset query int false "Number of users to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.User} "Users retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list users"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters with defaults
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Prepare list options
	listOptions := repository.ListOptions{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: repository.SortDescending,
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Optional filtering
	filters := []repository.FilterOption{}
	if role := c.Query("role"); role != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "role",
			Operator: "=",
			Value:    role,
		})
	}
	listOptions.Filters = filters

	// Retrieve users
	users, err := h.userService.ListUsers(c.Request.Context(), listOptions)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users", err.Error(), "")
		return
	}

	// Count total users for pagination
	totalUsers, err := h.userService.CountUsers(c.Request.Context(), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to count users", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalUsers,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalUsers) / float64(limit))),
		HasNextPage: offset+limit < totalUsers,
	}

	// Respond with users and pagination
	response.SuccessOK(c, users, "Users retrieved successfully", pagination)
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by various criteria
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (email, name, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.APIResponse{data=[]models.User} "Users found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	// Get search query
	query := c.Query("query")
	if query == "" {
		response.BadRequest(c, "Search query is required", "Empty search query", "")
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Prepare list options for search
	listOptions := repository.ListOptions{
		Limit:  limit,
		Offset: offset,
		Filters: []repository.FilterOption{
			{
				Field:    "email",
				Operator: "like",
				Value:    query,
			},
		},
	}

	// Search users
	users, err := h.userService.SearchUsers(c.Request.Context(), query, listOptions)
	if err != nil {
		response.InternalServerError(c, "Failed to search users", err.Error(), "")
		return
	}

	// Count total search results
	totalUsers, err := h.userService.CountUsers(c.Request.Context(), listOptions.Filters)
	if err != nil {
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalUsers,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalUsers) / float64(limit))),
		HasNextPage: offset+limit < totalUsers,
	}

	// Respond with users and pagination
	response.SuccessOK(c, users, "Users found successfully", pagination)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's details
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "User ID"
// @Param user body models.UserUpdate true "User Update Details"
// @Success 200 {object} response.APIResponse{data=models.User} "User updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid user update details"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", "Missing user ID", "")
		return
	}

	// Create a user model to bind request body
	var updateUser models.User

	// Bind input
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	updateUser.ID = userID

	// Update user
	if err := h.userService.UpdateUser(c.Request.Context(), &updateUser); err != nil {
		response.Conflict(c, "Failed to update user", err.Error(), "")
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user from the system by their unique identifier
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", "Missing user ID", "")
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Conflict(c, "Failed to delete user", err.Error(), "")
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User deleted successfully")
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user's details by their unique identifier
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse{data=models.User} "User retrieved successfully"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required", "Missing user ID", "")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "User not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, user, "User retrieved successfully")
}
