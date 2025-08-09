package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"

	_ "github.com/holycann/cultour-backend/pkg/response"
)

// LocationHandler handles HTTP requests related to locations
type LocationHandler struct {
	base.BaseHandler
	locationService services.LocationService
}

// NewLocationHandler creates a new instance of location handler
func NewLocationHandler(
	locationService services.LocationService,
	logger *logger.Logger,
) *LocationHandler {
	return &LocationHandler{
		BaseHandler:     *base.NewBaseHandler(logger),
		locationService: locationService,
	}
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Allows administrators to add a new geographical location to the system
// @Description Supports creating locations with detailed geospatial information
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param location body models.LocationCreate true "Location Creation Details"
// @Success 201 {object} response.APIResponse{data=models.LocationDTO} "Location successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid location creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges (admin role required)"
// @Failure 500 {object} response.APIResponse "Internal server error during location creation"
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var locationCreate models.LocationCreate

	// Validate request
	if err := h.ValidateRequest(c, &locationCreate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Create location
	createdLocation, err := h.locationService.CreateLocation(c.Request.Context(), &locationCreate)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with created location
	h.HandleCreated(c, createdLocation.ToDTO(), "Location created successfully")
}

// SearchLocations godoc
// @Summary Search locations
// @Description Performs a full-text search across location details with advanced filtering
// @Description Allows finding locations by keywords, city, coordinates, and other attributes
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding locations" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,name,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param city_id query string false "Filter locations by specific city"
// @Param latitude query float64 false "Latitude for proximity search"
// @Param longitude query float64 false "Longitude for proximity search"
// @Param radius query float64 false "Search radius in kilometers for proximity search" minimum(0)
// @Success 200 {object} response.APIResponse{data=[]models.LocationDTO} "Successfully completed location search"
// @Success 204 {object} response.APIResponse "No locations match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during location search"
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	// Manually set list options with default values
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		perPage = 10
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	listOptions := base.ListOptions{
		Page:      page,
		PerPage:   perPage,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters:   []base.FilterOption{},
	}

	// Prepare search filter
	searchFilter := models.LocationDTO{
		Name: c.Query("query"),
	}
	if cityID := c.Query("city_id"); cityID != "" {
		cityUUID, err := uuid.Parse(cityID)
		if err != nil {
			h.HandleError(c, errors.New(
				errors.ErrValidation,
				"Invalid city ID",
				err,
			))
			return
		}
		searchFilter.CityID = cityUUID
	}

	// Build filters dynamically
	filters := base.BuildFilterFromStruct(searchFilter)

	// Update list options with filters
	listOptions.Filters = filters

	// Search locations
	locations, err := h.locationService.SearchLocations(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Handle pagination (todo: return accurate total from service)
	h.HandlePagination(c, locations, len(locations), listOptions)
}

// UpdateLocation godoc
// @Summary Update an existing location
// @Description Allows administrators to modify location details
// @Description Supports partial updates with optional fields
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Location Identifier" format(uuid)
// @Param location body models.LocationUpdate true "Location Update Payload"
// @Success 200 {object} response.APIResponse{data=models.LocationDTO} "Location successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid location update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error during location update"
// @Router /locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	// Validate and parse location ID
	locationID, err := h.ValidateUUID(c.Param("id"), "Location ID")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Validate request
	var locationUpdate models.LocationUpdate
	if err := h.ValidateRequest(c, &locationUpdate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Set ID for update
	locationUpdate.ID = locationID

	// Update location
	updatedLocation, err := h.locationService.UpdateLocation(c.Request.Context(), &locationUpdate)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with updated location
	h.HandleSuccess(c, updatedLocation.ToDTO(), "Location updated successfully")
}

// DeleteLocation godoc
// @Summary Delete a location
// @Description Allows administrators to permanently remove a location from the system
// @Description Deletes the location and its associated resources
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Location Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Location successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid location ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error during location deletion"
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	// Validate and parse location ID
	locationID := c.Param("id")
	if locationID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Location ID cannot be empty",
			nil,
		))
		return
	}

	// Delete location
	if err := h.locationService.DeleteLocation(c.Request.Context(), locationID); err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with success
	h.HandleSuccess(c, nil, "Location deleted successfully")
}

// ListLocations godoc
// @Summary Retrieve locations list
// @Description Fetches a paginated list of locations with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of locations per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort locations by" default("created_at)" Enum(created_at,name)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param city_id query string false "Filter locations by specific city"
// @Success 200 {object} response.APIResponse{data=[]models.LocationDTO} "Successfully retrieved locations list"
// @Success 204 {object} response.APIResponse "No locations found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during locations retrieval"
// @Router /locations [get]
func (h *LocationHandler) ListLocations(c *gin.Context) {
	// Manually set list options with default values
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		perPage = 10
	}

	listOptions := base.ListOptions{
		Page:    page,
		PerPage: perPage,
	}

	// Prepare filter
	var filters []base.FilterOption
	if cityID := c.Query("city_id"); cityID != "" {
		cityUUID, err := uuid.Parse(cityID)
		if err != nil {
			h.HandleError(c, errors.New(
				errors.ErrValidation,
				"Invalid city ID",
				err,
			))
			return
		}
		filters = append(filters, base.FilterOption{
			Field:    "city_id",
			Operator: base.OperatorEqual,
			Value:    cityUUID,
		})
	}
	listOptions.Filters = filters

	// Retrieve locations
	locations, err := h.locationService.ListLocations(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Handle pagination (todo: return accurate total from service)
	h.HandlePagination(c, locations, len(locations), listOptions)
}

// GetLocationByID godoc
// @Summary Retrieve a specific location
// @Description Fetches comprehensive details of a location by its unique identifier
// @Description Returns full location information including city details and geospatial data
// @Tags Locations
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Location Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.LocationDTO} "Successfully retrieved location details"
// @Failure 400 {object} response.APIResponse "Invalid location ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "Location not found"
// @Failure 500 {object} response.APIResponse "Internal server error during location retrieval"
// @Router /locations/{id} [get]
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	// Validate and parse location ID
	locationID := c.Param("id")
	if locationID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Location ID cannot be empty",
			nil,
		))
		return
	}

	// Retrieve location
	location, err := h.locationService.GetLocationByID(c.Request.Context(), locationID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Respond with location details
	h.HandleSuccess(c, location, "Location detail retrieved successfully")
}
