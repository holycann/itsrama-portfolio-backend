package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/handlers"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/supabase-community/supabase-go"
)

func RegisterThreadRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	threadRepository := repositories.NewThreadRepository(supabaseClient, *repositories.DefaultThreadConfig())
	threadService := services.NewThreadService(threadRepository)
	threadHandler := handlers.NewThreadHandler(threadService, logger)

	thread := r.Group("/threads")
	{
		thread.POST("/", threadHandler.CreateThread)
		thread.GET("/", threadHandler.ListThreads)
		thread.GET("/search", threadHandler.SearchThread)
		thread.PUT("/:id", threadHandler.UpdateThread)
		thread.DELETE("/:id", threadHandler.DeleteThread)
	}
}

func RegisterMessageRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	messageRepository := repositories.NewMessageRepository(supabaseClient, *repositories.DefaultMessageConfig())
	messageService := services.NewMessageService(messageRepository)
	messageHandler := handlers.NewMessageHandler(messageService, logger)

	message := r.Group("/messages")
	{
		message.POST("/", messageHandler.CreateMessage)
		message.GET("/", messageHandler.ListMessages)
		message.GET("/search", messageHandler.SearchMessages)
		message.PUT("/:id", messageHandler.UpdateMessage)
		message.DELETE("/:id", messageHandler.DeleteMessage)
	}
}
