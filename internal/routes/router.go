package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/gemini"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(router *gin.Engine, geminiClient gemini.AIClient) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/ask", geminiClient.AskAI)

	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status":    "healthy",
	// 		"timestamp": time.Now().UTC(),
	// 	})
	// })
}
