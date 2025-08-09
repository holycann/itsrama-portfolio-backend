package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/achievement/handlers"
	"github.com/holycann/cultour-backend/internal/middleware"
)

// RegisterBadgeRoutes sets up routes for badge-related operations
// @Description Configures HTTP routes for badge management with authentication and authorization
// @Tags Badges
func RegisterBadgeRoutes(
	r *gin.Engine,
	badgeHandler *handlers.BadgeHandler,
	routerMiddleware *middleware.Middleware,
) {
	// Create a route group for badges
	badges := r.Group("/badges")
	{
		// Create a new badge (admin only)
		// @Summary Create a new badge
		// @Description Allows administrators to create a new achievement badge
		// @Tags Badges
		badges.POST("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.CreateBadge,
		)

		// List badges
		// @Summary Retrieve a list of badges
		// @Description Fetches paginated list of badges with optional filtering
		// @Tags Badges
		badges.GET("",
			routerMiddleware.VerifyJWT(),
			badgeHandler.ListBadges,
		)

		// Search badges
		// @Summary Search badges by query
		// @Description Performs full-text search across badge fields
		// @Tags Badges
		badges.GET("/search",
			routerMiddleware.VerifyJWT(),
			badgeHandler.SearchBadges,
		)

		// Count total badges
		// @Summary Count total number of badges
		// @Description Returns the total count of badges in the system
		// @Tags Badges
		badges.GET("/count",
			routerMiddleware.VerifyJWT(),
			badgeHandler.CountBadges,
		)

		// Get a specific badge by ID
		// @Summary Retrieve a specific badge
		// @Description Fetches detailed information about a badge by its unique identifier
		// @Tags Badges
		badges.GET("/:id",
			routerMiddleware.VerifyJWT(),
			badgeHandler.GetBadgeByID,
		)

		// Update a badge (admin only)
		// @Summary Update an existing badge
		// @Description Allows administrators to modify badge details
		// @Tags Badges
		badges.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.UpdateBadge,
		)

		// Delete a badge (admin only)
		// @Summary Delete a badge
		// @Description Allows administrators to remove a badge from the system
		// @Tags Badges
		badges.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.DeleteBadge,
		)
	}
}
