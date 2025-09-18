package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/internal/middleware"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
)

// RegisterTechStackRoutes sets up routes for tech stack operations
func RegisterTechStackRoutes(
	r *gin.RouterGroup,
	techStackHandler *tech_stack.TechStackHandler,
	routerMiddleware *middleware.Middleware,
) {
	// Create a route group for tech stacks
	techStacks := r.Group("/tech-stacks")
	{
		// Create a new tech stack
		techStacks.POST("",
			routerMiddleware.VerifyJWT(),
			techStackHandler.CreateTechStack,
		)

		// List tech stacks
		techStacks.GET("",
			techStackHandler.ListTechStacks,
		)

		// Get a specific tech stack by ID
		techStacks.GET("/:id",
			techStackHandler.GetTechStackByID,
		)

		// Update a tech stack
		techStacks.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			techStackHandler.UpdateTechStack,
		)

		// Delete a tech stack
		techStacks.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			techStackHandler.DeleteTechStack,
		)

		// Bulk create tech stacks
		techStacks.POST("/bulk",
			routerMiddleware.VerifyJWT(),
			techStackHandler.BulkCreateTechStacks,
		)

		// Bulk update tech stacks
		techStacks.PUT("/bulk",
			routerMiddleware.VerifyJWT(),
			techStackHandler.BulkUpdateTechStacks,
		)

		// Bulk delete tech stacks
		techStacks.DELETE("/bulk",
			routerMiddleware.VerifyJWT(),
			techStackHandler.BulkDeleteTechStacks,
		)
	}
}
