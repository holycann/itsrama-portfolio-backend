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
	// Create a route group for experience
	experience := r.Group("/experience")
	{
		// Create a new experience
		experience.POST("",
			routerMiddleware.VerifyJWT(),
			experienceHandler.CreateExperience,
		)

		// List experience
		experience.GET("",
			experienceHandler.ListExperiences,
		)

		// Get experience by company
		experience.GET("/company/:company",
			experienceHandler.GetExperiencesByCompany,
		)

		// Get a specific experience by ID
		experience.GET("/:id",
			experienceHandler.GetExperienceByID,
		)

		// Update an experience
		experience.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			experienceHandler.UpdateExperience,
		)

		// Delete an experience
		experience.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			experienceHandler.DeleteExperience,
		)

		// Search experience
		experience.GET("/search",
			experienceHandler.SearchExperiences,
		)

		// Bulk create experience
		experience.POST("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkCreateExperiences,
		)

		// Bulk update experience
		experience.PUT("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkUpdateExperiences,
		)

		// Bulk delete experience
		experience.DELETE("/bulk",
			routerMiddleware.VerifyJWT(),
			experienceHandler.BulkDeleteExperiences,
		)
	}
}
