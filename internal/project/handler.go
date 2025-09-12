package project

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/pkg/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
	"github.com/holycann/itsrama-portfolio-backend/pkg/response"
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
	if err := h.validateAndExtractProjectInput(c, &projectInput); err != nil {
		h.HandleError(c, err)
		return
	}

	projectInput.UploadedImages = c.Request.MultipartForm.File["uploaded_images"]

	// Create project via service
	project, err := h.projectService.CreateProject(c.Request.Context(), &projectInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, project, "Project created successfully")
}

// GetProjectByID retrieves a specific project
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
	if err := h.validateAndExtractProjectInput(c, &projectInput); err != nil {
		h.HandleError(c, err)
		return
	}

	projectInput.UploadedImages = c.Request.MultipartForm.File["uploaded_images"]

	// Update project via service
	project, err := h.projectService.UpdateProject(c.Request.Context(), &projectInput)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, project, "Project updated successfully")
}

// DeleteProject deletes an existing project
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
func (h *ProjectHandler) BulkCreateProjects(c *gin.Context) {
	var projectsInput []*ProjectCreate

	// Bind input
	if err := c.ShouldBindJSON(&projectsInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
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
func (h *ProjectHandler) BulkUpdateProjects(c *gin.Context) {
	var projectsInput []*ProjectUpdate

	// Bind input
	if err := c.ShouldBindJSON(&projectsInput); err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid input",
			err,
		))
		return
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

// GetProjectsByCategory retrieves projects by category
func (h *ProjectHandler) GetProjectsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Category is required",
			nil,
		))
		return
	}

	projects, err := h.projectService.GetProjectsByCategory(c.Request.Context(), category)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, projects, "Projects retrieved successfully")
}

// validateAndExtractProjectInput handles validation and extraction of project input fields
func (h *ProjectHandler) validateAndExtractProjectInput(c *gin.Context, projectInput interface{}) error {
	// Extract JSON payload
	jsonPayload := c.PostForm("payload")
	if jsonPayload == "" {
		return errors.New(
			errors.ErrValidation,
			"Project payload is required",
			nil,
		)
	}

	// Unmarshal JSON payload
	if err := json.Unmarshal([]byte(jsonPayload), projectInput); err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid project payload format",
			err,
		)
	}

	return nil
}
