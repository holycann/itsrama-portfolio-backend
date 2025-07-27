package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
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
// @Success 201 {object} response.APIResponse{data=models.Thread} "Thread created successfully"
// @Failure 400 {object} response.APIResponse "Invalid thread creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /threads [post]
func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if thread.Title == "" || thread.EventID == uuid.Nil {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"title":    thread.Title == "",
			"event_id": thread.EventID == uuid.Nil,
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.threadService.CreateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error creating thread: %v", err)
		response.InternalServerError(c, "Failed to create thread", err.Error(), "")
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
// @Param query query string true "Search query (title, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.Thread} "Threads found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /threads/search [get]
func (h *ThreadHandler) SearchThreads(c *gin.Context) {
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
		Filters: []repository.FilterOption{
			{
				Field:    "title",
				Operator: "like",
				Value:    query,
			},
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Search threads
	threads, err := h.threadService.SearchThreads(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching threads: %v", err)
		response.InternalServerError(c, "Failed to search threads", err.Error(), "")
		return
	}

	// Count total search results
	totalThreads, err := h.threadService.CountThreads(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalThreads,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalThreads) / float64(limit))),
		HasNextPage: offset+limit < totalThreads,
	}

	// Respond with threads and pagination
	response.SuccessOK(c, threads, "Threads found successfully", pagination)
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
// @Success 200 {object} response.APIResponse{data=models.Thread} "Thread updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid thread update details"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /threads/{id} [put]
func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	// Get thread ID from path parameter
	threadID := c.Param("id")
	if threadID == "" {
		response.BadRequest(c, "Thread ID is required", "Missing thread ID", "")
		return
	}

	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(threadID)
	if err != nil {
		response.BadRequest(c, "Invalid Thread ID", "Invalid UUID format", "")
		return
	}
	thread.ID = parsedID

	if err := h.threadService.UpdateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error updating thread: %v", err)
		response.InternalServerError(c, "Failed to update thread", err.Error(), "")
		return
	}

	response.SuccessOK(c, thread, "Thread updated successfully")
}

// ListThreads godoc
// @Summary List threads
// @Description Retrieve a list of discussion threads with pagination and filtering
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of threads to retrieve" default(10)
// @Param offset query int false "Number of threads to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.Thread} "Threads retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list threads"
// @Router /threads [get]
func (h *ThreadHandler) ListThreads(c *gin.Context) {
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
	if status := c.Query("status"); status != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "status",
			Operator: "=",
			Value:    status,
		})
	}
	if eventID := c.Query("event_id"); eventID != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "event_id",
			Operator: "=",
			Value:    eventID,
		})
	}
	listOptions.Filters = filters

	// Retrieve threads
	threads, err := h.threadService.ListThreads(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving threads: %v", err)
		response.InternalServerError(c, "Failed to retrieve threads", err.Error(), "")
		return
	}

	// Count total threads for pagination
	totalThreads, err := h.threadService.CountThreads(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting threads: %v", err)
		response.InternalServerError(c, "Failed to count threads", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalThreads,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalThreads) / float64(limit))),
		HasNextPage: offset+limit < totalThreads,
	}

	// Respond with threads and pagination
	response.SuccessOK(c, threads, "Threads retrieved successfully", pagination)
}

// DeleteThread godoc
// @Summary Delete a thread
// @Description Remove a discussion thread from the system by its unique identifier
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Thread ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid thread ID"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /threads/{id} [delete]
func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	// Get thread ID from path parameter
	threadID := c.Param("id")
	if threadID == "" {
		response.BadRequest(c, "Thread ID is required", "Missing thread ID", "")
		return
	}

	if err := h.threadService.DeleteThread(c.Request.Context(), threadID); err != nil {
		h.logger.Error("Error deleting thread: %v", err)
		response.InternalServerError(c, "Failed to delete thread", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Thread deleted successfully")
}

// GetThreadByID godoc
// @Summary Get thread by ID
// @Description Retrieve a discussion thread's details by its unique identifier
// @Tags Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Thread ID"
// @Success 200 {object} response.APIResponse{data=models.Thread} "Thread retrieved successfully"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /threads/{id} [get]
func (h *ThreadHandler) GetThreadByID(c *gin.Context) {
	// Get thread ID from path parameter
	threadID := c.Param("id")
	if threadID == "" {
		response.BadRequest(c, "Thread ID is required", "Missing thread ID", "")
		return
	}

	thread, err := h.threadService.GetThreadByID(c.Request.Context(), threadID)
	if err != nil {
		h.logger.Error("Error finding thread by ID: %v", err)
		response.NotFound(c, "Thread not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, thread, "Thread detail retrieved successfully")
}
