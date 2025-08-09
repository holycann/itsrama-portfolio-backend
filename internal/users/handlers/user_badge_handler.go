package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
)

// UserBadgeHandler handles HTTP requests related to user badges
type UserBadgeHandler struct {
	base.BaseHandler
	service services.UserBadgeService
}

// NewUserBadgeHandler creates a new instance of UserBadgeHandler
func NewUserBadgeHandler(service services.UserBadgeService, logger *logger.Logger) *UserBadgeHandler {
	return &UserBadgeHandler{
		BaseHandler: *base.NewBaseHandler(logger),
		service:     service,
	}
}

// AssignBadge godoc
// @Summary Assign a badge to a user
// @Description Allows administrators to award a specific badge to a user
// @Description Creates a new user badge association in the system
// @Tags User Badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param badge body models.UserBadgePayload true "Badge Assignment Details"
// @Success 201 {object} response.APIResponse{data=models.UserBadgeDTO} "Badge successfully assigned to user"
// @Failure 400 {object} response.APIResponse "Invalid badge assignment payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges (admin role required)"
// @Failure 404 {object} response.APIResponse "User or badge not found"
// @Failure 409 {object} response.APIResponse "Badge already assigned to user"
// @Failure 500 {object} response.APIResponse "Internal server error during badge assignment"
// @Router /users/badges [post]
func (h *UserBadgeHandler) AssignBadge(c *gin.Context) {
	var badgePayload models.UserBadge

	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Use authenticated user's ID
	userUUID, err := h.ValidateUUID(userID, "user_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}
	badgePayload.UserID = userUUID
	if err := h.ValidateRequest(c, &badgePayload); err != nil {
		h.HandleError(c, err)
		return
	}

	// Add badge to user
	err = h.service.AddBadgeToUser(c.Request.Context(), badgePayload)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Retrieve the newly created badge
	badges, err := h.service.GetUserBadgesByUser(c.Request.Context(), userID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, badges, "Badge successfully assigned")
}

// GetUserBadges godoc
// @Summary Retrieve user's badges
// @Description Fetches a list of badges earned by a specific user
// @Description Returns comprehensive badge details with optional filtering
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "Unique User Identifier" format(uuid)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of badges per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort badges by" default("created_at)" Enum(created_at,badge_name)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.UserBadgeDTO} "Successfully retrieved user badges"
// @Success 204 {object} response.APIResponse "No badges found for the user"
// @Failure 400 {object} response.APIResponse "Invalid user ID or query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge retrieval"
// @Router /users/badges [get]
func (h *UserBadgeHandler) GetUserBadges(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
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
	listOptions := base.ListOptions{
		Page:      offset/limit + 1,
		PerPage:   limit,
		SortBy:    sortBy,
		SortOrder: base.SortDescending,
		Filters: []base.FilterOption{
			{
				Field:    "user_id",
				Operator: base.OperatorEqual,
				Value:    userID,
			},
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = base.SortAscending
	}

	// Retrieve user badges
	badges, totalBadges, err := h.service.ListUserBadges(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with badges and pagination
	h.HandlePagination(c, badges, totalBadges, listOptions)
}

// RemoveBadge godoc
// @Summary Remove a badge from a user
// @Description Allows administrators to revoke a specific badge from a user
// @Description Permanently deletes the user badge association
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "Unique User Identifier" format(uuid)
// @Param badge_id query string true "Unique Badge Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Badge successfully removed from user"
// @Failure 400 {object} response.APIResponse "Invalid user or badge ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "User badge association not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge removal"
// @Router /users/badges [delete]
func (h *UserBadgeHandler) RemoveBadge(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	badgeID := c.Query("badge_id")
	if badgeID == "" {
		h.HandleError(c, fmt.Errorf("Badge ID is required"))
		return
	}

	// Parse UUIDs
	userUUID, err := h.ValidateUUID(userID, "user_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}
	badgeUUID, err := h.ValidateUUID(badgeID, "badge_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Prepare payload
	payload := models.UserBadgePayload{
		UserID:  userUUID,
		BadgeID: badgeUUID,
	}

	// Remove badge from user
	err = h.service.RemoveBadgeFromUser(c.Request.Context(), payload)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Badge removed successfully")
}

// CountUserBadges godoc
// @Summary Count user's badges
// @Description Retrieves the total number of badges earned by a specific user
// @Description Supports optional filtering by badge type or name
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "Unique User Identifier" format(uuid)
// @Param badge_name query string false "Optional filter to count badges by name"
// @Success 200 {object} response.APIResponse{data=int} "Successfully retrieved badge count"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge count"
// @Router /users/badges/count [get]
func (h *UserBadgeHandler) CountUserBadges(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Prepare filter
	filters := []base.FilterOption{
		{
			Field:    "user_id",
			Operator: base.OperatorEqual,
			Value:    userID,
		},
	}

	count, err := h.service.CountUserBadges(c.Request.Context(), filters)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, count, "User badges counted successfully")
}

// GetUserBadgesByUser godoc
// @Summary Retrieve detailed user badge information
// @Description Fetches comprehensive badge details for a specific user
// @Description Returns full badge information including badge metadata
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param user_id path string true "Unique User Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=[]models.UserBadgeDTO} "Successfully retrieved user badge details"
// @Success 204 {object} response.APIResponse "No badges found for the user"
// @Failure 400 {object} response.APIResponse "Invalid user ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "User not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge retrieval"
// @Router /users/{user_id}/badges [get]
func (h *UserBadgeHandler) GetUserBadgesByUser(c *gin.Context) {
	// Get authenticated user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	userBadges, err := h.service.GetUserBadgesByUser(c.Request.Context(), userID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, userBadges, "User badges retrieved successfully")
}
