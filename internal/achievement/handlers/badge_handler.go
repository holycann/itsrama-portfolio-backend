package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
	"github.com/holycann/cultour-backend/pkg/validator"
)

// BadgeHandler handles HTTP requests related to badges
// @Description Manages badge-related operations such as creation, retrieval, update, and deletion
type BadgeHandler struct {
	base.BaseHandler
	service services.BadgeService
}

// NewBadgeHandler creates a new instance of BadgeHandler
// @Description Initializes a new BadgeHandler with the provided BadgeService and logger
func NewBadgeHandler(service services.BadgeService, logger *logger.Logger) *BadgeHandler {
	return &BadgeHandler{
		BaseHandler: *base.NewBaseHandler(logger),
		service:     service,
	}
}

// CreateBadge godoc
// @Summary Create a new badge in the system
// @Description Allows administrators to create a new achievement badge with detailed metadata
// @Description Requires admin authentication and authorization
// @Tags Badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param badge body models.BadgeCreate true "Badge Creation Details"
// @Success 201 {object} response.APIResponse{data=models.BadgeDTO} "Badge successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid badge creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges (admin role required)"
// @Failure 500 {object} response.APIResponse "Internal server error during badge creation"
// @Router /badges [post]
func (h *BadgeHandler) CreateBadge(c *gin.Context) {
	var badgeCreate models.BadgeCreate

	// Validate request body
	if err := c.ShouldBindJSON(&badgeCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	// Validate badge creation payload
	if err := validator.ValidateStruct(badgeCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Create badge through service
	createdBadge, err := h.service.CreateBadge(c.Request.Context(), &badgeCreate)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to create badge"))
		return
	}

	// Respond with created badge
	h.HandleCreated(c, createdBadge, "Badge created successfully")
}

// ListBadges godoc
// @Summary Retrieve a list of badges
// @Description Fetches a paginated list of badges with optional filtering and sorting
// @Description Supports pagination, sorting, and name-based filtering
// @Tags Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Maximum number of badges to return" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of badges to skip for pagination" default(0) minimum(0)
// @Param sort_by query string false "Field to sort badges by" default("created_at)" Enum(created_at,name,id)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param name query string false "Filter badges by partial name match"
// @Success 200 {object} response.APIResponse{data=[]models.BadgeDTO} "Successfully retrieved badge list"
// @Success 204 {object} response.APIResponse "No badges found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during badge retrieval"
// @Router /badges [get]
func (h *BadgeHandler) ListBadges(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Prepare list options
	listOptions := base.ListOptions{
		Page:      offset/limit + 1,
		PerPage:   limit,
		SortBy:    sortBy,
		SortOrder: base.SortDescending,
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = base.SortAscending
	}

	// Optional filtering
	if name := c.Query("name"); name != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "name",
			Operator: base.OperatorLike,
			Value:    name,
		})
	}

	// Retrieve badges
	badges, totalBadges, err := h.service.ListBadges(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to retrieve badges"))
		return
	}

	// Respond with badges and pagination
	h.HandlePagination(c, badges, totalBadges, listOptions)
}

// SearchBadges godoc
// @Summary Search badges by query
// @Description Performs a full-text search across badge name, description, and other relevant fields
// @Description Supports advanced search with pagination and relevance ranking
// @Tags Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding badges" minlength(2)
// @Param limit query int false "Maximum number of search results" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Success 200 {object} response.APIResponse{data=[]models.BadgeDTO} "Successfully found matching badges"
// @Success 204 {object} response.APIResponse "No badges match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during badge search"
// @Router /badges/search [get]
func (h *BadgeHandler) SearchBadges(c *gin.Context) {
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

	// Search badges
	badges, totalBadges, err := h.service.SearchBadges(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to search badges"))
		return
	}

	// Respond with badges and pagination
	h.HandlePagination(c, badges, totalBadges, listOptions)
}

