package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
)

// UserProfileHandler handles user profile-related HTTP requests
// @Description Manages user profile operations such as creation, retrieval, update, and deletion
type UserProfileHandler struct {
	userProfileService services.UserProfileService
}

// NewUserProfileHandler creates a new instance of UserProfileHandler
// @Description Initializes a new UserProfileHandler with the provided UserProfileService
func NewUserProfileHandler(userProfileService services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService: userProfileService,
	}
}

// CreateUserProfile godoc
// @Summary Create a new user profile
// @Description Register a new user profile in the system with provided details
// @Tags User Profiles
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param UserID header string true "User ID"
// @Param profile body models.UserProfileCreate true "User profile creation details"
// @Success 201 {object} models.UserProfile "User profile created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 409 {object} response.ErrorResponse "User profile already exists"
// @Router /profile [post]
// @Security ApiKeyAuth
func (h *UserProfileHandler) CreateUserProfile(c *gin.Context) {
	// Create a user model to bind request body
	var userProfile models.UserProfile

	// Bind and validate input
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Validate required fields
	if userProfile.UserID == "" || userProfile.Fullname == "" {
		response.BadRequest(c, "Missing required fields", map[string]interface{}{
			"user_id":  userProfile.UserID == "",
			"fullname": userProfile.Fullname == "",
		})
		return
	}

	// Create user through service
	if err := h.userProfileService.CreateProfile(c.Request.Context(), &userProfile); err != nil {
		response.Conflict(c, "Failed to create user profile", err.Error())
		return
	}

	// Respond with created user profile
	response.SuccessCreated(c, gin.H{
		"id":       userProfile.ID,
		"fullname": userProfile.Fullname,
	}, "User Profile created successfully")
}

// ListUsersProfile godoc
// @Summary List user profiles
// @Description Retrieve a paginated list of user profiles
// @Tags User Profiles
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param limit query int false "Number of user profiles to retrieve" default(10)
// @Param offset query int false "Starting offset for pagination" default(0)
// @Success 200 {object} map[string]interface{} "Successfully retrieved user profiles"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user profiles"
// @Router /profile [get]
// @Security ApiKeyAuth
func (h *UserProfileHandler) ListUsersProfile(c *gin.Context) {
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
	usersProfile, err := h.userProfileService.GetProfiles(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users profile", err.Error())
		return
	}

	// Use WithPagination for consistent pagination response
	response.WithPagination(c, usersProfile, len(usersProfile), offset/limit+1, limit)
}

// SearchUserProfile godoc
// @Summary Search user profiles
// @Description Retrieve user profiles by ID, User ID, or with pagination
// @Tags User Profiles
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id path string false "User Profile ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} map[string]interface{} "Successfully retrieved user profiles"
// @Failure 400 {object} response.ErrorResponse "Invalid input parameters"
// @Failure 404 {object} response.ErrorResponse "User profile not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user profiles"
// @Router /profile/search [get]
// @Security ApiKeyAuth
func (h *UserProfileHandler) SearchUserProfile(c *gin.Context) {
	// Check if specific ID is provided
	userProfileID := c.Param("id")
	userID := c.Query("user_id")

	// Validate search parameters
	if userProfileID == "" && userID == "" {
		response.BadRequest(c, "Either User Profile ID or User ID is required", nil)
		return
	}

	var userProfile *models.UserProfile
	var err error

	// Search by ID if provided
	if userID != "" {
		userProfile, err = h.userProfileService.GetProfileByID(c.Request.Context(), userID)
	} else {
		// Otherwise, search by email
		userProfile, err = h.userProfileService.GetProfileByUserID(c.Request.Context(), userID)
	}

	// Handle search errors
	if err != nil {
		response.NotFound(c, "User Profile not found", err.Error())
		return
	}

	response.SuccessOK(c, userProfile, "User Profile retrieved successfully")
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update an existing user profile's details
// @Tags User Profiles
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id path string true "User Profile ID"
// @Param profile body models.UserProfile true "User profile update details"
// @Success 200 {object} map[string]interface{} "User profile updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "User profile not found"
// @Failure 409 {object} response.ErrorResponse "Failed to update user profile"
// @Router /profile/{id} [put]
// @Security ApiKeyAuth
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User Profile ID is required", nil)
		return
	}

	// Create a user model to bind request body
	var updateUserProfile models.UserProfile

	// Bind input
	if err := c.ShouldBindJSON(&updateUserProfile); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set the ID from path parameter
	updateUserProfile.ID = userID

	// Update user
	if err := h.userProfileService.UpdateProfile(c.Request.Context(), &updateUserProfile); err != nil {
		response.Conflict(c, "Failed to update user profile", err.Error())
		return
	}

	// Respond with success
	response.SuccessOK(c, updateUserProfile, "User Profile updated successfully")
}

// DeleteUserProfile godoc
// @Summary Delete user profile
// @Description Soft delete a user profile by its unique identifier
// @Tags User Profiles
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param id path string true "User Profile ID"
// @Success 200 {object} map[string]string "User profile deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user profile ID"
// @Failure 404 {object} response.ErrorResponse "User profile not found"
// @Failure 409 {object} response.ErrorResponse "Failed to delete user profile"
// @Router /profile/{id} [delete]
// @Security ApiKeyAuth
func (h *UserProfileHandler) DeleteUserProfile(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User Profile ID is required", nil)
		return
	}

	// Delete user
	if err := h.userProfileService.DeleteProfile(c.Request.Context(), userID); err != nil {
		response.Conflict(c, "Failed to delete user profile", err.Error())
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User Profile deleted successfully")
}
