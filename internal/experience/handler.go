package experience

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/pkg/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
	"github.com/holycann/itsrama-portfolio-backend/pkg/response"
)

type ExperienceHandler struct {
	base.BaseHandler
	experienceService ExperienceService
}

func NewExperienceHandler(experienceService ExperienceService, logger *logger.Logger) *ExperienceHandler {
	return &ExperienceHandler{
		BaseHandler:       *base.NewBaseHandler(logger),
		experienceService: experienceService,
	}
}

// CreateExperience creates a new experience
func (h *ExperienceHandler) CreateExperience(c *gin.Context) {
	var experienceInput ExperienceCreate

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
	if err := h.validateAndExtractExperienceInput(c, &experienceInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get logo file
	logoFileHeaders, err := h.extractFileHeaders(c, "logo_image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(logoFileHeaders) > 0 {
		experienceInput.LogoImage = logoFileHeaders[0]
	}

	// Get image files
	imageFileHeaders, err := h.extractFileHeaders(c, "images", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(imageFileHeaders) > 0 {
		experienceInput.Images = imageFileHeaders
	}

	// Create experience via service
	experience, err := h.experienceService.CreateExperience(c.Request.Context(), &experienceInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experience created successfully")
}

// GetExperienceByID retrieves a specific experience
func (h *ExperienceHandler) GetExperienceByID(c *gin.Context) {
	experienceID := c.Param("id")
	if experienceID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Experience ID is required",
			nil,
		))
		return
	}

	experience, err := h.experienceService.GetExperienceByID(c.Request.Context(), experienceID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experience retrieved successfully")
}

// UpdateExperience updates an existing experience
func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	var experienceInput ExperienceUpdate

	// Extract experience ID from path parameter
	experienceIDStr := c.Param("id")
	experienceID, err := uuid.Parse(experienceIDStr)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid experience ID",
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

	// Set the ID in the payload
	experienceInput.ID = experienceID

	// Extract and validate form fields
	if err := h.validateAndExtractExperienceInput(c, &experienceInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get logo file
	logoFileHeaders, err := h.extractFileHeaders(c, "logo_image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	experienceInput.LogoImage = logoFileHeaders[0]

	// Get image files
	imageFileHeaders, err := h.extractFileHeaders(c, "images", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	experienceInput.Images = imageFileHeaders

	// Update experience via service
	experience, err := h.experienceService.UpdateExperience(c.Request.Context(), &experienceInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experience updated successfully")
}

// DeleteExperience deletes an existing experience
func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	experienceID := c.Param("id")
	if experienceID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Experience ID is required",
			nil,
		))
		return
	}

	err := h.experienceService.DeleteExperience(c.Request.Context(), experienceID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Experience deleted successfully")
}

// ListExperiences retrieves a paginated list of experience
func (h *ExperienceHandler) ListExperiences(c *gin.Context) {
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

	// Optional company filter
	company := c.Query("company")
	if company != "" {
		opts.Filters = append(opts.Filters, base.FilterOption{
			Field:    "company",
			Operator: base.OperatorEqual,
			Value:    company,
		})
	}

	// List experience
	experience, err := h.experienceService.ListExperiences(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Count total experience for pagination
	total, err := h.experienceService.CountExperiences(c.Request.Context(), opts.Filters)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experiences retrieved successfully",
		response.WithPagination(total, opts.Page, opts.PerPage))
}

// GetExperiencesByCompany retrieves experience by company
func (h *ExperienceHandler) GetExperiencesByCompany(c *gin.Context) {
	company := c.Param("company")
	if company == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Company is required",
			nil,
		))
		return
	}

	experience, err := h.experienceService.GetExperiencesByCompany(c.Request.Context(), company)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experiences retrieved successfully")
}

// SearchExperiences performs a full-text search on experience
func (h *ExperienceHandler) SearchExperiences(c *gin.Context) {
	query := c.Query("query")
	opts, err := base.ParsePaginationParams(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid query parameters",
			err,
		))
		return
	}

	// Attach search term
	opts.Search = query

	experience, total, err := h.experienceService.SearchExperiences(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experiences search completed successfully",
		response.WithPagination(total, opts.Page, opts.PerPage))
}

// BulkCreateExperiences creates multiple experience in bulk
func (h *ExperienceHandler) BulkCreateExperiences(c *gin.Context) {
	var experiencesInput []*ExperienceCreate

	// Bind input
	if err := c.ShouldBindJSON(&experiencesInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
	}

	// Bulk create experience
	experience, err := h.experienceService.BulkCreateExperiences(c.Request.Context(), experiencesInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experiences created successfully")
}

// BulkUpdateExperiences updates multiple experience in bulk
func (h *ExperienceHandler) BulkUpdateExperiences(c *gin.Context) {
	var experiencesInput []*ExperienceUpdate

	// Bind input
	if err := c.ShouldBindJSON(&experiencesInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
	}

	// Bulk update experience
	updatedExperiences, err := h.experienceService.BulkUpdateExperiences(c.Request.Context(), experiencesInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, updatedExperiences, "Experiences updated successfully")
}

// BulkDeleteExperiences deletes multiple experience in bulk
func (h *ExperienceHandler) BulkDeleteExperiences(c *gin.Context) {
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

	// Bulk delete experience
	err := h.experienceService.BulkDeleteExperiences(c.Request.Context(), idsInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Experiences deleted successfully")
}

// validateAndExtractExperienceInput handles validation and extraction of experience input fields
func (h *ExperienceHandler) validateAndExtractExperienceInput(c *gin.Context, experienceInput interface{}) error {
	// Extract JSON payload
	jsonPayload := c.PostForm("payload")
	if jsonPayload == "" {
		return errors.New(
			errors.ErrValidation,
			"Experience payload is required",
			nil,
		)
	}

	// Unmarshal JSON payload
	if err := json.Unmarshal([]byte(jsonPayload), experienceInput); err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid experience payload format: "+err.Error(),
			err,
		)
	}

	return nil
}

// extractFileHeaders handles extracting file headers from request
func (h *ExperienceHandler) extractFileHeaders(c *gin.Context, fieldName string, maxSizeInMB int64) ([]*multipart.FileHeader, error) {
	form, _ := c.MultipartForm()
	var fileHeaders []*multipart.FileHeader

	if form != nil && len(form.File[fieldName]) > 0 {
		fileHeaders = form.File[fieldName]
	} else {
		// fallback to single file
		f, err := c.FormFile(fieldName)
		if err == nil {
			fileHeaders = []*multipart.FileHeader{f}
		}
	}

	// Validate file sizes
	for _, fileHeader := range fileHeaders {
		if fileHeader.Size > maxSizeInMB*1024*1024 {
			return nil, errors.New(
				errors.ErrValidation,
				fmt.Sprintf("%s files must be less than %dMB each", fieldName, maxSizeInMB),
				nil,
			)
		}
	}

	return fileHeaders, nil
}
