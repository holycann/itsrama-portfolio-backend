package project

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

type ProjectHandler struct {
	base.BaseHandler
	projectService ProjectService
}

func NewProjectHandler(projectService ProjectService, logger *logger.Logger) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler:    *base.NewBaseHandler(logger),
		projectService: projectService,
	}
}

// CreateProject creates a new project
// @Summary Create a new project
// @Description Create a new project with details and optional images
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param uploaded_images formData []file false "Project Images"
// @Param payload formData string true "Project Details in JSON format (See ProjectCreate Model)"
// @Success 200 {object} response.APIResponse{data=Project} "Project created successfully"
// @Failure 400 {object} response.APIResponse{data=ProjectCreate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var projectInput ProjectCreate

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
	if err := utils.ExtractFormDataPayload(c, &projectInput); err != nil {
		h.HandleError(c, err)
		return
	}

	uploadedImages, err := utils.ExtractFileHeaders(c, "uploaded_images", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(uploadedImages) > 0 {
		projectInput.UploadedImages = uploadedImages
	}

	// Create project via service
	project, err := h.projectService.CreateProject(c.Request.Context(), &projectInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, project, "Project created successfully")
}

// GetProjectByID retrieves a specific project
// @Summary Get a project by ID
// @Description Retrieve a project using its unique identifier
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} response.APIResponse{data=Project} "Project retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Project not found"
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Project ID is required",
			nil,
		))
		return
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, project, "Project retrieved successfully")
}

// UpdateProject updates an existing project
// @Summary Update a project
// @Description Update an existing project with new details and optional images
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Project ID"
// @Param uploaded_images formData []file false "Project Images"
// @Param payload formData string true "Project Update Details in JSON format (See ProjectUpdate Model)"
// @Success 200 {object} response.APIResponse{data=Project} "Project updated successfully"
// @Failure 400 {object} response.APIResponse{data=ProjectUpdate} "Bad Request"
// @Failure 404 {object} response.APIResponse "Project not found"
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var projectInput ProjectUpdate

	// Extract project ID from path parameter
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid project ID",
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
	projectInput.ID = projectID

	// Extract and validate form fields
	if err := utils.ExtractFormDataPayload(c, &projectInput); err != nil {
		h.HandleError(c, err)
		return
	}

	uploadedImages, err := utils.ExtractFileHeaders(c, "uploaded_images", 2)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	if len(uploadedImages) > 0 {
		projectInput.UploadedImages = uploadedImages
	}

	// Update project via service
	project, err := h.projectService.UpdateProject(c.Request.Context(), &projectInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, project, "Project updated successfully")
}

// DeleteProject deletes an existing project
// @Summary Delete a project
// @Description Delete a project by its unique identifier
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} response.APIResponse "Project deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Failure 404 {object} response.APIResponse "Project not found"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Project ID is required",
			nil,
		))
		return
	}

	err := h.projectService.DeleteProject(c.Request.Context(), projectID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Project deleted successfully")
}

// ListProjects retrieves a paginated list of projects
// @Summary List projects
// @Description Retrieve a paginated list of projects with optional filtering
// @Tags Projects
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param category query string false "Filter by project category"
// @Success 200 {object} response.APIResponse{data=[]Project} "Projects retrieved successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
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

	// Optional category filter
	slug := c.Query("slug")
	if slug != "" {
		opts.Filters = append(opts.Filters, base.FilterOption{
			Field:    "slug",
			Operator: base.OperatorEqual,
			Value:    slug,
		})
	}

	// List projects
	projects, err := h.projectService.ListProjects(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Count total projects for pagination
	total, err := h.projectService.CountProjects(c.Request.Context(), opts.Filters)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, projects, "Projects retrieved successfully",
		response.WithPagination(total, opts.Page, opts.PerPage))
}

