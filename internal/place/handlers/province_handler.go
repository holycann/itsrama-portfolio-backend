package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/holycann/cultour-backend/pkg/response"
)

// ProvinceHandler handles HTTP requests related to provinces
type ProvinceHandler struct {
	base.BaseHandler
	provinceService services.ProvinceService
}

// NewProvinceHandler creates a new instance of province handler
func NewProvinceHandler(
	provinceService services.ProvinceService,
	appLogger *logger.Logger,
) *ProvinceHandler {
	return &ProvinceHandler{
		BaseHandler:     *base.NewBaseHandler(appLogger),
		provinceService: provinceService,
	}
}

// CreateProvince godoc
// @Summary Create a new province
// @Description Allows administrators to add a new administrative province to the system
// @Description Supports creating provinces with detailed geographical information
// @Tags Provinces
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param province body models.ProvinceCreate true "Province Creation Details"
// @Success 201 {object} response.APIResponse{data=models.ProvinceDTO} "Province successfully created with full details"
// @Failure 400 {object} response.APIResponse "Invalid province creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges (admin role required)"
// @Failure 500 {object} response.APIResponse "Internal server error during province creation"
// @Router /provinces [post]
func (h *ProvinceHandler) CreateProvince(c *gin.Context) {
	var provinceCreate models.ProvinceCreate
	if err := c.ShouldBindJSON(&provinceCreate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid province creation details"))
		return
	}

	// Validate the input
	if err := base.ValidateModel(provinceCreate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Convert ProvinceCreate to Province
	province := &models.Province{
		Name:        provinceCreate.Name,
		Description: provinceCreate.Description,
	}

	createdProvince, err := h.provinceService.CreateProvince(c.Request.Context(), province)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to create province"))
		return
	}

	h.HandleSuccess(c, createdProvince, "Province created successfully")
}

// GetProvinceByID godoc
// @Summary Retrieve a specific province
// @Description Fetches comprehensive details of a province by its unique identifier
// @Description Returns full province information including associated cities
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Province Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=models.ProvinceDTO} "Successfully retrieved province details"
// @Failure 400 {object} response.APIResponse "Invalid province ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error during province retrieval"
// @Router /provinces/{id} [get]
func (h *ProvinceHandler) GetProvinceByID(c *gin.Context) {
	provinceID := c.Param("id")

	// Validate UUID
	parsedID, err := h.ValidateUUID(provinceID, "province_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	province, err := h.provinceService.GetProvinceByID(c.Request.Context(), parsedID.String())
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrNotFound, "Province not found"))
		return
	}

	h.HandleSuccess(c, province, "Province retrieved successfully")
}

// SearchProvinces godoc
// @Summary Search provinces
// @Description Performs a full-text search across province details with advanced filtering
// @Description Allows finding provinces by keywords and other attributes
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding provinces" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,name,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.ProvinceDTO} "Successfully completed province search"
// @Success 204 {object} response.APIResponse "No provinces match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during province search"
// @Router /provinces/search [get]
func (h *ProvinceHandler) SearchProvinces(c *gin.Context) {
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
	provinces, _, err := h.provinceService.SearchProvinces(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to search provinces"))
		return
	}

	// Create pagination
	data, pagination := base.PaginateResults(provinces, listOptions.Page, listOptions.PerPage)

	h.HandleSuccess(c, data, "Provinces retrieved successfully",
		response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}

// UpdateProvince godoc
// @Summary Update an existing province
// @Description Allows administrators to modify province details
// @Description Supports partial updates with optional fields
// @Tags Provinces
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Province Identifier" format(uuid)
// @Param province body models.ProvinceUpdate true "Province Update Payload"
// @Success 200 {object} response.APIResponse{data=models.ProvinceDTO} "Province successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid province update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error during province update"
// @Router /provinces/{id} [put]
func (h *ProvinceHandler) UpdateProvince(c *gin.Context) {
	// Validate and parse province ID
	provinceID := c.Param("id")
	parsedID, err := h.ValidateUUID(provinceID, "province_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Bind and validate update payload
	var provinceUpdate models.ProvinceUpdate
	if err := c.ShouldBindJSON(&provinceUpdate); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrValidation, "Invalid province update details"))
		return
	}

	// Validate the input
	if err := base.ValidateModel(provinceUpdate); err != nil {
		h.HandleError(c, err)
		return
	}

	// Prepare province for update
	province := &models.Province{
		ID:          parsedID,
		Name:        provinceUpdate.Name,
		Description: provinceUpdate.Description,
	}

	// Perform update
	updatedProvince, err := h.provinceService.UpdateProvince(c.Request.Context(), province)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to update province"))
		return
	}

	h.HandleSuccess(c, updatedProvince, "Province updated successfully")
}

// DeleteProvince godoc
// @Summary Delete a province
// @Description Allows administrators to permanently remove a province from the system
// @Description Deletes the province and its associated resources
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Admin JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Province Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Province successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid province ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient privileges"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error during province deletion"
// @Router /provinces/{id} [delete]
func (h *ProvinceHandler) DeleteProvince(c *gin.Context) {
	// Validate and parse province ID
	provinceID := c.Param("id")
	parsedID, err := h.ValidateUUID(provinceID, "province_id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Perform deletion
	if err := h.provinceService.DeleteProvince(c.Request.Context(), parsedID.String()); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to delete province"))
		return
	}

	h.HandleSuccess(c, nil, "Province deleted successfully")
}

// ListProvinces godoc
// @Summary Retrieve provinces list
// @Description Fetches a paginated list of provinces with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of provinces per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort provinces by" default("created_at)" Enum(created_at,name)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.ProvinceDTO} "Successfully retrieved provinces list"
// @Success 204 {object} response.APIResponse "No provinces found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during provinces retrieval"
// @Router /provinces [get]
func (h *ProvinceHandler) ListProvinces(c *gin.Context) {
	// Parse list options manually
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

	// Optional name filter
	if name := c.Query("name"); name != "" {
		listOptions.Filters = append(listOptions.Filters, base.FilterOption{
			Field:    "name",
			Operator: base.OperatorLike,
			Value:    name,
		})
	}

	// Retrieve provinces
	provinces, _, err := h.provinceService.SearchProvinces(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrDatabase, "Failed to retrieve provinces"))
		return
	}

	// Create pagination
	data, pagination := base.PaginateResults(provinces, listOptions.Page, listOptions.PerPage)

	h.HandleSuccess(c, data, "Provinces retrieved successfully",
		response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}
