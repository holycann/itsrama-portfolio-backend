package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/handlers"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/supabase-community/supabase-go"
)

func RegisterEventRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	eventRepository := repositories.NewEventRepository(supabaseClient, *repositories.DefaultEventConfig())
	eventService := services.NewEventService(eventRepository)
	eventHandler := handlers.NewEventHandler(eventService, logger)

	event := r.Group("/events")
	{
		event.POST("/", eventHandler.CreateEvent)
		event.GET("/", eventHandler.ListEvent)
		event.GET("/search", eventHandler.SearchEvents)
		event.PUT("/:id", eventHandler.UpdateEvent)
		event.DELETE("/:id", eventHandler.DeleteEvent)
	}
}

func RegisterLocalRepositoryRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	localStoryRepository := repositories.NewLocalStoryRepository(supabaseClient, *repositories.DefaultLocalStoryConfig())
	localStoryService := services.NewLocalStoryService(localStoryRepository)
	localStoryHandler := handlers.NewLocalStoryHandler(localStoryService, logger)

	localStory := r.Group("/local-story")
	{
		localStory.POST("/", localStoryHandler.CreateLocalStory)
		localStory.GET("/", localStoryHandler.ListLocalStories)
		localStory.GET("/search", localStoryHandler.SearchLocalStories)
		localStory.PUT("/:id", localStoryHandler.UpdateLocalStory)
		localStory.DELETE("/:id", localStoryHandler.DeleteLocalStory)
	}
}
