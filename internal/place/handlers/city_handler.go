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

// CityHandler handles HTTP requests related to cities
type CityHandler struct {
	cityService services.CityService
	logger      *logger.Logger
}

// NewCityHandler creates a new instance of city handler
func NewCityHandler(cityService services.CityService, logger *logger.Logger) *CityHandler {
	return &CityHandler{
		cityService: cityService,
		logger:      logger,
	}
}

// CreateCity godoc
// @Summary Create a new city
// @Description Add a new city to the system
// @Tags Cities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param city body models.City true "City Information"
// @Success 201 {object} response.APIResponse{data=models.City} "City created successfully"
// @Failure 400 {object} response.APIResponse "Invalid city creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /cities [post]
func (h *CityHandler) CreateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if city.Name == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name": city.Name == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.cityService.CreateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error creating city: %v", err)
		response.InternalServerError(c, "Failed to create city", err.Error(), "")
		return
	}

	response.SuccessCreated(c, city, "City created successfully")
}

// SearchCities godoc
// @Summary Search cities
// @Description Search cities by various criteria
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (name, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.City} "Cities found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /cities/search [get]
func (h *CityHandler) SearchCities(c *gin.Context) {
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

	// Search cities
	cities, err := h.cityService.SearchCities(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching cities: %v", err)
		response.InternalServerError(c, "Failed to search cities", err.Error(), "")
		return
	}

	// Count total search results
	totalCities, err := h.cityService.CountCities(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalCities,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalCities) / float64(limit))),
		HasNextPage: offset+limit < totalCities,
	}

	// Respond with cities and pagination
	response.SuccessOK(c, cities, "Cities found successfully", pagination)
}

// UpdateCity godoc
// @Summary Update a city
// @Description Update an existing city's details
// @Tags Cities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "City ID"
// @Param city body models.City true "City Update Details"
// @Success 200 {object} response.APIResponse{data=models.City} "City updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid city update details"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /cities/{id} [put]
func (h *CityHandler) UpdateCity(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		response.BadRequest(c, "City ID is required", "Missing city ID", "")
		return
	}

	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(cityID)
	if err != nil {
		response.BadRequest(c, "Invalid City ID", "Invalid UUID format", "")
		return
	}
	city.ID = parsedID

	if err := h.cityService.UpdateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error updating city: %v", err)
		response.InternalServerError(c, "Failed to update city", err.Error(), "")
		return
	}

	response.SuccessOK(c, city, "City updated successfully")
}

// DeleteCity godoc
// @Summary Delete a city
// @Description Remove a city from the system by its unique identifier
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "City ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid city ID"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /cities/{id} [delete]
func (h *CityHandler) DeleteCity(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		response.BadRequest(c, "City ID is required", "Missing city ID", "")
		return
	}

	if err := h.cityService.DeleteCity(c.Request.Context(), cityID); err != nil {
		h.logger.Error("Error deleting city: %v", err)
		response.InternalServerError(c, "Failed to delete city", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "City deleted successfully")
}

// ListCities godoc
// @Summary List cities
// @Description Retrieve a list of cities with pagination and filtering
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of cities to retrieve" default(10)
// @Param offset query int false "Number of cities to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Param province_id query string false "Filter by province ID"
// @Success 200 {object} response.APIResponse{data=[]models.City} "Cities retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list cities"
// @Router /cities [get]
func (h *CityHandler) ListCities(c *gin.Context) {
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
	if provinceID := c.Query("province_id"); provinceID != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "province_id",
			Operator: "eq",
			Value:    provinceID,
		})
	}
	listOptions.Filters = filters

	// Retrieve cities
	cities, err := h.cityService.ListCities(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error(), "")
		return
	}

	// Count total cities for pagination
	totalCities, err := h.cityService.CountCities(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalCities,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalCities) / float64(limit))),
		HasNextPage: offset+limit < totalCities,
	}

	// Respond with cities and pagination
	response.SuccessOK(c, cities, "Cities retrieved successfully", pagination)
}

// GetCityByID godoc
// @Summary Get city by ID
// @Description Retrieve a city's details by its unique identifier
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "City ID"
// @Success 200 {object} response.APIResponse{data=models.City} "City retrieved successfully"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /cities/{id} [get]
func (h *CityHandler) GetCityByID(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		response.BadRequest(c, "City ID is required", "Missing city ID", "")
		return
	}

	city, err := h.cityService.GetCityByID(c.Request.Context(), cityID)
	if err != nil {
		h.logger.Error("Error finding city by ID: %v", err)
		response.NotFound(c, "City not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, city, "City detail retrieved successfully")
}
