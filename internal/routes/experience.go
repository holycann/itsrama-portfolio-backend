package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/internal/experience"
	"github.com/holycann/itsrama-portfolio-backend/internal/middleware"
)

// RegisterExperienceRoutes sets up routes for experience operations
func RegisterExperienceRoutes(
	r *gin.RouterGroup,
	experienceHandler *experience.ExperienceHandler,
	routerMiddleware *middleware.Middleware,
) {
	// Create a route group for experiences
	experiences := r.Group("/experiences")
	{
		// Create a new experience
		experiences.POST("",
			routerMiddleware.VerifyJWT(),
			experienceHandler.CreateExperience,
		)

		// List experiences
		experiences.GET("",
			experienceHandler.ListExperiences,
		)

		// Get a specific experience by ID
		experiences.GET("/:id",
			experienceHandler.GetExperienceByID,
		)

		// Update an experience
		experiences.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			experienceHandler.UpdateExperience,
		)

		// Delete an experience
		experiences.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			experienceHandler.DeleteExperience,
		)

		// Search experiences
		experiences.GET("/search",
			experienceHandler.SearchExperiences,
		)

		// Bulk create experiences
		experiences.POST("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkCreateExperiences,
		)

		// Bulk update experiences
		experiences.PUT("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkUpdateExperiences,
		)

		// Bulk delete experiences
		experiences.DELETE("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkDeleteExperiences,
		)
	}
}
