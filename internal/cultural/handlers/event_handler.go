package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/middleware"
	placeModel "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
)

// EventHandler handles HTTP requests related to events
type EventHandler struct {
	base.BaseHandler
	eventService services.EventService
}

// NewEventHandler creates a new instance of EventHandler
func NewEventHandler(eventService services.EventService, logger *logger.Logger) *EventHandler {
	return &EventHandler{
		BaseHandler:  *base.NewBaseHandler(logger),
		eventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create a new cultural event
// @Description Allows authenticated users to add a new cultural event to the platform
// @Description Supports multipart form data for event details and optional image upload
// @Tags Events
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param name formData string true "Event Name" minlength(2) maxlength(100)
// @Param description formData string true "Event Description" maxlength(500)
// @Param location formData string true "Location object as JSON (name, latitude, longitude)"
// @Param start_date formData string true "Start Date (RFC3339 or YYYY-MM-DD format)"
// @Param end_date formData string true "End Date (RFC3339 or YYYY-MM-DD format)"
// @Param is_kid_friendly formData bool false "Indicates if the event is suitable for children"
// @Param image formData file true "Event Cover Image (max 2MB)"
// @Success 201 {object} response.APIResponse{data=models.EventDTO} "Event successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid event creation details or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 500 {object} response.APIResponse "Internal server error during event creation"
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var eventInput models.EventPayload

	// Get user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrAuthentication,
			"Failed to get user ID",
			err,
		))
		return
	}
	eventInput.UserID = uuid.MustParse(userID)

	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrBadRequest,
			"Failed to parse multipart form",
			err,
		))
		return
	}

	// Extract and validate form fields
	if err := h.validateAndExtractEventInput(c, &eventInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get image file
	var fileHeader *multipart.FileHeader
	fileHeader, err = h.extractImageFile(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Validate image size if present
	if fileHeader != nil && fileHeader.Size > 2*1024*1024 {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Image file must be less than 2MB",
			nil,
		))
		return
	}

	// Create event via service
	event, err := h.eventService.CreateEvent(c.Request.Context(), &eventInput, fileHeader)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, event, "Event created successfully")
}

// validateAndExtractEventInput handles validation and extraction of event input fields
func (h *EventHandler) validateAndExtractEventInput(c *gin.Context, eventInput *models.EventPayload) error {
	// Extract required text fields
	eventInput.Name = c.PostForm("name")
	eventInput.Description = c.PostForm("description")
	eventInput.IsKidFriendly = c.PostForm("is_kid_friendly") == "true"

	// Parse and validate dates
	startDate, err := h.parseDateField(c, "start_date")
	if err != nil {
		return err
	}
	endDate, err := h.parseDateField(c, "end_date")
	if err != nil {
		return err
	}

	// Validate date order
	if endDate.Before(startDate) {
		return errors.New(
			errors.ErrValidation,
			"End date must be after or equal to start date",
			nil,
		)
	}

	eventInput.StartDate = startDate
	eventInput.EndDate = endDate

	// Parse and validate location
	locationStr := c.PostForm("location")
	if locationStr == "" {
		return errors.New(
			errors.ErrValidation,
			"Location field is required",
			nil,
		)
	}

	fmt.Println("locationStr:", locationStr)

	var loc placeModel.LocationCreate
	if err := json.Unmarshal([]byte(locationStr), &loc); err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid location format",
			err,
		)
	}
	eventInput.Location = &loc

	return nil
}

// parseDateField parses a date field with multiple format support
func (h *EventHandler) parseDateField(c *gin.Context, fieldName string) (time.Time, error) {
	dateStr := c.PostForm(fieldName)
	if dateStr == "" {
		return time.Time{}, errors.New(
			errors.ErrValidation,
			fmt.Sprintf("%s field is required", fieldName),
			nil,
		)
	}

	// Try parsing RFC3339 first
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If RFC3339 fails, try "2006-01-02" format
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return time.Time{}, errors.New(
				errors.ErrValidation,
				fmt.Sprintf("Invalid %s format: must be RFC3339 or YYYY-MM-DD", fieldName),
				err,
			)
		}
	}

	return date, nil
}

// extractImageFile handles extracting image file from request
func (h *EventHandler) extractImageFile(c *gin.Context) (*multipart.FileHeader, error) {
	form, _ := c.MultipartForm()
	var fileHeader *multipart.FileHeader

	if form != nil && len(form.File["image"]) > 0 {
		fileHeader = form.File["image"][0]
	} else {
		// fallback to single file
		f, err := c.FormFile("image")
		if err != nil {
			return nil, nil // No image is okay
		}
		fileHeader = f
	}

	return fileHeader, nil
}

// DeleteEvent godoc
// @Summary Delete an existing event
// @Description Allows event creator or administrator to remove an event from the platform
// @Description Permanently deletes the event and associated resources
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Event Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Event successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid event ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error during event deletion"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Event ID is required",
			nil,
		))
		return
	}

	err := h.eventService.DeleteEvent(c.Request.Context(), eventID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Event deleted successfully")
}

