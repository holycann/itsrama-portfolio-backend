package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/repository"
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
// @Description Register a new user profile in the system
// @Tags User Profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param profile body models.UserProfileCreate true "User Profile Creation Details"
// @Success 201 {object} response.APIResponse{data=models.UserProfile} "User profile created successfully"
// @Failure 400 {object} response.APIResponse "Invalid user profile creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /profile [post]
func (h *UserProfileHandler) CreateUserProfile(c *gin.Context) {
	// Create a user model to bind request body
	var userProfile models.UserProfile

	// Bind and validate input
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error(), "")
		return
	}

	// Validate required fields
	if userProfile.UserID.String() == "" || userProfile.Fullname == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"user_id":  userProfile.UserID.String() == "",
			"fullname": userProfile.Fullname == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	// Create user through service
	if err := h.userProfileService.CreateProfile(c.Request.Context(), &userProfile); err != nil {
		response.Conflict(c, "Failed to create user profile", err.Error(), "")
		return
	}

	// Respond with created user profile
	response.SuccessCreated(c, userProfile, "User Profile created successfully")
}

// ListUsersProfile godoc
// @Summary List user profiles
// @Description Retrieve a list of user profiles with pagination and filtering
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of user profiles to retrieve" default(10)
// @Param offset query int false "Number of user profiles to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.UserProfile} "User profiles retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list user profiles"
// @Router /profile [get]
func (h *UserProfileHandler) ListUsersProfile(c *gin.Context) {
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
	if fullname := c.Query("fullname"); fullname != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "fullname",
			Operator: "like",
			Value:    fullname,
		})
	}
	listOptions.Filters = filters

	// Retrieve users
	usersProfile, err := h.userProfileService.ListProfiles(c.Request.Context(), listOptions)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users profile", err.Error(), "")
		return
	}

	// Count total users for pagination
	totalProfiles, err := h.userProfileService.CountProfiles(c.Request.Context(), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to count user profiles", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalProfiles,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalProfiles) / float64(limit))),
		HasNextPage: offset+limit < totalProfiles,
	}

	// Respond with users and pagination
	response.SuccessOK(c, usersProfile, "User Profiles retrieved successfully", pagination)
}

// SearchUserProfile godoc
// @Summary Search user profiles
// @Description Search user profiles by various criteria
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (fullname, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.APIResponse{data=[]models.UserProfile} "User profiles found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /profile/search [get]
func (h *UserProfileHandler) SearchUserProfile(c *gin.Context) {
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
				Field:    "fullname",
				Operator: "like",
				Value:    query,
			},
		},
	}

	// Search user profiles
	userProfiles, err := h.userProfileService.SearchProfiles(c.Request.Context(), query, listOptions)
	if err != nil {
		response.InternalServerError(c, "Failed to search user profiles", err.Error(), "")
		return
	}

	// Count total search results
	totalProfiles, err := h.userProfileService.CountProfiles(c.Request.Context(), listOptions.Filters)
	if err != nil {
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalProfiles,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalProfiles) / float64(limit))),
		HasNextPage: offset+limit < totalProfiles,
	}

	// Respond with users and pagination
	response.SuccessOK(c, userProfiles, "User Profiles found successfully", pagination)
}

// UpdateUserProfile godoc
// @Summary Update a user profile
// @Description Update an existing user profile's details
// @Tags User Profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "User Profile ID"
// @Param profile body models.UserProfile true "User Profile Update Details"
// @Success 200 {object} response.APIResponse{data=models.UserProfile} "User profile updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid user profile update details"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /profile/{id} [put]
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User Profile ID is required", "Missing user profile ID", "")
		return
	}

	// Create a user model to bind request body
	var updateUserProfile models.UserProfile

	// Bind input
	if err := c.ShouldBindJSON(&updateUserProfile); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		response.BadRequest(c, "Invalid User Profile ID", "Invalid UUID format", "")
		return
	}
	updateUserProfile.ID = parsedID

	// Update user
	if err := h.userProfileService.UpdateProfile(c.Request.Context(), &updateUserProfile); err != nil {
		response.Conflict(c, "Failed to update user profile", err.Error(), "")
		return
	}

	// Respond with success
	response.SuccessOK(c, updateUserProfile, "User Profile updated successfully")
}

// DeleteUserProfile godoc
// @Summary Delete a user profile
// @Description Remove a user profile from the system by its unique identifier
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "User Profile ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid user profile ID"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /profile/{id} [delete]
func (h *UserProfileHandler) DeleteUserProfile(c *gin.Context) {
	// Get user ID from path parameter
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "User Profile ID is required", "Missing user profile ID", "")
		return
	}

	// Delete user
	if err := h.userProfileService.DeleteProfile(c.Request.Context(), userID); err != nil {
		response.Conflict(c, "Failed to delete user profile", err.Error(), "")
		return
	}

	// Respond with success
	response.SuccessOK(c, gin.H{
		"id": userID,
	}, "User Profile deleted successfully")
}
