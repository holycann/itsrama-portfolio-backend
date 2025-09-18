package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/internal/middleware"
	"github.com/holycann/itsrama-portfolio-backend/internal/project"
)

// RegisterProjectRoutes sets up routes for project operations
func RegisterProjectRoutes(
	r *gin.RouterGroup,
	projectHandler *project.ProjectHandler,
	routerMiddleware *middleware.Middleware,
) {
	// Create a route group for projects
	projects := r.Group("/projects")
	{
		// Create a new project
		projects.POST("",
			routerMiddleware.VerifyJWT(),
			projectHandler.CreateProject,
		)

		// List projects
		projects.GET("",
			projectHandler.ListProjects,
		)

		// Get a specific project by ID
		projects.GET("/:id",
			projectHandler.GetProjectByID,
		)

		// Update a project
		projects.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			projectHandler.UpdateProject,
		)

		// Delete a project
		projects.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			projectHandler.DeleteProject,
		)

		// Search projects
		projects.GET("/search",
			projectHandler.SearchProjects,
		)

		// Bulk create projects
		projects.POST("/bulk",
			routerMiddleware.VerifyJWT(),
			projectHandler.BulkCreateProjects,
		)

		// Bulk update projects
		projects.PUT("/bulk",
			routerMiddleware.VerifyJWT(),
			projectHandler.BulkUpdateProjects,
		)

		// Bulk delete projects
		projects.DELETE("/bulk",
			routerMiddleware.VerifyJWT(),
			projectHandler.BulkDeleteProjects,
		)
	}
}
