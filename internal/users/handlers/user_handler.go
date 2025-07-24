package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
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
// @Description Register a new user in the system with provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param user body models.UserCreate true "User creation details"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Router /users [post]
// @Security ApiKeyAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Create a user model to bind request body
	var userCreate models.UserCreate

	// Bind and validate input
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
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
		response.BadRequest(c, "Missing required fields", map[string]interface{}{
			"email":    user.Email == "",
			"password": user.Password == "",
			"role":     user.Role == "",
		})
		return
	}

	// Create user through service
	if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
		response.Conflict(c, "Failed to create user", err.Error())
		return
	}

	// Respond with created user (excluding sensitive info)
	response.SuccessCreated(c, user, "User created successfully")
}

// ListUsers godoc
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param limit query int false "Number of users to retrieve" default(10)
// @Param offset query int false "Starting offset for pagination" default(0)
// @Success 200 {object} models.User "Successfully retrieved users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve users"
// @Router /users [get]
// @Security ApiKeyAuth
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters with defaults
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve users
	users, err := h.userService.GetUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users", err.Error())
		return
	}

	// Use WithPagination for consistent pagination response
	response.WithPagination(c, users, len(users), offset/limit+1, limit)
}

// SearchUser godoc
// @Summary Search user by ID or email
// @Description Retrieve a specific user's details by their unique identifier or email address
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id query string false "User ID"
// @Param email query string false "User Email"
// @Success 200 {object} models.User "User details retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/search [get]
// @Security ApiKeyAuth
func (h *UserHandler) SearchUser(c *gin.Context) {
	// Get search parameters
	userID := c.Query("id")
	email := c.Query("email")

	// Validate search parameters
	if userID == "" && email == "" {
		response.BadRequest(c, "Either user ID or email is required", nil)
		return
	}

	var user *models.User
	var err error

	// Search by ID if provided
	if userID != "" {
		user, err = h.userService.GetUserByID(c.Request.Context(), userID)
	} else {
		// Otherwise, search by email
		user, err = h.userService.GetUserByEmail(c.Request.Context(), email)
	}

	// Handle search errors
	if err != nil {
		response.NotFound(c, "User not found", err.Error())
		return
	}

	response.SuccessOK(c, user, "User retrieved successfully")
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user's details
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id path string true "User ID"
// @Param user body models.UserUpdate true "User update details"
// @Success 200 {object} map[string]string "User updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 409 {object} response.ErrorResponse "Failed to update user"
// @Router /users/{id} [put]
// @Security ApiKeyAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	// Create a user model to bind request body
	var updateUser models.User

	// Bind input
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set the ID from path parameter
	updateUser.ID = userID

	// Update user
	if err := h.userService.UpdateUser(c.Request.Context(), &updateUser); err != nil {
		response.Conflict(c, "Failed to update user", err.Error())
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft delete a user by their unique identifier
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 409 {object} response.ErrorResponse "Failed to delete user"
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Conflict(c, "Failed to delete user", err.Error())
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User deleted successfully")
}
