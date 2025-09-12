package tech_stack

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/pkg/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
	"github.com/holycann/itsrama-portfolio-backend/pkg/response"
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
func (h *TechStackHandler) CreateTechStack(c *gin.Context) {
	var techStackInput TechStackCreate

	// Bind input
	if err := c.ShouldBindJSON(&techStackInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
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

	// Bind input
	if err := c.ShouldBindJSON(&techStackInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
	}

	// Set the ID from path
	techStackInput.ID = techStackID

	// Update tech stack
	updatedTechStack, err := h.techStackService.UpdateTechStack(c.Request.Context(), &techStackInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, updatedTechStack, "Tech stack updated successfully")
}

// DeleteTechStack deletes an existing tech stack
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

// GetTechStacksByCategory retrieves tech stacks by category
func (h *TechStackHandler) GetTechStacksByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Category is required",
			nil,
		))
		return
	}

	techStacks, err := h.techStackService.GetTechStacksByCategory(c.Request.Context(), category)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, techStacks, "Tech stacks retrieved successfully")
}

// BulkCreateTechStacks creates multiple tech stacks in bulk
func (h *TechStackHandler) BulkCreateTechStacks(c *gin.Context) {
	var techStacksInput []*TechStackCreate

	// Bind input
	if err := c.ShouldBindJSON(&techStacksInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
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
func (h *TechStackHandler) BulkUpdateTechStacks(c *gin.Context) {
	var techStacksInput []*TechStackUpdate

	// Bind input
	if err := c.ShouldBindJSON(&techStacksInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
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
