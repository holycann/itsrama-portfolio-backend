package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/handlers"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/supabase-community/supabase-go"
)

func RegisterEventRoutes(
	r *gin.Engine,
	appLogger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware *middleware.Middleware,
) {
	eventRepository := repositories.NewEventRepository(supabaseClient, *repositories.DefaultEventConfig())
	eventService := services.NewEventService(eventRepository)
	eventHandler := handlers.NewEventHandler(eventService, appLogger)

	event := r.Group("/events")
	{
		// Only allow users with "admin" or "user" role or "penjelajah" or "warlok" badge to create/update/delete
		event.POST("/",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.CreateEvent,
		)
		event.GET("/", eventHandler.ListEvent)
		event.GET("/search", eventHandler.SearchEvents)
		event.GET("/trending", eventHandler.TrendingEvents)
		event.GET("/:id", eventHandler.GetEventByID) // detail by id
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
	supabaseClient *supabase.Client,
	routeMiddleware *middleware.Middleware,
) {
	localStoryRepository := repositories.NewLocalStoryRepository(supabaseClient, *repositories.DefaultLocalStoryConfig())
	localStoryService := services.NewLocalStoryService(localStoryRepository)
	localStoryHandler := handlers.NewLocalStoryHandler(localStoryService, appLogger)

	localStory := r.Group("/local-stories")
	{
		// Only allow users with "admin" or "user" role or "penjelajah" or "warlok" badge to create/update/delete
		localStory.POST("/",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			localStoryHandler.CreateLocalStory,
		)
		localStory.GET("/", localStoryHandler.ListLocalStories)
		localStory.GET("/search", localStoryHandler.SearchLocalStories)
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
