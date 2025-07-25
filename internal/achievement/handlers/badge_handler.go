package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
)

// BadgeHandler handles HTTP requests related to badges
type BadgeHandler struct {
	service services.BadgeService
	logger  *logger.Logger
}

// NewBadgeHandler creates a new instance of BadgeHandler
func NewBadgeHandler(service services.BadgeService, logger *logger.Logger) *BadgeHandler {
	return &BadgeHandler{
		service: service,
		logger:  logger,
	}
}

// CreateBadge godoc
// @Summary Create a new badge
// @Description Add a new badge to the system
// @Tags badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param badge body models.BadgeCreate true "Badge Information"
// @Success 201 {object} response.Response{data=models.Badge} "Badge created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid badge creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /badges [post]
func (h *BadgeHandler) CreateBadge(c *gin.Context) {
	var badgeCreate models.BadgeCreate
	if err := c.ShouldBindJSON(&badgeCreate); err != nil {
		response.BadRequest(c, "Invalid badge creation details", err)
		return
	}

	badge, err := h.service.CreateBadge(c.Request.Context(), &badgeCreate)
	if err != nil {
		h.logger.Error("Failed to create badge", err)
		response.InternalServerError(c, "Failed to create badge", err)
		return
	}

	response.SuccessCreated(c, badge, "Badge created successfully")
}

// GetBadgeByID godoc
// @Summary Get a specific badge
// @Description Retrieve a badge by its unique identifier
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Badge ID"
// @Success 200 {object} response.Response{data=models.Badge} "Badge retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Badge ID is required"
// @Failure 404 {object} response.ErrorResponse "Badge not found"
// @Router /badges/{id} [get]
func (h *BadgeHandler) GetBadgeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Badge ID is required", nil)
		return
	}

	badge, err := h.service.GetBadgeByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to retrieve badge", err)
		response.NotFound(c, "Badge not found", err)
		return
	}

	response.Success(c, 200, badge, "Badge retrieved successfully")
}

// ListBadges godoc
// @Summary List badges
// @Description Retrieve a list of badges with pagination
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of badges to retrieve" default(10)
// @Param offset query int false "Number of badges to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Badge} "Badges retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list badges"
// @Router /badges [get]
func (h *BadgeHandler) ListBadges(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	badges, err := h.service.ListBadges(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list badges", err)
		response.InternalServerError(c, "Failed to list badges", err)
		return
	}

	response.Success(c, 200, badges, "Badges retrieved successfully")
}

// UpdateBadge godoc
// @Summary Update a badge
// @Description Update an existing badge by its ID
// @Tags badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Badge ID"
// @Param badge body models.BadgeCreate true "Badge Update Information"
// @Success 200 {object} response.Response{data=models.Badge} "Badge updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid badge update details or missing ID"
// @Failure 500 {object} response.ErrorResponse "Failed to update badge"
// @Router /badges/{id} [put]
func (h *BadgeHandler) UpdateBadge(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Badge ID is required", nil)
		return
	}

	var badgeUpdate models.BadgeCreate
	if err := c.ShouldBindJSON(&badgeUpdate); err != nil {
		response.BadRequest(c, "Invalid badge update details", err)
		return
	}

	badge, err := h.service.UpdateBadge(c.Request.Context(), id, &badgeUpdate)
	if err != nil {
		h.logger.Error("Failed to update badge", err)
		response.InternalServerError(c, "Failed to update badge", err)
		return
	}

	response.Success(c, 200, badge, "Badge updated successfully")
}

// DeleteBadge godoc
// @Summary Delete a badge
// @Description Remove a badge from the system by its ID
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Badge ID"
// @Success 200 {object} response.Response "Badge deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Badge ID is required"
// @Failure 500 {object} response.ErrorResponse "Failed to delete badge"
// @Router /badges/{id} [delete]
func (h *BadgeHandler) DeleteBadge(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Badge ID is required", nil)
		return
	}

	err := h.service.DeleteBadge(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete badge", err)
		response.InternalServerError(c, "Failed to delete badge", err)
		return
	}

	response.Success(c, 200, nil, "Badge deleted successfully")
}

// CountBadges godoc
// @Summary Count badges
// @Description Retrieve the total number of badges in the system
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Success 200 {object} response.Response{data=int} "Badge count retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to count badges"
// @Router /badges/count [get]
func (h *BadgeHandler) CountBadges(c *gin.Context) {
	count, err := h.service.CountBadges(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to count badges", err)
		response.InternalServerError(c, "Failed to count badges", err)
		return
	}

	response.Success(c, 200, count, "Badge count retrieved successfully")
}
