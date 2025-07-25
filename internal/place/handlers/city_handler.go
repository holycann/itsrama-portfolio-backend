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
// @Success 201 {object} response.Response{data=models.City} "City created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid city creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /cities [post]
func (h *CityHandler) CreateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.cityService.CreateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error creating city: %v", err)
		response.InternalServerError(c, "Failed to create city", err.Error())
		return
	}

	response.SuccessCreated(c, city, "City created successfully")
}

// SearchCity godoc
// @Summary Search cities
// @Description Search cities by various criteria
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id query string false "City ID"
// @Param name query string false "City Name"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.City} "Cities found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /cities/search [get]
func (h *CityHandler) SearchCity(c *gin.Context) {
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
		city, err := h.cityService.GetCityByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding city by ID: %v", err)
			response.NotFound(c, "City not found", err.Error())
			return
		}
		response.SuccessOK(c, city, "City found")
		return
	}

	// If name is provided, search by name
	if name != "" {
		city, err := h.cityService.GetCityByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding city by name: %v", err)
			response.NotFound(c, "City not found", err.Error())
			return
		}
		response.SuccessOK(c, city, "City found")
		return
	}

	// If no specific parameters are provided, return a list of cities
	cities, err := h.cityService.GetCities(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Count total cities for pagination
	total, err := h.cityService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, cities, total, offset/limit+1, limit)
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
// @Success 200 {object} response.Response{data=models.City} "City updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid city update details"
// @Failure 404 {object} response.ErrorResponse "City not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /cities/{id} [put]
func (h *CityHandler) UpdateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.cityService.UpdateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error updating city: %v", err)
		response.InternalServerError(c, "Failed to update city", err.Error())
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
// @Success 200 {object} response.Response "City deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid city ID"
// @Failure 404 {object} response.ErrorResponse "City not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /cities/{id} [delete]
func (h *CityHandler) DeleteCity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "City ID is required", nil)
		return
	}

	if err := h.cityService.DeleteCity(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting city: %v", err)
		response.InternalServerError(c, "Failed to delete city", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "City deleted successfully")
}

// ListCities godoc
// @Summary List cities
// @Description Retrieve a list of cities with pagination
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of cities to retrieve" default(10)
// @Param offset query int false "Number of cities to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.City} "Cities retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list cities"
// @Router /cities [get]
func (h *CityHandler) ListCities(c *gin.Context) {
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

	// Get list of cities
	cities, err := h.cityService.GetCities(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Count total cities for pagination
	total, err := h.cityService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, cities, total, offset/limit+1, limit)
}

// GetCityByID godoc
// @Summary Get city by ID
// @Description Retrieve a city's details by its unique identifier
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "City ID"
// @Success 200 {object} response.Response{data=models.City} "City retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "City not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /cities/{id} [get]
func (h *CityHandler) GetCityByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "City ID is required", nil)
		return
	}
	city, err := h.cityService.GetCityByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Error finding city by ID: %v", err)
		response.NotFound(c, "City not found", err.Error())
		return
	}
	response.SuccessOK(c, city, "City detail retrieved successfully")
}
