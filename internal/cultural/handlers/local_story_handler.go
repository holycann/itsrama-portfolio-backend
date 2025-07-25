package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
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
// @Success 201 {object} response.Response{data=models.LocalStory} "Local story created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid local story creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /local-stories [post]
func (h *LocalStoryHandler) CreateLocalStory(c *gin.Context) {
	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.localStoryService.CreateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error creating local story: %v", err)
		response.InternalServerError(c, "Failed to create local story", err.Error())
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
// @Param id query string false "Local Story ID"
// @Param title query string false "Local Story Title"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.LocalStory} "Local stories found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /local-stories/search [get]
func (h *LocalStoryHandler) SearchLocalStories(c *gin.Context) {
	// Get query parameters
	id := c.Query("id")
	title := c.Query("title")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parse limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// If ID is provided, search by ID
	if id != "" {
		localStory, err := h.localStoryService.GetLocalStoryByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding local story by ID: %v", err)
			response.NotFound(c, "Local story not found", err.Error())
			return
		}
		response.SuccessOK(c, localStory, "Local story found")
		return
	}

	// If title is provided, search by title
	if title != "" {
		localStory, err := h.localStoryService.GetLocalStoryByTitle(c.Request.Context(), title)
		if err != nil {
			h.logger.Error("Error finding local story by title: %v", err)
			response.NotFound(c, "Local story not found", err.Error())
			return
		}
		response.SuccessOK(c, localStory, "Local story found")
		return
	}

	// If no specific parameters are provided, return a list of local stories
	localStories, err := h.localStoryService.GetLocalStories(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving local stories: %v", err)
		response.InternalServerError(c, "Failed to retrieve local stories", err.Error())
		return
	}

	// Count total local stories for pagination
	total, err := h.localStoryService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting local stories: %v", err)
		response.InternalServerError(c, "Failed to count local stories", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, localStories, total, offset/limit+1, limit)
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
// @Success 200 {object} response.Response{data=models.LocalStory} "Local story updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid local story update details"
// @Failure 404 {object} response.ErrorResponse "Local story not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /local-stories/{id} [put]
func (h *LocalStoryHandler) UpdateLocalStory(c *gin.Context) {
	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.localStoryService.UpdateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error updating local story: %v", err)
		response.InternalServerError(c, "Failed to update local story", err.Error())
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
// @Success 200 {object} response.Response "Local story deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid local story ID"
// @Failure 404 {object} response.ErrorResponse "Local story not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /local-stories/{id} [delete]
func (h *LocalStoryHandler) DeleteLocalStory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Local story ID is required", nil)
		return
	}

	if err := h.localStoryService.DeleteLocalStory(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting local story: %v", err)
		response.InternalServerError(c, "Failed to delete local story", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Local story deleted successfully")
}

// ListLocalStories godoc
// @Summary List local stories
// @Description Retrieve a list of local cultural stories with pagination
// @Tags Local Stories
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of local stories to retrieve" default(10)
// @Param offset query int false "Number of local stories to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.LocalStory} "Local stories retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list local stories"
// @Router /local-stories [get]
func (h *LocalStoryHandler) ListLocalStories(c *gin.Context) {
	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parse limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// Get list of local stories
	localStories, err := h.localStoryService.GetLocalStories(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving local stories: %v", err)
		response.InternalServerError(c, "Failed to retrieve local stories", err.Error())
		return
	}

	// Count total local stories for pagination
	total, err := h.localStoryService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting local stories: %v", err)
		response.InternalServerError(c, "Failed to count local stories", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, localStories, total, offset/limit+1, limit)
}
