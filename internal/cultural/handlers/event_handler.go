package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// EventHandler handles HTTP requests related to events
type EventHandler struct {
	eventService services.EventService
	logger       *logger.Logger
}

// NewEventHandler creates a new instance of EventHandler
func NewEventHandler(eventService services.EventService, logger *logger.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Add a new cultural event to the system
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param event body models.Event true "Event Information"
// @Success 201 {object} response.APIResponse{data=models.ResponseEvent} "Event created successfully"
// @Failure 400 {object} response.APIResponse "Invalid event creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var eventInput struct {
		models.Event
		StartDate string `json:"start_date"` // format: "2006-01-02"
		StartTime string `json:"start_time"` // format: "15:04"
	}
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Combine start_date and start_time into timestamp
	var startTimestamp *time.Time
	if eventInput.StartDate != "" && eventInput.StartTime != "" {
		combined := eventInput.StartDate + "T" + eventInput.StartTime + ":00"
		parsed, err := time.Parse("2006-01-02T15:04:05", combined)
		if err != nil {
			h.logger.Error("Error parsing start_date and start_time: %v", err)
			response.BadRequest(c, "Invalid start_date or start_time format", err.Error(), "")
			return
		}
		startTimestamp = &parsed
	}
	event := eventInput.Event
	if startTimestamp != nil {
		event.StartDate = *startTimestamp
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Error creating event: %v", err)
		response.InternalServerError(c, "Failed to create event", err.Error(), "")
		return
	}

	response.SuccessCreated(c, event, "Event created successfully")
}

// SearchEvents godoc
// @Summary Search events
// @Description Search cultural events by various criteria
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (name, description, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.ResponseEvent} "Events found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /events/search [get]
func (h *EventHandler) SearchEvents(c *gin.Context) {
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
				Field:    "name",
				Operator: "like",
				Value:    query,
			},
			{
				Field:    "description",
				Operator: "like",
				Value:    query,
			},
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Search events
	events, err := h.eventService.SearchEvents(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching events: %v", err)
		response.InternalServerError(c, "Failed to search events", err.Error(), "")
		return
	}

	// Count total search results
	totalEvents, err := h.eventService.CountEvents(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalEvents,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalEvents) / float64(limit))),
		HasNextPage: offset+limit < totalEvents,
	}

	// Respond with events and pagination
	response.SuccessOK(c, events, "Events found successfully", pagination)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an existing cultural event's details
// @Tags Events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Event ID"
// @Param event body models.Event true "Event Update Details"
// @Success 200 {object} response.APIResponse{data=models.ResponseEvent} "Event updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid event update details"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	// Get event ID from path parameter
	eventID := c.Param("id")
	if eventID == "" {
		response.BadRequest(c, "Event ID is required", "Missing event ID", "")
		return
	}

	var eventInput struct {
		models.Event
		StartDate string `json:"start_date"` // format: "2006-01-02"
		StartTime string `json:"start_time"` // format: "15:04"
	}
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	var startTimestamp *time.Time
	if eventInput.StartDate != "" && eventInput.StartTime != "" {
		combined := eventInput.StartDate + "T" + eventInput.StartTime + ":00"
		parsed, err := time.Parse("2006-01-02T15:04:05", combined)
		if err != nil {
			h.logger.Error("Error parsing start_date and start_time: %v", err)
			response.BadRequest(c, "Invalid start_date or start_time format", err.Error(), "")
			return
		}
		startTimestamp = &parsed
	}
	event := eventInput.Event
	if startTimestamp != nil {
		event.StartDate = *startTimestamp
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(eventID)
	if err != nil {
		response.BadRequest(c, "Invalid Event ID", "Invalid UUID format", "")
		return
	}
	event.ID = parsedID

	if err := h.eventService.UpdateEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Error updating event: %v", err)
		response.InternalServerError(c, "Failed to update event", err.Error(), "")
		return
	}

	response.SuccessOK(c, event, "Event updated successfully")
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Remove a cultural event from the system by its unique identifier
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Event ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid event ID"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	// Get event ID from path parameter
	eventID := c.Param("id")
	if eventID == "" {
		response.BadRequest(c, "Event ID is required", "Missing event ID", "")
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), eventID); err != nil {
		h.logger.Error("Error deleting event: %v", err)
		response.InternalServerError(c, "Failed to delete event", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Event deleted successfully")
}

// ListEvent godoc
// @Summary List events
// @Description Retrieve a list of cultural events with pagination and filtering
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of events to retrieve" default(10)
// @Param offset query int false "Number of events to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.ResponseEvent} "Events retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list events"
// @Router /events [get]
func (h *EventHandler) ListEvent(c *gin.Context) {
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
	if isKidFriendly := c.Query("is_kid_friendly"); isKidFriendly != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "is_kid_friendly",
			Operator: "=",
			Value:    isKidFriendly,
		})
	}
	listOptions.Filters = filters

	// Retrieve events
	events, err := h.eventService.ListEvents(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving events: %v", err)
		response.InternalServerError(c, "Failed to retrieve events", err.Error(), "")
		return
	}

	// Count total events for pagination
	totalEvents, err := h.eventService.CountEvents(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting events: %v", err)
		response.InternalServerError(c, "Failed to count events", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalEvents,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalEvents) / float64(limit))),
		HasNextPage: offset+limit < totalEvents,
	}

	// Respond with events and pagination
	response.SuccessOK(c, events, "Events retrieved successfully", pagination)
}

// TrendingEvents godoc
// @Summary Get trending events
// @Description Retrieve a list of trending events based on views
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of trending events to retrieve" default(10)
// @Success 200 {object} response.APIResponse{data=[]models.ResponseEvent} "Trending events retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to retrieve trending events"
// @Router /events/trending [get]
func (h *EventHandler) TrendingEvents(c *gin.Context) {
	// Parse limit parameter
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	// Retrieve trending events
	events, err := h.eventService.GetTrendingEvents(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Error retrieving trending events: %v", err)
		response.InternalServerError(c, "Failed to retrieve trending events", err.Error(), "")
		return
	}

	response.SuccessOK(c, events, "Trending events retrieved successfully")
}

// GetEventByID godoc
// @Summary Get event by ID
// @Description Retrieve a cultural event's details by its unique identifier
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Event ID"
// @Success 200 {object} response.APIResponse{data=models.ResponseEvent} "Event retrieved successfully"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	// Get event ID from path parameter
	eventID := c.Param("id")
	if eventID == "" {
		response.BadRequest(c, "Event ID is required", "Missing event ID", "")
		return
	}

	// Retrieve event by ID
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		h.logger.Error("Error finding event by ID: %v", err)
		response.NotFound(c, "Event not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, event, "Event detail retrieved successfully")
}
