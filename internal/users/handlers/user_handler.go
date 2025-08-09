package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
	"github.com/holycann/cultour-backend/pkg/validator"
)

// UserHandler handles user-related HTTP requests
// @Description Manages user-related operations such as creation, retrieval, update, and deletion
type UserHandler struct {
	base.BaseHandler
	userService services.UserService
}

// NewUserHandler creates a new instance of UserHandler
// @Description Initializes a new UserHandler with the provided UserService and logger
func NewUserHandler(userService services.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: *base.NewBaseHandler(logger),
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user account
// @Description Allows user registration with email or third-party authentication
// @Description Supports creating user accounts with various authentication providers
// @Tags User Management
// @Accept json
// @Produce json
// @Param user body models.UserCreate true "User Account Creation Details"
// @Success 201 {object} response.APIResponse{data=models.UserDTO} "User account successfully created"
// @Failure 400 {object} response.APIResponse "Invalid user creation payload or validation error"
// @Failure 409 {object} response.APIResponse "User already exists with the provided email"
// @Failure 500 {object} response.APIResponse "Internal server error during user account creation"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Create a user model to bind request body
	var userCreate models.UserCreate

	// Validate request body
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	// Validate user creation payload
	if err := validator.ValidateStruct(userCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Create user through service
	createdUser, err := h.userService.CreateUser(c.Request.Context(), &userCreate)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to create user"))
		return
	}

	// Respond with created user
	h.HandleCreated(c, createdUser, "User created successfully")
}

// ListUsers godoc
// @Summary Retrieve users list
// @Description Fetches a paginated list of user accounts with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags User Management
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of users per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort users by" default("created_at)" Enum(created_at,email)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param role query string false "Filter users by system role" Enum(user,admin,moderator)
// @Param status query string false "Filter users by account status" Enum(active,inactive,suspended)
// @Success 200 {object} response.APIResponse{data=[]models.UserDTO} "Successfully retrieved users list"
// @Success 204 {object} response.APIResponse "No users found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 500 {object} response.APIResponse "Internal server error during users retrieval"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters using base.ParsePaginationParams
	listOptions, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid pagination parameters"))
		return
	}

	// Optional filtering
	if role := c.Query("role"); role != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "role",
			Operator: base.OperatorEqual,
			Value:    role,
		})
	}

	// Retrieve users
	users, totalUsers, err := h.userService.ListUsers(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to retrieve users"))
		return
	}

	// Respond with users and pagination
	h.HandlePagination(c, users, totalUsers, listOptions)
}

// SearchUsers godoc
// @Summary Search users
// @Description Performs a full-text search across user details with advanced filtering
// @Description Allows finding users by keywords, email, and other attributes
// @Tags User Management
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding users" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,email,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param role query string false "Filter users by system role" Enum(user,admin,moderator)
// @Param status query string false "Filter users by account status" Enum(active,inactive,suspended)
// @Success 200 {object} response.APIResponse{data=[]models.UserDTO} "Successfully completed user search"
// @Success 204 {object} response.APIResponse "No users match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 500 {object} response.APIResponse "Internal server error during user search"
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	// Get search query
	query := c.Query("query")
	if query == "" {
		h.HandleError(c, errors.New(errors.ErrValidation, "Search query is required", nil))
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Prepare list options for search
	listOptions := base.ListOptions{
		Page:    offset/limit + 1,
		PerPage: limit,
		Search:  query,
	}

	// Search users
	users, totalUsers, err := h.userService.SearchUsers(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to search users"))
		return
	}

	// Respond with users and pagination
	h.HandlePagination(c, users, totalUsers, listOptions)
}

// UpdateUser godoc
// @Summary Update an existing user account
// @Description Allows administrators to modify user account details
// @Description Supports partial updates with optional fields
// @Tags User Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param user body models.UserUpdate true "User Account Update Payload"
// @Success 200 {object} response.APIResponse{data=models.UserDTO} "User account successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid user update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "User account not found"
// @Failure 500 {object} response.APIResponse "Internal server error during user account update"
// @Router /users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Get user ID from path parameter
	userIDStr := c.Param("id")
	userID, err := h.ValidateUUID(userIDStr, "user_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user ID"))
		return
	}

	// Create a user update model to bind request body
	var updateUser models.UserUpdate

	// Bind and validate input
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	// Validate user update payload
	if err := validator.ValidateStruct(updateUser); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Set the ID from path parameter
	updateUser.ID = userID

	// Update user
	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), &updateUser)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update user"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, updatedUser, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete a user account
// @Description Allows administrators to permanently remove a user account
// @Description Deletes the user account and associated resources
// @Tags User Management
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique User Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "User account successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid user ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "User account not found"
// @Failure 500 {object} response.APIResponse "Internal server error during user account deletion"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from path parameter
	userIDStr := c.Param("id")
	userID, err := h.ValidateUUID(userIDStr, "user_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user ID"))
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(c.Request.Context(), userID.String()); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to delete user"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, gin.H{
		"id": userID,
	}, "User deleted successfully")
}

// GetUserByID godoc
// @Summary Retrieve a specific user account
// @Description Fetches comprehensive details of a user account by its unique identifier
// @Description Returns full user information including profile and role details
// @Tags User Management
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique User Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.UserDTO} "Successfully retrieved user account details"
// @Failure 400 {object} response.APIResponse "Invalid user ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "User account not found"
// @Failure 500 {object} response.APIResponse "Internal server error during user account retrieval"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Get user ID from path parameter
	userIDStr := c.Param("id")
	userID, err := h.ValidateUUID(userIDStr, "user_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user ID"))
		return
	}

	// Retrieve user
	user, err := h.userService.GetUserByID(c.Request.Context(), userID.String())
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "User not found"))
		return
	}

	// Respond with user details
	h.HandleSuccess(c, user, "User retrieved successfully")
}
