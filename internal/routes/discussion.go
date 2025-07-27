package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/handlers"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/supabase-community/supabase-go"
)

func RegisterThreadRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware *middleware.Middleware,
) {
	threadRepository := repositories.NewThreadRepository(supabaseClient)
	threadService := services.NewThreadService(threadRepository)
	threadHandler := handlers.NewThreadHandler(threadService, logger)

	thread := r.Group("/threads")
	{
		thread.POST("", routeMiddleware.VerifyJWT(), threadHandler.CreateThread)
		thread.GET("", threadHandler.ListThreads)
		thread.GET("/search", routeMiddleware.VerifyJWT(), threadHandler.SearchThreads)
		thread.GET("/:id", routeMiddleware.VerifyJWT(), threadHandler.GetThreadByID)
		thread.PUT("/:id", routeMiddleware.VerifyJWT(), threadHandler.UpdateThread)
		thread.DELETE("/:id", routeMiddleware.VerifyJWT(), threadHandler.DeleteThread)
	}
}

func RegisterMessageRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware *middleware.Middleware,
) {
	messageRepository := repositories.NewMessageRepository(supabaseClient)
	messageService := services.NewMessageService(messageRepository)
	messageHandler := handlers.NewMessageHandler(messageService, logger)

	message := r.Group("/messages")
	{
		message.POST("", routeMiddleware.VerifyJWT(), messageHandler.CreateMessage)
		message.GET("", messageHandler.ListMessages)
		message.GET("/search", messageHandler.SearchMessages)
		message.GET("/:id", routeMiddleware.VerifyJWT(), messageHandler.ListMessages) // Use ListMessages with ID filter
		message.PUT("/:id", routeMiddleware.VerifyJWT(), messageHandler.UpdateMessage)
		message.DELETE("/:id", routeMiddleware.VerifyJWT(), messageHandler.DeleteMessage)
	}
}
