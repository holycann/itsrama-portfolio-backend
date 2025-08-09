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
	"github.com/holycann/cultour-backend/pkg/response"
)

// CityHandler handles HTTP requests related to cities
type CityHandler struct {
	*base.BaseHandler
	cityService services.CityService
}

// NewCityHandler creates a new instance of city handler
func NewCityHandler(cityService services.CityService, logger *logger.Logger) *CityHandler {
	return &CityHandler{
		BaseHandler: base.NewBaseHandler(logger),
		cityService: cityService,
	}
}

// CreateCity godoc
// @Summary Create a new city
// @Description Allows administrators to add a new city to the system
// @Description Supports creating cities with detailed information and optional image
// @Tags Cities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param city body models.CityCreate true "City Creation Details"
// @Success 201 {object} response.APIResponse{data=models.CityDTO} "City successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid city creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges (admin role required)"
// @Failure 500 {object} response.APIResponse "Internal server error during city creation"
// @Router /cities [post]
func (h *CityHandler) CreateCity(c *gin.Context) {
	var cityCreate models.CityCreate
	if err := h.ValidateRequest(c, &cityCreate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Convert CityCreate to City
	city := &models.City{
		Name:        cityCreate.Name,
		Description: cityCreate.Description,
		ProvinceID:  cityCreate.ProvinceID,
		ImageURL:    cityCreate.ImageURL,
	}

	createdCity, err := h.cityService.CreateCity(c.Request.Context(), city)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to create city"))
		return
	}

	h.HandleSuccess(c, createdCity, "City created successfully")
}

// SearchCities godoc
// @Summary Search cities
// @Description Performs a full-text search across city details with advanced filtering
// @Description Allows finding cities by keywords, province, and other attributes
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding cities" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,name,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param province_id query string false "Filter cities by specific province"
// @Success 200 {object} response.APIResponse{data=[]models.CityDTO} "Successfully completed city search"
// @Success 204 {object} response.APIResponse "No cities match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during city search"
// @Router /cities/search [get]
func (h *CityHandler) SearchCities(c *gin.Context) {
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

	// Add optional search filter
	if query := c.Query("query"); query != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "name",
			Operator: base.OperatorLike,
			Value:    query,
		})
	}

	// Perform search
	cities, total, err := h.cityService.SearchCities(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to search cities"))
		return
	}

	// Respond with pagination based on total and requested options
	h.HandleSuccess(c, cities, "Cities retrieved successfully",
		response.WithPagination(total, listOptions.Page, listOptions.PerPage))
}

// UpdateCity godoc
// @Summary Update an existing city
// @Description Allows administrators to modify city details
// @Description Supports partial updates with optional fields
// @Tags Cities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique City Identifier" format(uuid)
// @Param city body models.CityUpdate true "City Update Payload"
// @Success 200 {object} response.APIResponse{data=models.CityDTO} "City successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid city update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error during city update"
// @Router /cities/{id} [put]
func (h *CityHandler) UpdateCity(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		h.HandleError(c, errors.New(errors.ErrValidation, "City ID is required", nil))
		return
	}

	var cityUpdate models.CityUpdate
	if err := h.ValidateRequest(c, &cityUpdate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Parse the ID
	parsedID, err := uuid.Parse(cityID)
	if err != nil {
		h.HandleError(c, errors.New(errors.ErrValidation, "Invalid City ID", err))
		return
	}

	// Convert CityUpdate to City
	city := &models.City{
		ID:          parsedID,
		Name:        cityUpdate.Name,
		Description: cityUpdate.Description,
		ProvinceID:  cityUpdate.ProvinceID,
		ImageURL:    cityUpdate.ImageURL,
	}

	updatedCity, err := h.cityService.UpdateCity(c.Request.Context(), city)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update city"))
		return
	}

	h.HandleSuccess(c, updatedCity, "City updated successfully")
}

// DeleteCity godoc
// @Summary Delete a city
// @Description Allows administrators to permanently remove a city from the system
// @Description Deletes the city and its associated resources
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique City Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "City successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid city ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error during city deletion"
// @Router /cities/{id} [delete]
func (h *CityHandler) DeleteCity(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		h.HandleError(c, errors.New(errors.ErrValidation, "City ID is required", nil))
		return
	}

	if err := h.cityService.DeleteCity(c.Request.Context(), cityID); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to delete city"))
		return
	}

	h.HandleSuccess(c, nil, "City deleted successfully")
}

// ListCities godoc
// @Summary Retrieve cities list
// @Description Fetches a paginated list of cities with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of cities per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort cities by" default("created_at)" Enum(created_at,name)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param province_id query string false "Filter cities by specific province"
// @Success 200 {object} response.APIResponse{data=[]models.CityDTO} "Successfully retrieved cities list"
// @Success 204 {object} response.APIResponse "No cities found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during cities retrieval"
// @Router /cities [get]
func (h *CityHandler) ListCities(c *gin.Context) {
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

	// Optional filtering by province
	if provinceID := c.Query("province_id"); provinceID != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "province_id",
			Operator: base.OperatorEqual,
			Value:    provinceID,
		})
	}

	// Retrieve cities
	cities, total, err := h.cityService.ListCities(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to retrieve cities"))
		return
	}

	// Respond with pagination based on total and requested options
	h.HandleSuccess(c, cities, "Cities retrieved successfully",
		response.WithPagination(total, listOptions.Page, listOptions.PerPage))
}

// GetCityByID godoc
// @Summary Retrieve a specific city
// @Description Fetches comprehensive details of a city by its unique identifier
// @Description Returns full city information including province details
// @Tags Cities
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique City Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.CityDTO} "Successfully retrieved city details"
// @Failure 400 {object} response.APIResponse "Invalid city ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "City not found"
// @Failure 500 {object} response.APIResponse "Internal server error during city retrieval"
// @Router /cities/{id} [get]
func (h *CityHandler) GetCityByID(c *gin.Context) {
	// Get city ID from path parameter
	cityID := c.Param("id")
	if cityID == "" {
		h.HandleError(c, errors.New(errors.ErrValidation, "City ID is required", nil))
		return
	}

	city, err := h.cityService.GetCityByID(c.Request.Context(), cityID)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "City not found"))
		return
	}

	h.HandleSuccess(c, city, "City detail retrieved successfully")
}
