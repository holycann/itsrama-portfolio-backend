package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/handlers"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	placeRepo "github.com/holycann/cultour-backend/internal/place/repositories"
	placeSvc "github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/supabase"
)

func RegisterEventRoutes(
	r *gin.Engine,
	appLogger *logger.Logger,
	supabaseClient *supabase.SupabaseClient,
	supabaseStorage *supabase.SupabaseStorage,
	routeMiddleware *middleware.Middleware,
) {
	eventRepository := repositories.NewEventRepository(supabaseClient.GetClient())
	placeRepository := placeRepo.NewLocationRepository(supabaseClient.GetClient())
	placeService := placeSvc.NewLocationService(placeRepository)
	eventService := services.NewEventService(eventRepository, placeService, supabaseStorage)
	eventHandler := handlers.NewEventHandler(eventService, appLogger)

	event := r.Group("/events")
	{
		// Only allow users with "admin" or "user" role or "penjelajah" or "warlok" badge to create/update/delete
		event.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.CreateEvent,
		)
		event.GET("", eventHandler.ListEvent)
		event.GET("/search", eventHandler.SearchEvents)
		event.GET("/trending", eventHandler.TrendingEvents)
		event.GET("/:id", eventHandler.GetEventByID)
		event.POST("/:id/views", eventHandler.UpdateEventViews)
		event.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.UpdateEvent,
		)
		event.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.DeleteEvent,
		)
	}
}

func RegisterLocalStoryRoutes(
	r *gin.Engine,
	appLogger *logger.Logger,
	supabaseClient *supabase.SupabaseClient,
	routeMiddleware *middleware.Middleware,
) {
	localStoryRepository := repositories.NewLocalStoryRepository(supabaseClient.GetClient())
	localStoryService := services.NewLocalStoryService(localStoryRepository)
	localStoryHandler := handlers.NewLocalStoryHandler(localStoryService, appLogger)

	localStory := r.Group("/local-stories")
	{
		// Only allow users with "admin" or "user" role or "penjelajah" or "warlok" badge to create/update/delete
		localStory.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			localStoryHandler.CreateLocalStory,
		)
		localStory.GET("", localStoryHandler.ListLocalStories)
		localStory.GET("/search", localStoryHandler.SearchLocalStories)
		localStory.GET("/:id", localStoryHandler.ListLocalStories) // Use ListLocalStories with ID filter
		localStory.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			localStoryHandler.UpdateLocalStory,
		)
		localStory.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			localStoryHandler.DeleteLocalStory,
		)
	}
}