// UpdateBadge godoc
// @Summary Update an existing badge
// @Description Allows administrators to modify badge details by its unique identifier
// @Description Supports partial updates with optional fields
// @Tags Badges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Badge Identifier" format(uuid)
// @Param badge body models.BadgeUpdate true "Badge Update Payload"
// @Success 200 {object} response.APIResponse{data=models.BadgeDTO} "Badge successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid badge update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Badge not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge update"
// @Router /badges/{id} [put]
func (h *BadgeHandler) UpdateBadge(c *gin.Context) {
	// Get badge ID from path parameter
	badgeIDStr := c.Param("id")
	badgeID, err := h.ValidateUUID(badgeIDStr, "badge_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid badge ID"))
		return
	}

	// Create a badge update model to bind request body
	var updateBadge models.BadgeUpdate

	// Bind and validate input
	if err := c.ShouldBindJSON(&updateBadge); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid request body"))
		return
	}

	// Validate badge update payload
	if err := validator.ValidateStruct(updateBadge); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Validation failed"))
		return
	}

	// Set the ID from path parameter
	updateBadge.ID = badgeID

	// Update badge
	updatedBadge, err := h.service.UpdateBadge(c.Request.Context(), badgeIDStr, &updateBadge)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update badge"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, updatedBadge, "Badge updated successfully")
}

// DeleteBadge godoc
// @Summary Delete a badge
// @Description Permanently removes a badge from the system by its unique identifier
// @Description Requires administrative privileges
// @Tags Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Badge Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Badge successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid badge ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Badge not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge deletion"
// @Router /badges/{id} [delete]
func (h *BadgeHandler) DeleteBadge(c *gin.Context) {
	// Get badge ID from path parameter
	badgeIDStr := c.Param("id")
	badgeID, err := h.ValidateUUID(badgeIDStr, "badge_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid badge ID"))
		return
	}

	// Delete badge
	if err := h.service.DeleteBadge(c.Request.Context(), badgeIDStr); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to delete badge"))
		return
	}

	// Respond with success
	h.HandleSuccess(c, gin.H{
		"id": badgeID,
	}, "Badge deleted successfully")
}

// GetBadgeByID godoc
// @Summary Retrieve a specific badge
// @Description Fetches detailed information about a badge by its unique identifier
// @Tags Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Badge Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.BadgeDTO} "Successfully retrieved badge details"
// @Failure 400 {object} response.APIResponse "Invalid badge ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "Badge not found"
// @Failure 500 {object} response.APIResponse "Internal server error during badge retrieval"
// @Router /badges/{id} [get]
func (h *BadgeHandler) GetBadgeByID(c *gin.Context) {
	// Get badge ID from path parameter
	badgeIDStr := c.Param("id")
	_, err := h.ValidateUUID(badgeIDStr, "badge_id")
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid badge ID"))
		return
	}

	// Retrieve badge
	badge, err := h.service.GetBadgeByID(c.Request.Context(), badgeIDStr)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Badge not found"))
		return
	}

	// Respond with badge details
	h.HandleSuccess(c, badge, "Badge retrieved successfully")
}

// CountBadges godoc
// @Summary Count total number of badges
// @Description Returns the total count of badges in the system, with optional name filtering
// @Tags Badges
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param name query string false "Optional filter to count badges by name"
// @Success 200 {object} response.APIResponse{data=int} "Successfully retrieved badge count"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during badge count"
// @Router /badges/count [get]
func (h *BadgeHandler) CountBadges(c *gin.Context) {
	// Optional filtering
	filters := []base.FilterOption{}
	if name := c.Query("name"); name != "" {
		filters = append(filters, base.FilterOption{
			Field:    "name",
			Operator: base.OperatorLike,
			Value:    name,
		})
	}

	count, err := h.service.CountBadges(c.Request.Context(), filters)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to count badges"))
		return
	}

	h.HandleSuccess(c, count, "Badge count retrieved successfully")
}
