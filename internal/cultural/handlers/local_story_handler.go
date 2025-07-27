package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// LocalStoryHandler handles HTTP requests related to local stories
type LocalStoryHandler struct {
	localStoryService services.LocalStoryService
	logger            *logger.Logger
}

// NewLocalStoryHandler creates a new instance of local story handler
func NewLocalStoryHandler(localStoryService services.LocalStoryService, logger *logger.Logger) *LocalStoryHandler {
	return &LocalStoryHandler{
		localStoryService: localStoryService,
		logger:            logger,
	}
}

// CreateLocalStory godoc
// @Summary Create a new local story
// @Description Add a new local cultural story to the system
// @Tags Local Stories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param local_story body models.LocalStory true "Local Story Information"
// @Success 201 {object} response.APIResponse{data=models.LocalStory} "Local story created successfully"
// @Failure 400 {object} response.APIResponse "Invalid local story creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /local-stories [post]
func (h *LocalStoryHandler) CreateLocalStory(c *gin.Context) {
	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if localStory.Title == "" || localStory.StoryText == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"title":      localStory.Title == "",
			"story_text": localStory.StoryText == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.localStoryService.CreateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error creating local story: %v", err)
		response.InternalServerError(c, "Failed to create local story", err.Error(), "")
		return
	}

	response.SuccessCreated(c, localStory, "Local story created successfully")
}

// SearchLocalStories godoc
// @Summary Search local stories
// @Description Search local cultural stories by various criteria
// @Tags Local Stories
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (title, story text, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.LocalStory} "Local stories found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /local-stories/search [get]
func (h *LocalStoryHandler) SearchLocalStories(c *gin.Context) {
	// Get search query
	query := c.Query("query")
	if query == "" {
		response.BadRequest(c, "Search query is required", "Empty search query", "")
		return
	}

	// Parse pagination parameters
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

	// Prepare list options for search
	listOptions := repository.ListOptions{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: repository.SortDescending,
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Search local stories
	localStories, err := h.localStoryService.SearchLocalStories(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching local stories: %v", err)
		response.InternalServerError(c, "Failed to search local stories", err.Error(), "")
		return
	}

	// Count total search results
	totalStories, err := h.localStoryService.CountLocalStories(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalStories,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalStories) / float64(limit))),
		HasNextPage: offset+limit < totalStories,
	}

	// Respond with local stories and pagination
	response.SuccessOK(c, localStories, "Local stories found successfully", pagination)
}

// UpdateLocalStory godoc
// @Summary Update a local story
// @Description Update an existing local cultural story's details
// @Tags Local Stories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Local Story ID"
// @Param local_story body models.LocalStory true "Local Story Update Details"
// @Success 200 {object} response.APIResponse{data=models.LocalStory} "Local story updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid local story update details"
// @Failure 404 {object} response.APIResponse "Local story not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /local-stories/{id} [put]
func (h *LocalStoryHandler) UpdateLocalStory(c *gin.Context) {
	// Get local story ID from path parameter
	storyID := c.Param("id")
	if storyID == "" {
		response.BadRequest(c, "Local Story ID is required", "Missing local story ID", "")
		return
	}

	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(storyID)
	if err != nil {
		response.BadRequest(c, "Invalid Local Story ID", "Invalid UUID format", "")
		return
	}
	localStory.ID = parsedID

	if err := h.localStoryService.UpdateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error updating local story: %v", err)
		response.InternalServerError(c, "Failed to update local story", err.Error(), "")
		return
	}

	response.SuccessOK(c, localStory, "Local story updated successfully")
}

// DeleteLocalStory godoc
// @Summary Delete a local story
// @Description Remove a local cultural story from the system by its unique identifier
// @Tags Local Stories
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Local Story ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid local story ID"
// @Failure 404 {object} response.APIResponse "Local story not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /local-stories/{id} [delete]
func (h *LocalStoryHandler) DeleteLocalStory(c *gin.Context) {
	// Get local story ID from path parameter
	storyID := c.Param("id")
	if storyID == "" {
		response.BadRequest(c, "Local Story ID is required", "Missing local story ID", "")
		return
	}

	if err := h.localStoryService.DeleteLocalStory(c.Request.Context(), storyID); err != nil {
		h.logger.Error("Error deleting local story: %v", err)
		response.InternalServerError(c, "Failed to delete local story", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Local story deleted successfully")
}

// ListLocalStories godoc
// @Summary List local stories
// @Description Retrieve a list of local cultural stories with pagination and filtering
// @Tags Local Stories
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of local stories to retrieve" default(10)
// @Param offset query int false "Number of local stories to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.LocalStory} "Local stories retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list local stories"
// @Router /local-stories [get]
func (h *LocalStoryHandler) ListLocalStories(c *gin.Context) {
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
	if isForKids := c.Query("is_for_kids"); isForKids != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "is_for_kids",
			Operator: "=",
			Value:    isForKids,
		})
	}
	if originCulture := c.Query("origin_culture"); originCulture != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "origin_culture",
			Operator: "=",
			Value:    originCulture,
		})
	}
	listOptions.Filters = filters

	// Retrieve local stories
	localStories, err := h.localStoryService.ListLocalStories(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving local stories: %v", err)
		response.InternalServerError(c, "Failed to retrieve local stories", err.Error(), "")
		return
	}

	// Count total local stories for pagination
	totalStories, err := h.localStoryService.CountLocalStories(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting local stories: %v", err)
		response.InternalServerError(c, "Failed to count local stories", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalStories,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalStories) / float64(limit))),
		HasNextPage: offset+limit < totalStories,
	}

	// Respond with local stories and pagination
	response.SuccessOK(c, localStories, "Local stories retrieved successfully", pagination)
}
