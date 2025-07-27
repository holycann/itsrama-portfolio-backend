package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
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
// @Success 201 {object} response.APIResponse{data=models.Badge} "Badge created successfully"
// @Failure 400 {object} response.APIResponse "Invalid badge creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /badges [post]
func (h *BadgeHandler) CreateBadge(c *gin.Context) {
	var badgeCreate models.BadgeCreate
	if err := c.ShouldBindJSON(&badgeCreate); err != nil {
		h.logger.Error("Error binding badge: %v", err)
		response.BadRequest(c, "Invalid badge creation details", err.Error(), "")
		return
	}

	// Validate required fields
	if badgeCreate.Name == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name": badgeCreate.Name == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	err := h.service.CreateBadge(c.Request.Context(), &badgeCreate)
	if err != nil {
		h.logger.Error("Failed to create badge: %v", err)
		response.InternalServerError(c, "Failed to create badge", err.Error(), "")
		return
	}

	response.SuccessCreated(c, badgeCreate, "Badge created successfully")
}

// GetBadgeByID godoc
// @Summary Get a specific badge
// @Description Retrieve a badge by its unique identifier
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Badge ID"
// @Success 200 {object} response.APIResponse{data=models.Badge} "Badge retrieved successfully"
// @Failure 400 {object} response.APIResponse "Badge ID is required"
// @Failure 404 {object} response.APIResponse "Badge not found"
// @Router /badges/{id} [get]
func (h *BadgeHandler) GetBadgeByID(c *gin.Context) {
	// Get badge ID from path parameter
	badgeID := c.Param("id")
	if badgeID == "" {
		response.BadRequest(c, "Badge ID is required", "Missing badge ID", "")
		return
	}

	badge, err := h.service.GetBadgeByID(c.Request.Context(), badgeID)
	if err != nil {
		h.logger.Error("Failed to retrieve badge: %v", err)
		response.NotFound(c, "Badge not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, badge, "Badge retrieved successfully")
}

// ListBadges godoc
// @Summary List badges
// @Description Retrieve a list of badges with pagination and filtering
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of badges to retrieve" default(10)
// @Param offset query int false "Number of badges to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.Badge} "Badges retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list badges"
// @Router /badges [get]
func (h *BadgeHandler) ListBadges(c *gin.Context) {
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

	// Optional filtering
	filters := []repository.FilterOption{}
	if name := c.Query("name"); name != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "name",
			Operator: "like",
			Value:    name,
		})
	}

	// Retrieve badges
	badges, err := h.service.ListBadges(c.Request.Context(), repository.ListOptions{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: repository.SortDescending,
		Filters:   filters,
	})
	if sortOrder == "asc" {
		badges, err = h.service.ListBadges(c.Request.Context(), repository.ListOptions{
			Limit:     limit,
			Offset:    offset,
			SortBy:    sortBy,
			SortOrder: repository.SortAscending,
			Filters:   filters,
		})
	}
	if err != nil {
		h.logger.Error("Failed to list badges: %v", err)
		response.InternalServerError(c, "Failed to list badges", err.Error(), "")
		return
	}

	// Count total badges for pagination
	totalBadges, err := h.service.CountBadges(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to count badges: %v", err)
		response.InternalServerError(c, "Failed to count badges", err.Error(), "")
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
	response.SuccessOK(c, badges, "Badges retrieved successfully", pagination)
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
// @Success 200 {object} response.APIResponse{data=models.Badge} "Badge updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid badge update details or missing ID"
// @Failure 500 {object} response.APIResponse "Failed to update badge"
// @Router /badges/{id} [put]
func (h *BadgeHandler) UpdateBadge(c *gin.Context) {
	// Get badge ID from path parameter
	badgeID := c.Param("id")
	if badgeID == "" {
		response.BadRequest(c, "Badge ID is required", "Missing badge ID", "")
		return
	}

	var badgeUpdate models.BadgeCreate
	if err := c.ShouldBindJSON(&badgeUpdate); err != nil {
		h.logger.Error("Error binding badge update: %v", err)
		response.BadRequest(c, "Invalid badge update details", err.Error(), "")
		return
	}

	// Validate required fields
	if badgeUpdate.Name == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name": badgeUpdate.Name == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	err := h.service.UpdateBadge(c.Request.Context(), badgeID, &badgeUpdate)
	if err != nil {
		h.logger.Error("Failed to update badge: %v", err)
		response.InternalServerError(c, "Failed to update badge", err.Error(), "")
		return
	}

	response.SuccessOK(c, badgeUpdate, "Badge updated successfully")
}

// DeleteBadge godoc
// @Summary Delete a badge
// @Description Remove a badge from the system by its ID
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Badge ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Badge ID is required"
// @Failure 500 {object} response.APIResponse "Failed to delete badge"
// @Router /badges/{id} [delete]
func (h *BadgeHandler) DeleteBadge(c *gin.Context) {
	// Get badge ID from path parameter
	badgeID := c.Param("id")
	if badgeID == "" {
		response.BadRequest(c, "Badge ID is required", "Missing badge ID", "")
		return
	}

	err := h.service.DeleteBadge(c.Request.Context(), badgeID)
	if err != nil {
		h.logger.Error("Failed to delete badge: %v", err)
		response.InternalServerError(c, "Failed to delete badge", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Badge deleted successfully")
}

// CountBadges godoc
// @Summary Count badges
// @Description Retrieve the total number of badges in the system
// @Tags badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Success 200 {object} response.APIResponse{data=int} "Badge count retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to count badges"
// @Router /badges/count [get]
func (h *BadgeHandler) CountBadges(c *gin.Context) {
	// Optional filtering
	filters := []repository.FilterOption{}
	if name := c.Query("name"); name != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "name",
			Operator: "like",
			Value:    name,
		})
	}

	count, err := h.service.CountBadges(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to count badges: %v", err)
		response.InternalServerError(c, "Failed to count badges", err.Error(), "")
		return
	}

	response.SuccessOK(c, count, "Badge count retrieved successfully")
}
