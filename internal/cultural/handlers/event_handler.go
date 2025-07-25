package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
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
// @Success 201 {object} response.Response{data=models.ResponseEvent} "Event created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid event creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var eventInput struct {
		models.Event
		StartDate string `json:"start_date"` // format: "2006-01-02"
		StartTime string `json:"start_time"` // format: "15:04"
	}
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	// Combine start_date and start_time into timestamp
	var startTimestamp *time.Time
	if eventInput.StartDate != "" && eventInput.StartTime != "" {
		combined := eventInput.StartDate + "T" + eventInput.StartTime + ":00"
		parsed, err := time.Parse("2006-01-02T15:04:05", combined)
		if err != nil {
			h.logger.Error("Error parsing start_date and start_time: %v", err)
			response.BadRequest(c, "Invalid start_date or start_time format", err.Error())
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
		response.InternalServerError(c, "Failed to create event", err.Error())
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
// @Param id query string false "Event ID"
// @Param name query string false "Event Name"
// @Param query query string false "Search query"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseEvent} "Events found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /events/search [get]
func (h *EventHandler) SearchEvents(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	query := c.Query("query")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

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

	if query != "" {
		events, err := h.eventService.SearchEvents(c.Request.Context(), query, limit, offset)
		if err != nil {
			h.logger.Error("Error searching events: %v", err)
			response.InternalServerError(c, "Failed to search events", err.Error())
			return
		}
		total := len(events)
		response.WithPagination(c, events, total, offset/limit+1, limit)
		return
	}

	if id != "" {
		event, err := h.eventService.GetEventByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding event by ID: %v", err)
			response.NotFound(c, "Event not found", err.Error())
			return
		}
		response.SuccessOK(c, event, "Event found")
		return
	}

	if name != "" {
		event, err := h.eventService.GetEventByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding event by name: %v", err)
			response.NotFound(c, "Event not found", err.Error())
			return
		}
		response.SuccessOK(c, event, "Event found")
		return
	}

	events, err := h.eventService.GetEvents(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving events: %v", err)
		response.InternalServerError(c, "Failed to retrieve events", err.Error())
		return
	}

	total, err := h.eventService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting events: %v", err)
		response.InternalServerError(c, "Failed to count events", err.Error())
		return
	}

	response.WithPagination(c, events, total, offset/limit+1, limit)
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
// @Success 200 {object} response.Response{data=models.ResponseEvent} "Event updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid event update details"
// @Failure 404 {object} response.ErrorResponse "Event not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var eventInput struct {
		models.Event
		StartDate string `json:"start_date"` // format: "2006-01-02"
		StartTime string `json:"start_time"` // format: "15:04"
	}
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	var startTimestamp *time.Time
	if eventInput.StartDate != "" && eventInput.StartTime != "" {
		combined := eventInput.StartDate + "T" + eventInput.StartTime + ":00"
		parsed, err := time.Parse("2006-01-02T15:04:05", combined)
		if err != nil {
			h.logger.Error("Error parsing start_date and start_time: %v", err)
			response.BadRequest(c, "Invalid start_date or start_time format", err.Error())
			return
		}
		startTimestamp = &parsed
	}
	event := eventInput.Event
	if startTimestamp != nil {
		event.StartDate = *startTimestamp
	}

	if err := h.eventService.UpdateEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Error updating event: %v", err)
		response.InternalServerError(c, "Failed to update event", err.Error())
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
// @Success 200 {object} response.Response "Event deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid event ID"
// @Failure 404 {object} response.ErrorResponse "Event not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Event ID is required", nil)
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting event: %v", err)
		response.InternalServerError(c, "Failed to delete event", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Event deleted successfully")
}

// ListEvent godoc
// @Summary List events
// @Description Retrieve a list of cultural events with pagination
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of events to retrieve" default(10)
// @Param offset query int false "Number of events to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseEvent} "Events retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list events"
// @Router /events [get]
func (h *EventHandler) ListEvent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

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

	events, err := h.eventService.GetEvents(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving events: %v", err)
		response.InternalServerError(c, "Failed to retrieve events", err.Error())
		return
	}

	total, err := h.eventService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting events: %v", err)
		response.InternalServerError(c, "Failed to count events", err.Error())
		return
	}

	response.WithPagination(c, events, total, offset/limit+1, limit)
}

// TrendingEvents godoc
// @Summary Get trending events
// @Description Retrieve a list of trending events based on views
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of trending events to retrieve" default(10)
// @Success 200 {object} response.Response{data=[]models.ResponseEvent} "Trending events retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trending events"
// @Router /events/trending [get]
func (h *EventHandler) TrendingEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	events, err := h.eventService.GetTrendingEvents(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Error retrieving trending events: %v", err)
		response.InternalServerError(c, "Failed to retrieve trending events", err.Error())
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
// @Success 200 {object} response.Response{data=models.ResponseEvent} "Event retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Event not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Event ID is required", nil)
		return
	}
	event, err := h.eventService.GetEventByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Error finding event by ID: %v", err)
		response.NotFound(c, "Event not found", err.Error())
		return
	}
	response.SuccessOK(c, event, "Event detail retrieved successfully")
}
