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

// ProvinceHandler menangani permintaan HTTP terkait provinsi
type ProvinceHandler struct {
	provinceService services.ProvinceService
	logger          *logger.Logger
}

// NewProvinceHandler membuat instance baru dari province handler
func NewProvinceHandler(provinceService services.ProvinceService, logger *logger.Logger) *ProvinceHandler {
	return &ProvinceHandler{
		provinceService: provinceService,
		logger:          logger,
	}
}

// CreateProvince godoc
// @Summary Create a new province
// @Description Add a new province to the system
// @Tags Provinces
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param province body models.Province true "Province Information"
// @Success 201 {object} response.APIResponse{data=models.Province} "Province created successfully"
// @Failure 400 {object} response.APIResponse "Invalid province creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /provinces [post]
func (h *ProvinceHandler) CreateProvince(c *gin.Context) {
	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		h.logger.Error("Error binding province: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if province.Name == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"name": province.Name == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.provinceService.CreateProvince(c.Request.Context(), &province); err != nil {
		h.logger.Error("Error creating province: %v", err)
		response.InternalServerError(c, "Failed to create province", err.Error(), "")
		return
	}

	response.SuccessCreated(c, province, "Province created successfully")
}

// GetProvinceByID godoc
// @Summary Get province by ID
// @Description Retrieve a province's details by its unique identifier
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Province ID"
// @Success 200 {object} response.APIResponse{data=models.Province} "Province retrieved successfully"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /provinces/{id} [get]
func (h *ProvinceHandler) GetProvinceByID(c *gin.Context) {
	// Get province ID from path parameter
	provinceID := c.Param("id")
	if provinceID == "" {
		response.BadRequest(c, "Province ID is required", "Missing province ID", "")
		return
	}

	province, err := h.provinceService.GetProvinceByID(c.Request.Context(), provinceID)
	if err != nil {
		h.logger.Error("Error finding province by ID: %v", err)
		response.NotFound(c, "Province not found", err.Error(), "")
		return
	}

	response.SuccessOK(c, province, "Province found")
}

// SearchProvinces godoc
// @Summary Search provinces
// @Description Search provinces by various criteria
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (name, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.Province} "Provinces found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /provinces/search [get]
func (h *ProvinceHandler) SearchProvinces(c *gin.Context) {
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

	// Search provinces
	provinces, err := h.provinceService.SearchProvinces(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching provinces: %v", err)
		response.InternalServerError(c, "Failed to search provinces", err.Error(), "")
		return
	}

	// Count total search results
	totalProvinces, err := h.provinceService.CountProvinces(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalProvinces,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalProvinces) / float64(limit))),
		HasNextPage: offset+limit < totalProvinces,
	}

	// Respond with provinces and pagination
	response.SuccessOK(c, provinces, "Provinces found successfully", pagination)
}

// UpdateProvince godoc
// @Summary Update a province
// @Description Update an existing province's details
// @Tags Provinces
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Province ID"
// @Param province body models.Province true "Province Update Details"
// @Success 200 {object} response.APIResponse{data=models.Province} "Province updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid province update details"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /provinces/{id} [put]
func (h *ProvinceHandler) UpdateProvince(c *gin.Context) {
	// Get province ID from path parameter
	provinceID := c.Param("id")
	if provinceID == "" {
		response.BadRequest(c, "Province ID is required", "Missing province ID", "")
		return
	}

	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		h.logger.Error("Error binding province: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(provinceID)
	if err != nil {
		response.BadRequest(c, "Invalid Province ID", "Invalid UUID format", "")
		return
	}
	province.ID = parsedID

	if err := h.provinceService.UpdateProvince(c.Request.Context(), &province); err != nil {
		h.logger.Error("Error updating province: %v", err)
		response.InternalServerError(c, "Failed to update province", err.Error(), "")
		return
	}

	response.SuccessOK(c, province, "Province updated successfully")
}

// DeleteProvince godoc
// @Summary Delete a province
// @Description Remove a province from the system by its unique identifier
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Province ID"
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid province ID"
// @Failure 404 {object} response.APIResponse "Province not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /provinces/{id} [delete]
func (h *ProvinceHandler) DeleteProvince(c *gin.Context) {
	// Get province ID from path parameter
	provinceID := c.Param("id")
	if provinceID == "" {
		response.BadRequest(c, "Province ID is required", "Missing province ID", "")
		return
	}

	if err := h.provinceService.DeleteProvince(c.Request.Context(), provinceID); err != nil {
		h.logger.Error("Error deleting province: %v", err)
		response.InternalServerError(c, "Failed to delete province", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Province deleted successfully")
}

// ListProvinces godoc
// @Summary List provinces
// @Description Retrieve a list of provinces with pagination and filtering
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of provinces to retrieve" default(10)
// @Param offset query int false "Number of provinces to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.Province} "Provinces retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list provinces"
// @Router /provinces [get]
func (h *ProvinceHandler) ListProvinces(c *gin.Context) {
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
	if name := c.Query("name"); name != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "name",
			Operator: "like",
			Value:    name,
		})
	}
	listOptions.Filters = filters

	// Retrieve provinces
	provinces, err := h.provinceService.ListProvinces(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving provinces: %v", err)
		response.InternalServerError(c, "Failed to retrieve provinces", err.Error(), "")
		return
	}

	// Count total provinces for pagination
	totalProvinces, err := h.provinceService.CountProvinces(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting provinces: %v", err)
		response.InternalServerError(c, "Failed to count provinces", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalProvinces,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalProvinces) / float64(limit))),
		HasNextPage: offset+limit < totalProvinces,
	}

	// Respond with provinces and pagination
	response.SuccessOK(c, provinces, "Provinces retrieved successfully", pagination)
}
