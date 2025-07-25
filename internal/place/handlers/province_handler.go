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
// @Success 201 {object} response.Response{data=models.Province} "Province created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid province creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /provinces [post]
func (h *ProvinceHandler) CreateProvince(c *gin.Context) {
	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		h.logger.Error("Error binding province: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.provinceService.CreateProvince(c.Request.Context(), &province); err != nil {
		h.logger.Error("Error creating province: %v", err)
		response.InternalServerError(c, "Failed to create province", err.Error())
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
// @Success 200 {object} response.Response{data=models.Province} "Province retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Province not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /provinces/{id} [get]
func (h *ProvinceHandler) GetProvinceByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Province ID is required", nil)
		return
	}

	province, err := h.provinceService.GetProvinceByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Error finding province by ID: %v", err)
		response.NotFound(c, "Province not found", err.Error())
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
// @Param id query string false "Province ID"
// @Param name query string false "Province Name"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Province} "Provinces found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /provinces/search [get]
func (h *ProvinceHandler) SearchProvinces(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
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

	if id != "" {
		province, err := h.provinceService.GetProvinceByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding province by ID: %v", err)
			response.NotFound(c, "Province not found", err.Error())
			return
		}
		response.SuccessOK(c, province, "Province found")
		return
	}

	if name != "" {
		province, err := h.provinceService.GetProvinceByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding province by name: %v", err)
			response.NotFound(c, "Province not found", err.Error())
			return
		}
		response.SuccessOK(c, province, "Province found")
		return
	}

	provinces, err := h.provinceService.GetProvinces(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving provinces: %v", err)
		response.InternalServerError(c, "Failed to retrieve provinces", err.Error())
		return
	}

	total, err := h.provinceService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting provinces: %v", err)
		response.InternalServerError(c, "Failed to count provinces", err.Error())
		return
	}

	response.WithPagination(c, provinces, total, offset/limit+1, limit)
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
// @Success 200 {object} response.Response{data=models.Province} "Province updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid province update details"
// @Failure 404 {object} response.ErrorResponse "Province not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /provinces/{id} [put]
func (h *ProvinceHandler) UpdateProvince(c *gin.Context) {
	var province models.Province
	if err := c.ShouldBindJSON(&province); err != nil {
		h.logger.Error("Error binding province: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.provinceService.UpdateProvince(c.Request.Context(), &province); err != nil {
		h.logger.Error("Error updating province: %v", err)
		response.InternalServerError(c, "Failed to update province", err.Error())
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
// @Success 200 {object} response.Response "Province deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid province ID"
// @Failure 404 {object} response.ErrorResponse "Province not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /provinces/{id} [delete]
func (h *ProvinceHandler) DeleteProvince(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Province ID is required", nil)
		return
	}

	if err := h.provinceService.DeleteProvince(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting province: %v", err)
		response.InternalServerError(c, "Failed to delete province", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Province deleted successfully")
}

// ListProvinces godoc
// @Summary List provinces
// @Description Retrieve a list of provinces with pagination
// @Tags Provinces
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of provinces to retrieve" default(10)
// @Param offset query int false "Number of provinces to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.Province} "Provinces retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list provinces"
// @Router /provinces [get]
func (h *ProvinceHandler) ListProvinces(c *gin.Context) {
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

	provinces, err := h.provinceService.GetProvinces(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving provinces: %v", err)
		response.InternalServerError(c, "Failed to retrieve provinces", err.Error())
		return
	}

	total, err := h.provinceService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting provinces: %v", err)
		response.InternalServerError(c, "Failed to count provinces", err.Error())
		return
	}

	response.WithPagination(c, provinces, total, offset/limit+1, limit)
}
