package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/holycann/cultour-backend/pkg/response"
)

// ThreadHandler handles HTTP requests related to threads
type ThreadHandler struct {
	base.BaseHandler
	threadService services.ThreadService
}

// NewThreadHandler creates a new instance of thread handler
func NewThreadHandler(
	threadService services.ThreadService,
	logger *logger.Logger,
) *ThreadHandler {
	return &ThreadHandler{
		BaseHandler:   *base.NewBaseHandler(logger),
		threadService: threadService,
	}
}

// CreateThread godoc
// @Summary Create a new discussion thread
// @Description Allows authenticated users to start a new discussion thread for a specific event
// @Description Supports creating threads with optional initial status
// @Tags Discussion Threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param thread body models.CreateThread true "Thread Creation Details"
// @Success 201 {object} response.APIResponse{data=models.ThreadDTO} "Thread successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid thread creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient event access privileges"
// @Failure 500 {object} response.APIResponse "Internal server error during thread creation"
// @Router /threads [post]
func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var thread models.CreateThread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.HandleError(c, errors.New(errors.ErrValidation, "Invalid request payload", err))
		return
	}

	// Get user context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Set creator ID if not provided
	if thread.CreatorID == uuid.Nil {
		thread.CreatorID, _ = uuid.Parse(userID)
	}

	// Validate thread
	if err := base.ValidateModel(thread); err != nil {
		h.HandleError(c, err)
		return
	}

	// Create thread
	threadData, err := h.threadService.CreateThread(c.Request.Context(), &thread)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessCreated(c, threadData, "Thread created successfully")
}

// SearchThreads godoc
// @Summary Search discussion threads
// @Description Performs a full-text search across thread details with advanced filtering
// @Description Allows finding threads by keywords, event, and other attributes
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding threads" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param status query string false "Filter threads by status" Enum(active,closed,archived)
// @Success 200 {object} response.APIResponse{data=[]models.ThreadDTO} "Successfully completed thread search"
// @Success 204 {object} response.APIResponse "No threads match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during thread search"
// @Router /threads/search [get]
func (h *ThreadHandler) SearchThreads(c *gin.Context) {
	// Prepare list options
	listOptions := &base.ListOptions{
		Page:      1,            // Default limit
		PerPage:   10,           // Default offset
		SortBy:    "created_at", // Default sort field
		SortOrder: "desc",       // Default sort order
	}

	// Parse limit from query parameter
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil {
			listOptions.PerPage = perPage
		}
	}

	// Parse offset from query parameter
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			listOptions.Page = page
		}
	}

	// Parse sort_by from query parameter
	if sortBy := c.Query("sort_by"); sortBy != "" {
		listOptions.SortBy = sortBy
	}

	// Parse sort_order from query parameter
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		listOptions.SortOrder = sortOrder
	}
	// Get search query
	query := c.Query("query")
	if query == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Search query is required",
			nil,
		))
		return
	}

	// Add title search filter
	listOptions.Filters = append(listOptions.Filters, base.FilterOption{
		Field:    "title",
		Operator: base.OperatorLike,
		Value:    query,
	})

	// Search threads
	threads, _, err := h.threadService.SearchThreads(c.Request.Context(), query, *listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Paginate results using base.PaginateResults
	paginatedThreads, pagination := base.PaginateResults(threads, listOptions.PerPage, listOptions.Page)

	// Respond with threads and pagination
	h.HandleSuccess(c, paginatedThreads, "Threads found successfully", response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}

// UpdateThread godoc
// @Summary Update an existing discussion thread
// @Description Allows thread creator or event administrator to modify thread details
// @Description Supports partial updates with thread status changes
// @Tags Discussion Threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Thread Identifier" format(uuid)
// @Param thread body models.CreateThread true "Thread Update Payload"
// @Success 200 {object} response.APIResponse{data=models.ThreadDTO} "Thread successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid thread update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient thread modification privileges"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error during thread update"
// @Router /threads/{id} [put]
func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	// Parse thread ID
	threadID, err := h.ValidateUUID("id", "Thread ID")
	if err != nil {
		return
	}

	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.HandleError(c, errors.New(errors.ErrValidation, "Invalid request payload", err))
		return
	}

	// Set the ID from path parameter
	thread.ID = threadID

	// Validate thread
	if err := base.ValidateModel(thread); err != nil {
		h.HandleError(c, err)
		return
	}

	// Update thread
	threadData, err := h.threadService.UpdateThread(c.Request.Context(), &thread)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, threadData, "Thread updated successfully")
}

