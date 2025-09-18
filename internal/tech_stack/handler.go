package tech_stack

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/internal/response"
	"github.com/holycann/itsrama-portfolio-backend/internal/utils"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
)

type TechStackHandler struct {
	base.BaseHandler
	techStackService TechStackService
}

func NewTechStackHandler(techStackService TechStackService, logger *logger.Logger) *TechStackHandler {
	return &TechStackHandler{
		BaseHandler:      *base.NewBaseHandler(logger),
		techStackService: techStackService,
	}
}

// CreateTechStack creates a new tech stack
// @Summary Create a new tech stack
// @Description Create a new tech stack with image upload
// @Tags Tech Stacks
// @Accept multipart/form-data
// @Produce json
// @Param image formData file false "Tech Stack Image"
// @Param payload formData string true "Tech Stack Details in JSON format (See TechStackCreate Model)"
// @Success 200 {object} response.APIResponse{data=TechStack} "Tech stack created successfully"
// @Failure 400 {object} response.APIResponse{data=TechStackCreate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /tech-stacks [post]
func (h *TechStackHandler) CreateTechStack(c *gin.Context) {
	var techStackInput TechStackCreate

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
	if err := utils.ExtractFormDataPayload(c, &techStackInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get image file
	imageFileHeaders, err := utils.ExtractFileHeaders(c, "image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(imageFileHeaders) > 0 {
		techStackInput.Image = imageFileHeaders[0]
	}

	// Create tech stack
	techStack, err := h.techStackService.CreateTechStack(c.Request.Context(), &techStackInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, techStack, "Tech stack created successfully")
}

// GetTechStackByID retrieves a specific tech stack
// @Summary Get a tech stack by ID
// @Description Retrieve a specific tech stack using its unique identifier
// @Tags Tech Stacks
// @Produce json
// @Param id path string true "Tech Stack ID"
// @Success 200 {object} response.APIResponse{data=TechStack} "Tech stack retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Tech stack not found"
// @Router /tech-stacks/{id} [get]
func (h *TechStackHandler) GetTechStackByID(c *gin.Context) {
	techStackID := c.Param("id")
	if techStackID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Tech stack ID is required",
			nil,
		))
		return
	}

	techStack, err := h.techStackService.GetTechStackByID(c.Request.Context(), techStackID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, techStack, "Tech stack retrieved successfully")
}

// UpdateTechStack updates an existing tech stack
// @Summary Update a tech stack
// @Description Update an existing tech stack with new details and optional image
// @Tags Tech Stacks
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Tech Stack ID"
// @Param image formData file false "Tech Stack Image"
// @Param payload formData string true "Tech Stack Update Details in JSON format (See TechStackUpdate Model)"
// @Success 200 {object} response.APIResponse{data=TechStack} "Tech stack updated successfully"
// @Failure 400 {object} response.APIResponse{data=TechStackUpdate} "Bad Request"
// @Failure 404 {object} response.APIResponse "Tech stack not found"
// @Router /tech-stacks/{id} [put]
func (h *TechStackHandler) UpdateTechStack(c *gin.Context) {
	var techStackInput TechStackUpdate

	// Parse ID from path
	techStackIDStr := c.Param("id")
	techStackID, err := uuid.Parse(techStackIDStr)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid tech stack ID",
			err,
		))
		return
	}

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
	if err := utils.ExtractFormDataPayload(c, &techStackInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Set the ID from path
	techStackInput.ID = techStackID

	// Get image file
	imageFileHeaders, err := utils.ExtractFileHeaders(c, "image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(imageFileHeaders) > 0 {
		techStackInput.Image = imageFileHeaders[0]
	}

	// Update tech stack
	updatedTechStack, err := h.techStackService.UpdateTechStack(c.Request.Context(), &techStackInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, updatedTechStack, "Tech stack updated successfully")
}

// DeleteTechStack deletes an existing tech stack
// @Summary Delete a tech stack
// @Description Delete a tech stack by its unique identifier
// @Tags Tech Stacks
// @Produce json
// @Param id path string true "Tech Stack ID"
// @Success 200 {object} response.APIResponse "Tech stack deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Tech stack not found"
// @Router /tech-stacks/{id} [delete]
func (h *TechStackHandler) DeleteTechStack(c *gin.Context) {
	techStackID := c.Param("id")
	if techStackID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Tech stack ID is required",
			nil,
		))
		return
	}

	err := h.techStackService.DeleteTechStack(c.Request.Context(), techStackID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Tech stack deleted successfully")
}