// GetEventByID godoc
// @Summary Retrieve a specific event
// @Description Fetches comprehensive details of an event by its unique identifier
// @Description Returns full event information including location, creator, and related metadata
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Event Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.EventDTO} "Successfully retrieved event details"
// @Failure 400 {object} response.APIResponse "Invalid event ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error during event retrieval"
// @Router /events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Event ID is required",
			nil,
		))
		return
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, event, "Event retrieved successfully")
}

// ListEvents godoc
// @Summary Retrieve events list
// @Description Fetches a paginated list of events with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of events per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort events by" default("created_at)" Enum(created_at,start_date,name)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param is_kid_friendly query bool false "Filter events by kid-friendliness"
// @Success 200 {object} response.APIResponse{data=[]models.EventDTO} "Successfully retrieved events list"
// @Success 204 {object} response.APIResponse "No events found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during events retrieval"
// @Router /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	opts, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid query parameters",
			err,
		))
		return
	}

	events, err := h.eventService.ListEvents(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, events, "Events retrieved successfully")
}

// UpdateEvent godoc
// @Summary Update an existing event
// @Description Allows event creator or administrator to modify event details
// @Description Supports partial updates with multipart form data and optional image upload
// @Tags Events
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Event Identifier" format(uuid)
// @Param name formData string false "Updated Event Name" minlength(2) maxlength(100)
// @Param description formData string false "Updated Event Description" maxlength(500)
// @Param location formData string false "Updated Location object as JSON"
// @Param start_date formData string false "Updated Start Date (RFC3339 or YYYY-MM-DD format)"
// @Param end_date formData string false "Updated End Date (RFC3339 or YYYY-MM-DD format)"
// @Param is_kid_friendly formData bool false "Updated kid-friendly status"
// @Param image formData file false "New Event Cover Image"
// @Success 200 {object} response.APIResponse{data=models.EventDTO} "Event successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid event update details or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error during event update"
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var eventInput models.EventPayload

	// Parse and validate input
	if err := h.validateAndExtractEventInput(c, &eventInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get image file
	var fileHeader *multipart.FileHeader
	fileHeader, err := h.extractImageFile(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Update event via service
	event, err := h.eventService.UpdateEvent(c.Request.Context(), &eventInput, fileHeader)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, event, "Event updated successfully")
}

// GetTrendingEvents godoc
// @Summary Retrieve trending events
// @Description Fetches a list of most popular or recently viewed events
// @Description Ranks events based on view count and recency
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Maximum number of trending events to retrieve" default(10) minimum(1) maximum(50)
// @Success 200 {object} response.APIResponse{data=[]models.EventDTO} "Successfully retrieved trending events"
// @Success 204 {object} response.APIResponse "No trending events found"
// @Failure 400 {object} response.APIResponse "Invalid limit parameter"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during trending events retrieval"
// @Router /events/trending [get]
func (h *EventHandler) GetTrendingEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid limit parameter",
			err,
		))
		return
	}

	events, err := h.eventService.GetTrendingEvents(c.Request.Context(), limit)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, events, "Trending events retrieved successfully")
}

// GetRelatedEvents godoc
// @Summary Retrieve related events
// @Description Finds events similar to a specific event based on location and other criteria
// @Description Helps users discover nearby or thematically connected events
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Reference Event Identifier" format(uuid)
// @Param location_id query string false "Optional location filter for related events"
// @Param limit query int false "Maximum number of related events to retrieve" default(5) minimum(1) maximum(20)
// @Success 200 {object} response.APIResponse{data=[]models.EventDTO} "Successfully retrieved related events"
// @Success 204 {object} response.APIResponse "No related events found"
// @Failure 400 {object} response.APIResponse "Invalid event ID or location ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during related events retrieval"
// @Router /events/{id}/related [get]
func (h *EventHandler) GetRelatedEvents(c *gin.Context) {
	eventID := c.Param("id")
	locationID := c.Query("location_id")
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid limit parameter",
			err,
		))
		return
	}

	events, err := h.eventService.GetRelatedEvents(c.Request.Context(), eventID, locationID, limit)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, events, "Related events retrieved successfully")
}

// SearchEvents godoc
// @Summary Search events
// @Description Performs a full-text search across event details with advanced filtering
// @Description Allows finding events by keywords, location, and other attributes
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding events" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,created_at,start_date)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.EventDTO} "Successfully completed event search"
// @Success 204 {object} response.APIResponse "No events match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during event search"
// @Router /events/search [get]
func (h *EventHandler) SearchEvents(c *gin.Context) {
	query := c.Query("query")
	opts, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid query parameters",
			err,
		))
		return
	}

	events, err := h.eventService.SearchEvents(c.Request.Context(), query, opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, events, "Events search completed successfully")
}

// UpdateEventViews godoc
// @Summary Increment event view count
// @Description Tracks and updates the number of times an event has been viewed by users
// @Description Helps in calculating event popularity and trending status
// @Tags Events
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Event Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Event views successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid event ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during view count update"
// @Router /events/{id}/views [post]
func (h *EventHandler) UpdateEventViews(c *gin.Context) {
	userID := c.GetString("user_id")
	eventID := c.Param("id")

	result := h.eventService.UpdateEventViews(c.Request.Context(), userID, eventID)
	h.HandleSuccess(c, result, "Event views updated")
}
