package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/gemini"
	"github.com/holycann/cultour-backend/internal/middleware"
)

func SetupGeminiRoutes(
	router *gin.Engine,
	geminiHandler *gemini.GeminiHandler,
	routeMiddleware *middleware.Middleware,
) {
	aiGroup := router.Group("/ai")
	{
		// Create a new AI chat session (auth required)
		aiGroup.POST("/chat/session",
			routeMiddleware.VerifyJWT(),
			geminiHandler.CreateChatSession,
		)

		// Send a message in an existing chat session (auth required)
		aiGroup.POST("/chat/:sessionID/message",
			routeMiddleware.VerifyJWT(),
			geminiHandler.SendMessage,
		)

		// Generate an AI-powered event description from title/prompt (auth required)
		aiGroup.POST("/events/description",
			routeMiddleware.VerifyJWT(),
			geminiHandler.GenerateEventDescription,
		)
	}
}