// SearchProjects performs a full-text search on projects
// @Summary Search projects
// @Description Perform a full-text search on projects with pagination
// @Tags Projects
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.APIResponse{data=[]Project} "Projects search completed successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /projects/search [get]
func (h *ProjectHandler) SearchProjects(c *gin.Context) {
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

	projects, total, err := h.projectService.SearchProjects(c.Request.Context(), opts)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, projects, "Projects search completed successfully",
		response.WithPagination(total, opts.Page, opts.PerPage))
}

// BulkCreateProjects creates multiple projects in bulk
// @Summary Bulk create projects
// @Description Create multiple projects in a single request
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param uploaded_images formData []file false "Project Images"
// @Param payload formData string true "Project Details in JSON array format (See ProjectCreate Model)"
// @Success 200 {object} response.APIResponse{data=[]Project} "Projects created successfully"
// @Failure 400 {object} response.APIResponse{data=[]ProjectCreate} "Bad Request"
// @Failure 500 {object} response.APIResponse "Internal Server Error"
// @Router /projects/bulk [post]
func (h *ProjectHandler) BulkCreateProjects(c *gin.Context) {
	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrBadRequest,
			"Failed to parse multipart form",
			err,
		))
		return
	}

	var projectsInput []*ProjectCreate

	// Extract and validate form fields
	if err := utils.ExtractFormDataPayload(c, &projectsInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get logo files
	uploadedImages, err := utils.ExtractFileHeaders(c, "uploaded_images", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Match logo and images to projects if available
	if len(uploadedImages) > 0 {
		for i := range projectsInput {
			if i < len(uploadedImages) {
				projectsInput[i].UploadedImages = []*multipart.FileHeader{uploadedImages[i]}
			}
		}
	}

	// Bulk create projects
	projects, err := h.projectService.BulkCreateProjects(c.Request.Context(), projectsInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, projects, "Projects created successfully")
}

// BulkUpdateProjects updates multiple projects in bulk
// @Summary Bulk update projects
// @Description Update multiple projects in a single request
// @Tags Projects
// @Accept multipart/form-data
// @Produce json
// @Param uploaded_images formData []file false "Project Images"
// @Param payload formData string true "Project Update Details in JSON array format (See ProjectUpdate Model)"
// @Success 200 {object} response.APIResponse{data=[]Project} "Projects updated successfully"
// @Failure 400 {object} response.APIResponse{data=[]ProjectUpdate} "Bad Request"
// @Router /projects/bulk [put]
func (h *ProjectHandler) BulkUpdateProjects(c *gin.Context) {
	// Parse multipart form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrBadRequest,
			"Failed to parse multipart form",
			err,
		))
		return
	}

	var projectsInput []*ProjectUpdate

	// Extract and validate form fields
	if err := utils.ExtractFormDataPayload(c, &projectsInput); err != nil {
		h.HandleError(c, err)
		return
	}

	// Get uploaded images
	uploadedImages, err := utils.ExtractFileHeaders(c, "uploaded_images", 5)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Match images to projects if available
	if len(uploadedImages) > 0 {
		for i := range projectsInput {
			if i < len(uploadedImages) {
				projectsInput[i].UploadedImages = []*multipart.FileHeader{uploadedImages[i]}
			}
		}
	}

	// Bulk update projects
	updatedProjects, err := h.projectService.BulkUpdateProjects(c.Request.Context(), projectsInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, updatedProjects, "Projects updated successfully")
}

// BulkDeleteProjects deletes multiple projects in bulk
// @Summary Bulk delete projects
// @Description Delete multiple projects in a single request
// @Tags Projects
// @Accept json
// @Produce json
// @Param ids body []string true "Project IDs to delete"
// @Success 200 {object} response.APIResponse "Projects deleted successfully"
// @Failure 400 {object} response.APIResponse "Bad Request"
// @Router /projects/bulk [delete]
func (h *ProjectHandler) BulkDeleteProjects(c *gin.Context) {
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

	// Bulk delete projects
	err := h.projectService.BulkDeleteProjects(c.Request.Context(), idsInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, nil, "Projects deleted successfully")
}
