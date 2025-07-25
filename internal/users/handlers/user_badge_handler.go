package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
)

// UserBadgeHandler handles HTTP requests related to user badges
type UserBadgeHandler struct {
	service *services.UserBadgeService
	logger  *logger.Logger
}

// NewUserBadgeHandler creates a new instance of UserBadgeHandler
func NewUserBadgeHandler(service *services.UserBadgeService, logger *logger.Logger) *UserBadgeHandler {
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
// @Success 201 {object} response.Response{data=models.UserBadge} "Badge assigned successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid badge assignment details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/badges [post]
func (h *UserBadgeHandler) AssignBadge(c *gin.Context) {
	var badgeCreate models.UserBadgeCreate
	if err := c.ShouldBindJSON(&badgeCreate); err != nil {
		response.BadRequest(c, "Invalid badge assignment details", err)
		return
	}

	// You might want to add additional validation or get more badge details from a badge repository
	badge, err := h.service.AssignBadge(
		c.Request.Context(),
		badgeCreate.UserID,
		badgeCreate.BadgeID,
		"Penjelajah",         // Default badge name, replace with actual logic
		"Earned a new badge", // Default description, replace with actual logic
		"",                   // Default icon URL, replace with actual logic
	)
	if err != nil {
		h.logger.Error("Failed to assign badge", err)
		response.InternalServerError(c, "Failed to assign badge", err)
		return
	}

	response.SuccessCreated(c, badge, "Badge successfully assigned")
}

// GetUserBadges godoc
// @Summary Get user badges
// @Description Retrieve badges for a specific user
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "User ID"
// @Param limit query int false "Number of badges to retrieve" default(10)
// @Param offset query int false "Number of badges to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.UserBadge} "User badges retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/badges [get]
func (h *UserBadgeHandler) GetUserBadges(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	limit, offset := 10, 0
	var err error
	if limit, err = strconv.Atoi(c.DefaultQuery("limit", "10")); err != nil {
		limit = 10
	}
	if offset, err = strconv.Atoi(c.DefaultQuery("offset", "0")); err != nil {
		offset = 0
	}

	search := &models.UserBadgeSearch{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	badges, err := h.service.GetUserBadges(c.Request.Context(), userID, search)
	if err != nil {
		h.logger.Error("Failed to retrieve user badges", err)
		response.InternalServerError(c, "Failed to retrieve user badges", err)
		return
	}

	response.Success(c, 200, badges, "User badges retrieved successfully")
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
// @Success 200 {object} response.Response "Badge removed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID or badge ID"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/badges [delete]
func (h *UserBadgeHandler) RemoveBadge(c *gin.Context) {
	userID := c.Query("user_id")
	badgeID := c.Query("badge_id")

	if userID == "" || badgeID == "" {
		response.BadRequest(c, "User ID and Badge ID are required", nil)
		return
	}

	err := h.service.RemoveBadge(c.Request.Context(), userID, badgeID)
	if err != nil {
		h.logger.Error("Failed to remove badge", err)
		response.InternalServerError(c, "Failed to remove badge", err)
		return
	}

	response.Success(c, 200, nil, "Badge removed successfully")
}

// CountUserBadges godoc
// @Summary Count user badges
// @Description Retrieve the number of badges a user has
// @Tags User Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param user_id query string true "User ID"
// @Success 200 {object} response.Response{data=int} "User badge count retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/badges/count [get]
func (h *UserBadgeHandler) CountUserBadges(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		response.BadRequest(c, "User ID is required", nil)
		return
	}

	count, err := h.service.CountUserBadges(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to count user badges", err)
		response.InternalServerError(c, "Failed to count user badges", err)
		return
	}

	response.Success(c, 200, count, "User badges counted successfully")
}
