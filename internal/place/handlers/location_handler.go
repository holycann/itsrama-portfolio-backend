package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/response"
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
// @Success 201 {object} response.Response{data=models.Location} "Location created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid location creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.locationService.CreateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error creating location: %v", err)
		response.InternalServerError(c, "Failed to create location", err.Error())
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
// @Param id query string false "Location ID"
// @Param name query string false "Location Name"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Location} "Locations found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	// Get query parameters
	id := c.Query("id")
	name := c.Query("name")
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
		location, err := h.locationService.GetLocationByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding location by ID: %v", err)
			response.NotFound(c, "Location not found", err.Error())
			return
		}
		response.SuccessOK(c, location, "Location found")
		return
	}

	// If name is provided, search by name
	if name != "" {
		location, err := h.locationService.GetLocationByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding location by name: %v", err)
			response.NotFound(c, "Location not found", err.Error())
			return
		}
		response.SuccessOK(c, location, "Location found")
		return
	}

	// If no specific parameters are provided, return a list of locations
	locations, err := h.locationService.GetLocations(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving locations: %v", err)
		response.InternalServerError(c, "Failed to retrieve locations", err.Error())
		return
	}

	// Count total locations for pagination
	total, err := h.locationService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting locations: %v", err)
		response.InternalServerError(c, "Failed to count locations", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, locations, total, offset/limit+1, limit)
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
// @Success 200 {object} response.Response{data=models.Location} "Location updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid location update details"
// @Failure 404 {object} response.ErrorResponse "Location not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.locationService.UpdateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error updating location: %v", err)
		response.InternalServerError(c, "Failed to update location", err.Error())
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
// @Success 200 {object} response.Response "Location deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid location ID"
// @Failure 404 {object} response.ErrorResponse "Location not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Location ID is required", nil)
		return
	}

	if err := h.locationService.DeleteLocation(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting location: %v", err)
		response.InternalServerError(c, "Failed to delete location", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Location deleted successfully")
}

// ListLocation godoc
// @Summary List locations
// @Description Retrieve a list of locations with pagination
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of locations to retrieve" default(10)
// @Param offset query int false "Number of locations to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Location} "Locations retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list locations"
// @Router /locations [get]
func (h *LocationHandler) ListLocation(c *gin.Context) {
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

	// Get list of locations
	locations, err := h.locationService.GetLocations(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving locations: %v", err)
		response.InternalServerError(c, "Failed to retrieve locations", err.Error())
		return
	}

	// Count total locations for pagination
	total, err := h.locationService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting locations: %v", err)
		response.InternalServerError(c, "Failed to count locations", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, locations, total, offset/limit+1, limit)
}
