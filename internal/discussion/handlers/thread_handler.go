package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
)

// ThreadHandler handles HTTP requests related to threads
type ThreadHandler struct {
	threadService services.ThreadService
	logger        *logger.Logger
}

// NewThreadHandler creates a new instance of thread handler
func NewThreadHandler(threadService services.ThreadService, logger *logger.Logger) *ThreadHandler {
	return &ThreadHandler{
		threadService: threadService,
		logger:        logger,
	}
}

// CreateThread godoc
// @Summary Create a new thread
// @Description Add a new discussion thread to the system
// @Tags Threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param thread body models.Thread true "Thread Information"
// @Success 201 {object} response.Response{data=models.Thread} "Thread created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid thread creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /threads [post]
func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.threadService.CreateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error creating thread: %v", err)
		response.InternalServerError(c, "Failed to create thread", err.Error())
		return
	}

	response.SuccessCreated(c, thread, "Thread created successfully")
}

// SearchThreads godoc
// @Summary Search threads
// @Description Search discussion threads by various criteria
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id query string false "Thread ID"
// @Param title query string false "Thread Title"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Thread} "Threads found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /threads/search [get]
func (h *ThreadHandler) SearchThread(c *gin.Context) {
	// Get query parameters
	id := c.Query("id")
	// title := c.Query("title")
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
		thread, err := h.threadService.GetThreadByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding thread by ID: %v", err)
			response.NotFound(c, "Thread not found", err.Error())
			return
		}
		response.SuccessOK(c, thread, "Thread found")
		return
	}

	// If no specific parameters are provided, return a list of threads
	threads, err := h.threadService.GetThreads(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving threads: %v", err)
		response.InternalServerError(c, "Failed to retrieve threads", err.Error())
		return
	}

	// Count total threads for pagination
	total := len(threads)

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, threads, total, offset/limit+1, limit)
}

// UpdateThread godoc
// @Summary Update a thread
// @Description Update an existing discussion thread's details
// @Tags Threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Thread ID"
// @Param thread body models.Thread true "Thread Update Details"
// @Success 200 {object} response.Response{data=models.Thread} "Thread updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid thread update details"
// @Failure 404 {object} response.ErrorResponse "Thread not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /threads/{id} [put]
func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.threadService.UpdateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error updating thread: %v", err)
		response.InternalServerError(c, "Failed to update thread", err.Error())
		return
	}

	response.SuccessOK(c, thread, "Thread updated successfully")
}

// ListThreads godoc
// @Summary List threads
// @Description Retrieve a list of discussion threads with pagination
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of threads to retrieve" default(10)
// @Param offset query int false "Number of threads to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Thread} "Threads retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list threads"
// @Router /threads [get]
func (h *ThreadHandler) ListThreads(c *gin.Context) {
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

	// Get list of threads
	threads, err := h.threadService.GetThreads(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving threads: %v", err)
		response.InternalServerError(c, "Failed to retrieve threads", err.Error())
		return
	}

	// Count total threads for pagination
	total := len(threads)

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, threads, total, offset/limit+1, limit)
}

// DeleteThread godoc
// @Summary Delete a thread
// @Description Remove a discussion thread from the system by its unique identifier
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Thread ID"
// @Success 200 {object} response.Response "Thread deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid thread ID"
// @Failure 404 {object} response.ErrorResponse "Thread not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /threads/{id} [delete]
func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Thread ID is required", nil)
		return
	}

	if err := h.threadService.DeleteThread(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting thread: %v", err)
		response.InternalServerError(c, "Failed to delete thread", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Thread deleted successfully")
}

// GetThreadByID godoc
// @Summary Get thread by ID
// @Description Retrieve a discussion thread's details by its unique identifier
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Thread ID"
// @Success 200 {object} response.Response{data=models.Thread} "Thread retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Thread not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /threads/{id} [get]
func (h *ThreadHandler) GetThreadByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Thread ID is required", nil)
		return
	}
	thread, err := h.threadService.GetThreadByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Error finding thread by ID: %v", err)
		response.NotFound(c, "Thread not found", err.Error())
		return
	}
	response.SuccessOK(c, thread, "Thread detail retrieved successfully")
}
