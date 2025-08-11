package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
	"github.com/holycann/cultour-backend/pkg/validator"
)

// UserProfileHandler handles user profile-related HTTP requests
// @Description Manages user profile operations such as creation, retrieval, update, and deletion
type UserProfileHandler struct {
	base.BaseHandler
	userProfileService services.UserProfileService
}

// NewUserProfileHandler creates a new instance of UserProfileHandler
// @Description Initializes a new UserProfileHandler with the provided UserProfileService and logger
func NewUserProfileHandler(userProfileService services.UserProfileService, logger *logger.Logger) *UserProfileHandler {
	return &UserProfileHandler{
		BaseHandler:        *base.NewBaseHandler(logger),
		userProfileService: userProfileService,
	}
}

// CreateUserProfile godoc
// @Summary Create a new user profile
// @Description Allows administrators or users to create a detailed user profile
// @Description Supports initializing profile with optional personal information
// @Tags User Profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param profile body models.UserProfileCreate true "User Profile Creation Details"
// @Success 201 {object} response.APIResponse{data=models.UserProfileDTO} "User profile successfully created"
// @Failure 400 {object} response.APIResponse "Invalid profile creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 409 {object} response.APIResponse "Profile already exists for the user"
// @Failure 500 {object} response.APIResponse "Internal server error during profile creation"
// @Router /users/profiles [post]
func (h *UserProfileHandler) CreateUserProfile(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Create a user profile model to bind request body
	var userProfileCreate models.UserProfileCreate

	// Validate request body
	if err := c.ShouldBindJSON(&userProfileCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	if userProfileCreate.UserID == uuid.Nil {
		// Set user ID from context
		userProfileCreate.UserID, err = h.ValidateUUID(userID, "User ID")
		if err != nil {
			h.HandleError(c, err)
			return
		}
	}

	// Validate user profile creation payload
	if err := validator.ValidateStruct(userProfileCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Create user profile through service
	createdProfile, err := h.userProfileService.CreateProfile(c.Request.Context(), &userProfileCreate)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to create user profile"))
		return
	}

	// Respond with created user profile
	h.HandleCreated(c, createdProfile, "User profile created successfully")
}

// ListUsersProfile godoc
// @Summary Retrieve users profiles list
// @Description Fetches a paginated list of user profiles with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of profiles per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort profiles by" default("created_at)" Enum(created_at,fullname)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param fullname query string false "Filter profiles by partial full name match"
// @Success 200 {object} response.APIResponse{data=[]models.UserProfileDTO} "Successfully retrieved user profiles list"
// @Success 204 {object} response.APIResponse "No user profiles found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during profiles retrieval"
// @Router /users/profiles [get]
func (h *UserProfileHandler) ListUsersProfile(c *gin.Context) {
	listOptions, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.New(errors.ErrBadRequest, err.Error(), err))
		return
	}

	// Retrieve user profiles
	profiles, totalProfiles, err := h.userProfileService.ListProfiles(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to retrieve user profiles"))
		return
	}

	// Respond with user profiles and pagination
	h.HandlePagination(c, profiles, totalProfiles, listOptions)
}

// SearchUserProfile godoc
// @Summary Search user profiles
// @Description Performs a full-text search across user profile details with advanced filtering
// @Description Allows finding user profiles by keywords, name, and other attributes
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding user profiles" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,fullname,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.UserProfileDTO} "Successfully completed user profile search"
// @Success 204 {object} response.APIResponse "No user profiles match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during user profile search"
// @Router /users/profiles/search [get]
func (h *UserProfileHandler) SearchUserProfile(c *gin.Context) {
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

	// Search user profiles
	profiles, totalProfiles, err := h.userProfileService.SearchProfiles(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to search user profiles"))
		return
	}

	// Respond with user profiles and pagination
	h.HandlePagination(c, profiles, totalProfiles, listOptions)
}

// UpdateUserProfile godoc
// @Summary Update an existing user profile
// @Description Allows users to modify their profile details
// @Description Supports partial updates with optional fields
// @Tags User Profiles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param profile body models.UserProfileUpdate true "User Profile Update Payload"
// @Success 200 {object} response.APIResponse{data=models.UserProfileDTO} "User profile successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid profile update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - can only update own profile"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error during profile update"
// @Router /users/profiles/{id} [put]
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Get user profile ID from path parameter
	profileIDStr := c.Param("id")
	profileID, err := h.ValidateUUID(profileIDStr, "profile_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user profile ID"))
		return
	}

	// Create a user profile update model to bind request body
	var updateProfile models.UserProfileUpdate

	// Bind and validate input
	if err := c.ShouldBindJSON(&updateProfile); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	if updateProfile.UserID == uuid.Nil {
		// Set user ID from context
		updateProfile.UserID, err = h.ValidateUUID(userID, "User ID")
		if err != nil {
			h.HandleError(c, err)
			return
		}
	}

	// Set the ID from path parameter
	updateProfile.ID = profileID

	// Validate user profile update payload
	if err := validator.ValidateStruct(updateProfile); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Update user profile
	updatedProfile, err := h.userProfileService.UpdateProfile(c.Request.Context(), &updateProfile)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update user profile"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, updatedProfile, "User profile updated successfully")
}

