package experience

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/internal/response"
	"github.com/holycann/itsrama-portfolio-backend/internal/utils"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
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
// @Summary Create a new experience
// @Description Create a new experience with logo and images
// @Tags Experiences
// @Accept multipart/form-data
// @Produce json
// @Param logo_image formData file false "Logo Image"
// @Param images formData file false "Experience Images"
// @Param payload formData string true "Experience Details in JSON format (See ExperienceCreate Model)"
// @Success 200 {object} response.APIResponse{data=Experience} "Experience created successfully"
// @Failure 400 {object} response.APIResponse{data=ExperienceCreate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /experiences [post]
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
	if err := utils.ExtractFormDataPayload(c, &experienceInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get logo file
	logoFileHeaders, err := utils.ExtractFileHeaders(c, "logo_image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(logoFileHeaders) > 0 {
		experienceInput.LogoImage = logoFileHeaders[0]
	}

	// Get image files
	imageFileHeaders, err := utils.ExtractFileHeaders(c, "images", 2)
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
// @Summary Get an experience by ID
// @Description Retrieve a specific experience using its unique identifier
// @Tags Experiences
// @Produce json
// @Param id path string true "Experience ID"
// @Success 200 {object} response.APIResponse{data=Experience} "Experience retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Experience not found"
// @Router /experiences/{id} [get]
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
// @Summary Update an experience
// @Description Update an existing experience with new details and optional files
// @Tags Experiences
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Experience ID"
// @Param logo_image formData file false "Logo Image"
// @Param images formData file false "Experience Images"
// @Param payload formData string true "Experience Details in JSON format (See ExperienceUpdate Model)"
// @Success 200 {object} response.APIResponse{data=Experience} "Experience updated successfully"
// @Failure 400 {object} response.APIResponse{data=ExperienceUpdate} "Bad Request"
// @Failure 404 {object} response.APIResponse "Experience not found"
// @Router /experiences/{id} [put]
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
	if err := utils.ExtractFormDataPayload(c, &experienceInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get logo file
	logoFileHeaders, err := utils.ExtractFileHeaders(c, "logo_image", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(logoFileHeaders) > 0 {
		experienceInput.LogoImage = logoFileHeaders[0]
	}

	// Get image files
	imageFileHeaders, err := utils.ExtractFileHeaders(c, "images", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(imageFileHeaders) > 0 {
		experienceInput.Images = imageFileHeaders
	}

	// Update experience via service
	experience, err := h.experienceService.UpdateExperience(c.Request.Context(), &experienceInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experience updated successfully")
}

// DeleteExperience deletes an existing experience
// @Summary Delete an experience
// @Description Delete an experience by its unique identifier
// @Tags Experiences
// @Produce json
// @Param id path string true "Experience ID"
// @Success 200 {object} response.APIResponse "Experience deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Experience not found"
// @Router /experiences/{id} [delete]
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
// @Summary List experiences
// @Description Retrieve a paginated list of experiences with optional filtering
// @Tags Experiences
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param company query string false "Filter by company name"
// @Param is_featured query string false "Filter by featured status"
// @Success 200 {object} response.APIResponse{data=[]Experience} "Experiences retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /experiences [get]
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

	// Optional company filter
	isFeatured := c.Query("is_featured")
	if isFeatured != "" {
		opts.Filters = append(opts.Filters, base.FilterOption{
			Field:    "is_featured",
			Operator: base.OperatorEqual,
			Value:    isFeatured,
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

// SearchExperiences performs a full-text search on experience
// @Summary Search experiences
// @Description Perform a full-text search on experiences with pagination
// @Tags Experiences
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.APIResponse{data=[]Experience} "Experiences search completed successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /experiences/search [get]
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
// @Summary Bulk create experiences
// @Description Create multiple experiences in a single request
// @Tags Experiences
// @Accept multipart/form-data
// @Produce json
// @Param logo_image formData []file false "Logo Image"
// @Param images formData []file false "Experience Images"
// @Param payload formData string true "Experience Details in JSON array format (See ExperienceCreate Model)"
// @Success 200 {object} response.APIResponse{data=[]Experience} "Experiences created successfully"
// @Failure 400 {object} response.APIResponse{data=[]ExperienceCreate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /experiences/bulk [post]
func (h *ExperienceHandler) BulkCreateExperiences(c *gin.Context) {
	var experiencesInput []*ExperienceCreate

	// Extract JSON payload
	if err := utils.ExtractFormDataPayload(c, &experiencesInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Extract logo images
	logoImages, err := utils.ExtractFileHeaders(c, "logo_image", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Extract experience images
	experienceImages, err := utils.ExtractFileHeaders(c, "images", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(logoImages) > 0 || len(experienceImages) > 0 {
		for i := range experiencesInput {
			// Match logo image to experience input if available
			if i < len(logoImages) {
				experiencesInput[i].LogoImage = logoImages[i]
			}

			// Match experience images to experience input if available
			if i < len(experienceImages) {
				experiencesInput[i].Images = []*multipart.FileHeader{experienceImages[i]}
			}
		}
	}

	// Bulk create experience
	experience, err := h.experienceService.BulkCreateExperiences(c.Request.Context(), experiencesInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, experience, "Experiences created successfully")
}

// @Description Bulk update experiences with logo and images
// @Tags Experiences
// @Accept multipart/form-data
// @Produce json
// @Param logo_image formData file false "Logo Image"
// @Param images formData file false "Experience Images"
// @Param payload formData string true "Experience Details in JSON array format (See ExperienceUpdate Model)"
// @Success 200 {object} response.APIResponse{data=[]Experience} "Experiences updated successfully"
// @Failure 400 {object} response.APIResponse{data=[]ExperienceUpdate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /experiences/bulk [put]
func (h *ExperienceHandler) BulkUpdateExperiences(c *gin.Context) {
	var experiencesInput []*ExperienceUpdate

	// Extract JSON payload
	if err := utils.ExtractFormDataPayload(c, &experiencesInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Extract logo images
	logoImages, err := utils.ExtractFileHeaders(c, "logo_image", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Extract experience images
	experienceImages, err := utils.ExtractFileHeaders(c, "images", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(logoImages) > 0 || len(experienceImages) > 0 {
		for i := range experiencesInput {
			// Match logo image to experience input if available
			if i < len(logoImages) {
				experiencesInput[i].LogoImage = logoImages[i]
			}

			// Match experience images to experience input if available
			if i < len(experienceImages) {
				experiencesInput[i].Images = []*multipart.FileHeader{experienceImages[i]}
			}
		}
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
// @Summary Bulk delete experiences
// @Description Delete multiple experiences in a single request
// @Tags Experiences
// @Accept json
// @Produce json
// @Param ids body []string true "Experience IDs to delete"
// @Success 200 {object} response.APIResponse "Experiences deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /experiences/bulk [delete]
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
