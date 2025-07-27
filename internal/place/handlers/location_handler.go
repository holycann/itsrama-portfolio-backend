package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// LocationHandler handles HTTP requests related to locations
type LocationHandler struct {
	locationService services.LocationService
	logger          *logger.Logger
}

// NewLocationHandler creates a new instance of location handler
func NewLocationHandler(locationService services.LocationService, logger *logger.Logger) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
		logger:          logger,
	}
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Add a new location to the system
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param location body models.Location true "Location Information"
// @Success 201 {object} response.APIResponse{data=models.Location} "Location created successfully"
// @Failure 400 {object} response.APIResponse "Invalid location creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if location.Name == "" || location.CityID == uuid.Nil || location.Latitude == 0 || location.Longitude == 0 {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name":      location.Name == "",
			"city_id":   location.CityID == uuid.Nil,
			"latitude":  location.Latitude == 0,
			"longitude": location.Longitude == 0,
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.locationService.CreateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error creating location: %v", err)
		response.InternalServerError(c, "Failed to create location", err.Error(), "")
		return
	}

	response.SuccessCreated(c, location, "Location created successfully")
}

// SearchLocations godoc
// @Summary Search locations
// @Description Search locations by various criteria
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (name, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Param city_id query string false "Filter by city ID"
// @Success 200 {object} response.APIResponse{data=[]models.Location} "Locations found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
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
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Optional filtering by city_id
	if cityID := c.Query("city_id"); cityID != "" {
		listOptions.Filters = append(listOptions.Filters, repository.FilterOption{
			Field:    "city_id",
			Operator: "eq",
			Value:    cityID,
		})
	}

	// Search locations
	locations, err := h.locationService.SearchLocations(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching locations: %v", err)
		response.InternalServerError(c, "Failed to search locations", err.Error(), "")
		return
	}

	// Count total search results
	totalLocations, err := h.locationService.CountLocations(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalLocations,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalLocations) / float64(limit))),
		HasNextPage: offset+limit < totalLocations,
	}

	// Respond with locations and pagination
	response.SuccessOK(c, locations, "Locations found successfully", pagination)
}

// UpdateLocation godoc
// @Summary Update a location
// @Description Update an existing location's details
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Location ID"
// @Param location body models.Location true "Location Update Details"
// @Success 200 {object} response.APIResponse{data=models.Location} "Location updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid location update details"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	// Get location ID from path parameter
	locationID := c.Param("id")
	if locationID == "" {
		response.BadRequest(c, "Location ID is required", "Missing location ID", "")
		return
	}

	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(locationID)
	if err != nil {
		response.BadRequest(c, "Invalid Location ID", "Invalid UUID format", "")
		return
	}
	location.ID = parsedID

	// Validate required fields
	if location.Name == "" || location.CityID == uuid.Nil || location.Latitude == 0 || location.Longitude == 0 {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name":      location.Name == "",
			"city_id":   location.CityID == uuid.Nil,
			"latitude":  location.Latitude == 0,
			"longitude": location.Longitude == 0,
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.locationService.UpdateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error updating location: %v", err)
		response.InternalServerError(c, "Failed to update location", err.Error(), "")
		return
	}

	response.SuccessOK(c, location, "Location updated successfully")
}

// DeleteLocation godoc
// @Summary Delete a location
// @Description Remove a location from the system by its unique identifier
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Location ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid location ID"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	// Get location ID from path parameter
	locationID := c.Param("id")
	if locationID == "" {
		response.BadRequest(c, "Location ID is required", "Missing location ID", "")
		return
	}

	if err := h.locationService.DeleteLocation(c.Request.Context(), locationID); err != nil {
		h.logger.Error("Error deleting location: %v", err)
		response.InternalServerError(c, "Failed to delete location", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Location deleted successfully")
}

// ListLocations godoc
// @Summary List locations
// @Description Retrieve a list of locations with pagination and filtering
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of locations to retrieve" default(10)
// @Param offset query int false "Number of locations to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Param city_id query string false "Filter by city ID"
// @Success 200 {object} response.APIResponse{data=[]models.Location} "Locations retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list locations"
// @Router /locations [get]
func (h *LocationHandler) ListLocations(c *gin.Context) {
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
	if cityID := c.Query("city_id"); cityID != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "city_id",
			Operator: "eq",
			Value:    cityID,
		})
	}
	listOptions.Filters = filters

	// Retrieve locations
	locations, err := h.locationService.ListLocations(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving locations: %v", err)
		response.InternalServerError(c, "Failed to retrieve locations", err.Error(), "")
		return
	}

	// Count total locations for pagination
	totalLocations, err := h.locationService.CountLocations(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting locations: %v", err)
		response.InternalServerError(c, "Failed to count locations", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalLocations,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalLocations) / float64(limit))),
		HasNextPage: offset+limit < totalLocations,
	}

	// Respond with locations and pagination
	response.SuccessOK(c, locations, "Locations retrieved successfully", pagination)
}

// GetLocationByID godoc
// @Summary Get location by ID
// @Description Retrieve a location's details by its unique identifier
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Location ID"
// @Success 200 {object} response.APIResponse{data=models.Location} "Location retrieved successfully"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /locations/{id} [get]
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	// Get location ID from path parameter
	locationID := c.Param("id")
	if locationID == "" {
		response.BadRequest(c, "Location ID is required", "Missing location ID", "")
		return
	}

	location, err := h.locationService.GetLocationByID(c.Request.Context(), locationID)
	if err != nil {
		h.logger.Error("Error finding location by ID: %v", err)
		response.NotFound(c, "Location not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, location, "Location detail retrieved successfully")
}