// ListThreads godoc
// @Summary Retrieve discussion threads list
// @Description Fetches a paginated list of discussion threads with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of threads per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort threads by" default("created_at)" Enum(created_at,status)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param status query string false "Filter threads by status" Enum(active,closed,archived)
// @Param event_id query string false "Filter threads by specific event"
// @Success 200 {object} response.APIResponse{data=[]models.ThreadDTO} "Successfully retrieved threads list"
// @Success 204 {object} response.APIResponse "No threads found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during threads retrieval"
// @Router /threads [get]
func (h *ThreadHandler) ListThreads(c *gin.Context) {
	// Prepare list options with optional filters
	listOptions := &base.ListOptions{
		Page:      1,            // Default limit
		PerPage:   10,           // Default offset
		SortBy:    "created_at", // Default sort field
		SortOrder: "desc",       // Default sort order
	}

	// Optional filtering
	if status := c.Query("status"); status != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "status",
			Operator: base.OperatorEqual,
			Value:    status,
		})
	}
	if eventID := c.Query("event_id"); eventID != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "event_id",
			Operator: base.OperatorEqual,
			Value:    eventID,
		})
	}

	// Retrieve threads
	threads, err := h.threadService.ListThreads(c.Request.Context(), *listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	data, pagination := base.PaginateResults(threads, listOptions.Page, listOptions.PerPage)

	// Respond with threads and pagination
	h.HandleSuccess(c, data, "Threads retrieved successfully", response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}

// DeleteThread godoc
// @Summary Delete a discussion thread
// @Description Allows thread creator or event administrator to remove a specific thread
// @Description Permanently deletes the thread and associated messages
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Thread Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Thread successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid thread ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient thread deletion privileges"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error during thread deletion"
// @Router /threads/{id} [delete]
func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	// Parse thread ID
	threadID, err := h.ValidateUUID("id", "Thread ID")
	if err != nil {
		return
	}

	// Delete thread
	if err := h.threadService.DeleteThread(c.Request.Context(), threadID.String()); err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, nil, "Thread deleted successfully")
}

// GetThreadByID godoc
// @Summary Retrieve a specific discussion thread
// @Description Fetches comprehensive details of a thread by its unique identifier
// @Description Returns full thread information including creator, participants, and messages
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Thread Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.ThreadDTO} "Successfully retrieved thread details"
// @Failure 400 {object} response.APIResponse "Invalid thread ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error during thread retrieval"
// @Router /threads/{id} [get]
func (h *ThreadHandler) GetThreadByID(c *gin.Context) {
	// Parse thread ID
	threadID, err := h.ValidateUUID("id", "Thread ID")
	if err != nil {
		return
	}

	// Retrieve thread
	thread, err := h.threadService.GetThreadByID(c.Request.Context(), threadID.String())
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, thread, "Thread detail retrieved successfully")
}

// JoinThread godoc
// @Summary Join a discussion thread
// @Description Allows authenticated users to join an existing discussion thread
// @Description Adds the current user as a participant in the thread
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Thread Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.ThreadDTO} "Successfully joined thread"
// @Failure 400 {object} response.APIResponse "Invalid thread ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - thread joining not allowed"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 409 {object} response.APIResponse "User already a participant in the thread"
// @Failure 500 {object} response.APIResponse "Internal server error during thread joining"
// @Router /threads/{id}/join [post]
func (h *ThreadHandler) JoinThread(c *gin.Context) {
	// Parse thread ID from path parameter
	threadID := c.Param("id")
	if threadID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Thread ID is required or invalid",
			nil,
		))
		return
	}

	// Get user context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Join thread
	if err := h.threadService.JoinThread(c.Request.Context(), threadID, userID); err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, nil, "Successfully joined thread")
}

// GetThreadByEvent godoc
// @Summary Retrieve thread for a specific event
// @Description Fetches the discussion thread associated with a particular event
// @Description Returns the primary or most recent thread for the given event
// @Tags Discussion Threads
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param event_id path string true "Unique Event Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.ThreadDTO} "Successfully retrieved event thread"
// @Failure 400 {object} response.APIResponse "Invalid event ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "No thread found for the specified event"
// @Failure 500 {object} response.APIResponse "Internal server error during event thread retrieval"
// @Router /threads/event/{event_id} [get]
func (h *ThreadHandler) GetThreadByEvent(c *gin.Context) {
	eventID := c.Param("event_id")
	if eventID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Event ID is required or invalid",
			nil,
		))
		return
	}

	// Retrieve thread by event
	thread, err := h.threadService.GetThreadByEvent(c.Request.Context(), eventID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, thread, "Thread retrieved successfully")
}
