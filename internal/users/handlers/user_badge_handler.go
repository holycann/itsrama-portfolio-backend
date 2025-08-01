package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// UserBadgeHandler handles HTTP requests related to user badges
type UserBadgeHandler struct {
	service services.UserBadgeService
	logger  *logger.Logger
}

// NewUserBadgeHandler creates a new instance of UserBadgeHandler
func NewUserBadgeHandler(service services.UserBadgeService, logger *logger.Logger) *UserBadgeHandler {
	return &UserBadgeHandler{
		service: service,
		logger:  logger,
	}
}

// AssignBadge godoc
// @Summary Assign a badge to a user
// @Description Add a new badge to a user's profile
// @Tags User Badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param badge body models.UserBadgeCreate true "Badge Assignment Details"
// @Success 201 {object} response.APIResponse{data=models.UserBadge} "Badge assigned successfully"
// @Failure 400 {object} response.APIResponse "Invalid badge assignment details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/badges [post]
func (h *UserBadgeHandler) AssignBadge(c *gin.Context) {
	var badgeCreate models.UserBadgeCreate
	if err := c.ShouldBindJSON(&badgeCreate); err != nil {
		response.BadRequest(c, "Invalid badge assignment details", err.Error(), "")
		return
	}

	// Validate required fields
	if badgeCreate.UserID.String() == "" || badgeCreate.BadgeID.String() == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"user_id":  badgeCreate.UserID.String() == "",
			"badge_id": badgeCreate.BadgeID.String() == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	// Create user badge
	err := h.service.CreateUserBadge(c.Request.Context(), &badgeCreate)
	if err != nil {
		h.logger.Error("Failed to assign badge", err)
		response.InternalServerError(c, "Failed to assign badge", err.Error(), "")
		return
	}

	// Retrieve the newly created badge
	badge, err := h.service.GetUserBadgesByUser(c.Request.Context(), badgeCreate.UserID.String())
	if err != nil {
		h.logger.Error("Failed to retrieve assigned badge", err)
		response.InternalServerError(c, "Failed to retrieve assigned badge", err.Error(), "")
		return
	}

	response.SuccessCreated(c, badge, "Badge successfully assigned")
}

// GetUserBadges godoc
// @Summary Get user badges
// @Description Retrieve badges for a specific user with pagination and filtering
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "User ID"
// @Param limit query int false "Number of badges to retrieve" default(10)
// @Param offset query int false "Number of badges to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.UserBadge} "User badges retrieved successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/badges [get]
func (h *UserBadgeHandler) GetUserBadges(c *gin.Context) {
	// Get user ID
	userID := c.Query("user_id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", "Missing user ID", "")
		return
	}

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
		Filters: []repository.FilterOption{
			{
				Field:    "user_id",
				Operator: "=",
				Value:    userID,
			},
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Retrieve user badges
	badges, err := h.service.ListUserBadges(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Failed to retrieve user badges", err)
		response.InternalServerError(c, "Failed to retrieve user badges", err.Error(), "")
		return
	}

	// Count total user badges for pagination
	totalBadges, err := h.service.CountUserBadges(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Failed to count user badges", err)
		response.InternalServerError(c, "Failed to count user badges", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalBadges,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalBadges) / float64(limit))),
		HasNextPage: offset+limit < totalBadges,
	}

	// Respond with badges and pagination
	response.SuccessOK(c, badges, "User badges retrieved successfully", pagination)
}

// RemoveBadge godoc
// @Summary Remove a badge from a user
// @Description Delete a specific badge from a user's profile
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "User ID"
// @Param badge_id query string true "Badge ID"
// @Success 200 {object} response.APIResponse "Badge removed successfully"
// @Failure 400 {object} response.APIResponse "Invalid user ID or badge ID"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/badges [delete]
func (h *UserBadgeHandler) RemoveBadge(c *gin.Context) {
	userID := c.Query("user_id")
	badgeID := c.Query("badge_id")

	if userID == "" || badgeID == "" {
		response.BadRequest(c, "User ID and Badge ID are required", "Missing user or badge ID", "")
		return
	}

	// Parse UUIDs
	_, err := uuid.Parse(userID)
	if err != nil {
		response.BadRequest(c, "Invalid User ID", "Invalid UUID format for user ID", "")
		return
	}
	badgeUUID, err := uuid.Parse(badgeID)
	if err != nil {
		response.BadRequest(c, "Invalid Badge ID", "Invalid UUID format for badge ID", "")
		return
	}

	// Find the specific user badge to delete
	userBadges, err := h.service.GetUserBadgesByUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to find user badges", err)
		response.InternalServerError(c, "Failed to find user badges", err.Error(), "")
		return
	}

	// Find the specific badge to delete
	var badgeToDelete *models.UserBadge
	for _, badge := range userBadges {
		if badge.BadgeID == badgeUUID {
			badgeToDelete = &badge
			break
		}
	}

	if badgeToDelete == nil {
		response.NotFound(c, "Badge not found", "User does not have this badge", "")
		return
	}

	// Delete the badge
	err = h.service.DeleteUserBadge(c.Request.Context(), badgeToDelete.ID.String())
	if err != nil {
		h.logger.Error("Failed to remove badge", err)
		response.InternalServerError(c, "Failed to remove badge", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Badge removed successfully")
}

// CountUserBadges godoc
// @Summary Count user badges
// @Description Retrieve the number of badges a user has
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "User ID"
// @Success 200 {object} response.APIResponse{data=int} "User badge count retrieved successfully"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/badges/count [get]
func (h *UserBadgeHandler) CountUserBadges(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", "Missing user ID", "")
		return
	}

	// Prepare filter
	filters := []repository.FilterOption{
		{
			Field:    "user_id",
			Operator: "=",
			Value:    userID,
		},
	}

	count, err := h.service.CountUserBadges(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to count user badges", err)
		response.InternalServerError(c, "Failed to count user badges", err.Error(), "")
		return
	}

	response.SuccessOK(c, count, "User badges counted successfully")
}

// GetUserBadgesByUser godoc
// @Summary Get badges for a specific user
// @Description Retrieve all badges associated with a given user ID
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id path string true "User ID"
// @Success 200 {object} response.APIResponse{data=[]models.UserBadge} "User badges retrieved successfully"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /users/badges/{user_id} [get]
func (h *UserBadgeHandler) GetUserBadgesByUser(c *gin.Context) {
	// Verify the authenticated user's permission to access the badges
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "Authentication failed", "User not authenticated", "")
		return
	}

	userBadges, err := h.service.GetUserBadgesByUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to retrieve user badges", err)
		response.InternalServerError(c, "Failed to retrieve user badges", err.Error(), "")
		return
	}

	response.SuccessOK(c, userBadges, "User badges retrieved successfully")
}