// UpdateUserAvatar godoc
// @Summary Update user profile avatar
// @Description Allows users to upload a new profile picture
// @Description Supports multipart file upload or URL-based avatar update
// @Tags User Profiles
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id formData string true "Unique User Profile Identifier" format(uuid)
// @Param avatar_url formData string true "URL of the new avatar image"
// @Param image formData file true "New Avatar Image File"
// @Success 200 {object} response.APIResponse{data=models.UserProfileDTO} "Avatar successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid avatar update payload or file"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - can only update own avatar"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error during avatar update"
// @Router /users/profiles/{id}/avatar [put]
func (h *UserProfileHandler) UpdateUserAvatar(c *gin.Context) {
	// Get user profile ID from path parameter
	profileIDStr := c.Param("id")
	profileID, err := h.ValidateUUID(profileIDStr, "profile_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user profile ID"))
		return
	}

	// Multipart form handling with max file size of 2MB
	const maxAvatarFileSize = 2 * 1024 * 1024 // 2MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAvatarFileSize)

	if err := c.Request.ParseMultipartForm(maxAvatarFileSize); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "File size exceeds maximum limit of 2MB"))
		} else {
			h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid multipart form request"))
		}
		return
	}

	// Get the avatar file from multipart form
	avatarFile, err := c.FormFile("avatar")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Avatar file is required"))
		return
	}

	// Create a user profile avatar update model
	var updateAvatarProfile models.UserProfileAvatarUpdate
	updateAvatarProfile.ID = profileID
	updateAvatarProfile.Image = avatarFile

	// Update user profile avatar
	updatedProfile, err := h.userProfileService.UpdateProfileAvatar(c.Request.Context(), &updateAvatarProfile)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update user profile avatar"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, updatedProfile, "User profile avatar updated successfully")
}

// VerifyIdentity godoc
// @Summary Update user profile identity
// @Description Allows users to upload a new identity document
// @Description Supports multipart file upload for identity verification
// @Tags User Profiles
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique User Profile Identifier" format(uuid)
// @Param identity_image formData file true "Identity Document Image" format(binary)
// @Success 200 {object} response.APIResponse{data=models.UserProfileDTO} "Identity document successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid identity update payload or file"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - can only update own identity"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 413 {object} response.APIResponse "File size too large"
// @Failure 415 {object} response.APIResponse "Unsupported file type"
// @Failure 500 {object} response.APIResponse "Internal server error during identity update"
// @Router /users/profiles/{id}/verify [post]
func (h *UserProfileHandler) VerifyIdentity(c *gin.Context) {
	// Get user profile ID from path parameter
	profileIDStr := c.Param("id")
	profileID, err := h.ValidateUUID(profileIDStr, "profile_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user profile ID"))
		return
	}

	// Multipart form handling with max file size of 2MB
	const maxIdentityFileSize = 2 * 1024 * 1024 // 2MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxIdentityFileSize)

	if err := c.Request.ParseMultipartForm(maxIdentityFileSize); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "File size exceeds maximum limit of 2MB"))
		} else {
			h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid multipart form request"))
		}
		return
	}

	// Log form field keys
	for key := range c.Request.Form {
		fmt.Printf("Form field key: %s\n", key)
	}

	// Get the identity document file from multipart form
	identityFile, err := c.FormFile("identity_image")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Identity document image is required"))
		return
	}

	// Create a user profile identity update model
	var updateIdentityProfile models.UserProfileIdentityUpdate
	updateIdentityProfile.ID = profileID
	updateIdentityProfile.Image = identityFile

	// Update user profile identity
	updatedProfile, err := h.userProfileService.UpdateProfileIdentity(c.Request.Context(), &updateIdentityProfile)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update user profile identity document"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, updatedProfile, "User profile identity document updated successfully")
}

// DeleteUserProfile godoc
// @Summary Delete a user profile
// @Description Allows administrators to permanently remove a user profile
// @Description Deletes the profile and associated user information
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique User Profile Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "User profile successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid profile ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error during profile deletion"
// @Router /users/profiles/{id} [delete]
func (h *UserProfileHandler) DeleteUserProfile(c *gin.Context) {
	// Get user profile ID from path parameter
	profileIDStr := c.Param("id")
	profileID, err := h.ValidateUUID(profileIDStr, "profile_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user profile ID"))
		return
	}

	// Delete user profile
	if err := h.userProfileService.DeleteProfile(c.Request.Context(), profileID.String()); err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with success
	h.HandleSuccess(c, gin.H{
		"id": profileID,
	}, "User profile deleted successfully")
}

// GetUserProfileById godoc
// @Summary Retrieve a specific user profile
// @Description Fetches comprehensive details of a user profile by its unique identifier
// @Description Returns full profile information including user details
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique User Profile Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.UserProfileDTO} "Successfully retrieved user profile details"
// @Failure 400 {object} response.APIResponse "Invalid profile ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error during profile retrieval"
// @Router /users/profiles/{id} [get]
func (h *UserProfileHandler) GetUserProfileById(c *gin.Context) {
	// Get user profile ID from path parameter
	profileIDStr := c.Param("id")
	profileID, err := h.ValidateUUID(profileIDStr, "profile_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid user profile ID"))
		return
	}

	// Retrieve user profile
	profile, err := h.userProfileService.GetProfileByID(c.Request.Context(), profileID.String())
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "User profile not found"))
		return
	}

	// Respond with user profile details
	h.HandleSuccess(c, profile, "User profile retrieved successfully")
}

// GetAuthenticatedUserProfile godoc
// @Summary Retrieve the current user's profile
// @Description Fetches the profile details of the authenticated user
// @Description Returns comprehensive profile information for the logged-in user
// @Tags User Profiles
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Success 200 {object} response.APIResponse{data=models.UserProfileDTO} "Successfully retrieved authenticated user profile"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "User profile not found"
// @Failure 500 {object} response.APIResponse "Internal server error during profile retrieval"
// @Router /users/profiles/me [get]
func (h *UserProfileHandler) GetAuthenticatedUserProfile(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Retrieve user profile
	profile, err := h.userProfileService.GetProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with user profile details
	h.HandleSuccess(c, profile, "User profile retrieved successfully")
}