// ListTechStacks retrieves a paginated list of tech stacks
// @Summary List tech stacks
// @Description Retrieve a paginated list of tech stacks with optional filtering
// @Tags Tech Stacks
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Success 200 {object} response.APIResponse{data=[]TechStack} "Tech stacks retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /tech-stacks [get]
func (h *TechStackHandler) ListTechStacks(c *gin.Context) {
	// Parse pagination and filtering options
	opts, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid query parameters",
			err,
		))
		return
	}

	// Optional category filter
	category := c.Query("category")
	if category != "" {
		opts.Filters = append(opts.Filters, base.FilterOption{
			Field:    "category",
			Operator: base.OperatorEqual,
			Value:    category,
		})
	}

	// List tech stacks
	techStacks, err := h.techStackService.ListTechStacks(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Count total tech stacks for pagination
	total, err := h.techStackService.CountTechStacks(c.Request.Context(), opts.Filters)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, techStacks, "Tech stacks retrieved successfully",
		response.WithPagination(total, opts.Page, opts.PerPage))
}

// BulkCreateTechStacks creates multiple tech stacks in bulk
// @Summary Bulk create tech stacks
// @Description Create multiple tech stacks in a single request
// @Tags Tech Stacks
// @Accept multipart/form-data
// @Produce json
// @Param image formData []file false "Tech Stack Images"
// @Param payload formData string true "Tech Stack Details in JSON array format (See TechStackCreate Model)"
// @Success 200 {object} response.APIResponse{data=[]TechStack} "Tech stacks created successfully"
// @Failure 400 {object} response.APIResponse{data=[]TechStackCreate} "Bad Request"
// @Router /tech-stacks/bulk [post]
func (h *TechStackHandler) BulkCreateTechStacks(c *gin.Context) {
	var techStacksInput []*TechStackCreate

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
	if err := utils.ExtractFormDataPayload(c, &techStacksInput); err != nil {
		h.HandleError(c, err)
		return
	}

	imageFileHeaders, err := utils.ExtractFileHeaders(c, "image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Get image files
	if len(imageFileHeaders) > 0 {
		for i := range techStacksInput {
			if i < len(imageFileHeaders) {
				techStacksInput[i].Image = imageFileHeaders[i]
			}
		}
	}

	// Bulk create tech stacks
	techStacks, err := h.techStackService.BulkCreateTechStacks(c.Request.Context(), techStacksInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, techStacks, "Tech stacks created successfully")
}

// BulkUpdateTechStacks updates multiple tech stacks in bulk
// @Summary Bulk update tech stacks
// @Description Update multiple tech stacks in a single request
// @Tags Tech Stacks
// @Accept multipart/form-data
// @Produce json
// @Param images formData []file false "Tech Stack Images"
// @Param payload formData string true "Tech Stack Update Details in JSON array format (See TechStackUpdate Model)"
// @Success 200 {object} response.APIResponse{data=[]TechStack} "Tech stacks updated successfully"
// @Failure 400 {object} response.APIResponse{data=[]TechStackUpdate} "Bad Request"
// @Router /tech-stacks/bulk [put]
func (h *TechStackHandler) BulkUpdateTechStacks(c *gin.Context) {
	var techStacksInput []*TechStackUpdate

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
	if err := utils.ExtractFormDataPayload(c, &techStacksInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get image files
	imageFileHeaders, err := utils.ExtractFileHeaders(c, "image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(imageFileHeaders) > 0 {
		for i := range techStacksInput {
			techStacksInput[i].Image = imageFileHeaders[0]
		}
	}

	// Bulk update tech stacks
	updatedTechStacks, err := h.techStackService.BulkUpdateTechStacks(c.Request.Context(), techStacksInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, updatedTechStacks, "Tech stacks updated successfully")
}

// BulkDeleteTechStacks deletes multiple tech stacks in bulk
// @Summary Bulk delete tech stacks
// @Description Delete multiple tech stacks in a single request
// @Tags Tech Stacks
// @Accept json
// @Produce json
// @Param ids body []string true "Tech Stack IDs to delete"
// @Success 200 {object} response.APIResponse "Tech stacks deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /tech-stacks/bulk [delete]
func (h *TechStackHandler) BulkDeleteTechStacks(c *gin.Context) {
	var idsInput []string

	// Bind input
	if err := c.ShouldBindJSON(&idsInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
	}

	// Bulk delete tech stacks
	err := h.techStackService.BulkDeleteTechStacks(c.Request.Context(), idsInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Tech stacks deleted successfully")
}
