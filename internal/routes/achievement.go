package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/achievement/handlers"
	"github.com/holycann/cultour-backend/internal/achievement/repositories"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/supabase-community/supabase-go"
)

// RegisterBadgeRoutes sets up routes for badge-related operations
func RegisterBadgeRoutes(
	r *gin.Engine,
	supabaseClient *supabase.Client,
	routerMiddleware *middleware.Middleware,
	appLogger *logger.Logger,
) {
	// Initialize repository
	badgeRepository := repositories.NewSupabaseBadgeRepository(supabaseClient, "badges")

	// Initialize service
	badgeService := services.NewBadgeService(badgeRepository)

	// Initialize handler
	badgeHandler := handlers.NewBadgeHandler(badgeService, appLogger)

	// Create a route group for badges
	badges := r.Group("/badges")
	{
		// Create a new badge (admin only)
		badges.POST("/",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.CreateBadge,
		)

		// List badges (optional authentication)
		badges.GET("/",
			routerMiddleware.VerifyJWT(),
			badgeHandler.ListBadges,
		)

		// Count total badges (optional authentication)
		badges.GET("/count",
			routerMiddleware.VerifyJWT(),
			badgeHandler.CountBadges,
		)

		// Get, update, or delete a specific badge by ID (admin only)
		badges.GET("/:id",
			routerMiddleware.VerifyJWT(),
			badgeHandler.GetBadgeByID,
		)
		badges.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.UpdateBadge,
		)
		badges.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			badgeHandler.DeleteBadge,
		)
	}
}
