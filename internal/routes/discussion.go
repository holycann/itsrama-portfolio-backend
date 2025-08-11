package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/handlers"
	"github.com/holycann/cultour-backend/internal/middleware"
)

func RegisterThreadRoutes(
	r *gin.Engine,
	threadHandler *handlers.ThreadHandler,
	routeMiddleware *middleware.Middleware,
) {
	thread := r.Group("/threads")
	{
		// Create a new thread (requires authentication)
		thread.POST("",
			routeMiddleware.VerifyJWT(),
			threadHandler.CreateThread,
		)

		// List threads with optional filtering
		thread.GET("",
			routeMiddleware.VerifyJWT(),
			threadHandler.ListThreads,
		)

		// Search threads by query (requires authentication)
		thread.GET("/search",
			routeMiddleware.VerifyJWT(),
			threadHandler.SearchThreads,
		)

		// Get thread details by ID (requires authentication)
		thread.GET("/:id",
			routeMiddleware.VerifyJWT(),
			threadHandler.GetThreadByID,
		)

		// Update an existing thread (requires authentication)
		thread.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			threadHandler.UpdateThread,
		)

		// Delete a thread (requires authentication)
		thread.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			threadHandler.DeleteThread,
		)

		// Join a thread (requires authentication)
		thread.POST("/:id/join",
			routeMiddleware.VerifyJWT(),
			threadHandler.JoinThread,
		)

		// Get thread by associated event ID (requires authentication)
		thread.GET("/event/:event_id",
			routeMiddleware.VerifyJWT(),
			threadHandler.GetThreadByEvent,
		)
	}
}

func RegisterMessageRoutes(
	r *gin.Engine,
	messageHandler *handlers.MessageHandler,
	routeMiddleware *middleware.Middleware,
) {
	message := r.Group("/messages")
	{
		// Create a new message (requires authentication)
		message.POST("",
			routeMiddleware.VerifyJWT(),
			messageHandler.CreateMessage,
		)

		// List messages with optional filtering
		message.GET("",
			routeMiddleware.VerifyJWT(),
			messageHandler.ListMessages,
		)

		// Search messages by query (requires authentication)
		message.GET("/search",
			routeMiddleware.VerifyJWT(),
			messageHandler.SearchMessages,
		)

		// Get messages for a specific thread (requires authentication)
		message.GET("/thread/:thread_id",
			routeMiddleware.VerifyJWT(),
			messageHandler.GetMessagesByThread,
		)

		// Update an existing message (requires authentication)
		message.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			messageHandler.UpdateMessage,
		)

		// Delete a message (requires authentication)
		message.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			messageHandler.DeleteMessage,
		)
	}
}
