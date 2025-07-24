package gemini

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/response"
)

// AIRequest represents the structure of an incoming AI request
type AIRequest struct {
	// Prompt is the user's input text for the AI
	// Example: "Explain quantum computing in simple terms"
	// Required: true
	// Min length: 1
	// Max length: 2000
	Prompt string `json:"prompt" binding:"required,min=1,max=2000" example:"Explain quantum computing in simple terms"`
}

// AIResponse represents the structured response from the AI
type AIResponse struct {
	// Generated text from the AI
	Response string `json:"response" example:"Quantum computing is a type of computing that uses quantum-mechanical phenomena..."`

	// Metadata about the AI-generated content
	Metadata struct {
		// Length of the generated text
		Length int `json:"length" example:"250"`

		// Tokens used in the generation
		TokensUsed int `json:"tokens_used" example:"60"`
	} `json:"metadata"`
}

// AskAI handles AI text generation requests
// @Summary Generate AI text response
// @Description Generate text using Gemini AI
// @Tags AI
// @Accept json
// @Produce json
// @Param request body AIRequest true "AI Text Generation Request"
// @Success 200 {object} AIResponse "Successful AI text generation"
// @Failure 400 {object} response.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error during AI generation"
// @Router /ask [post]
func (a *AIClient) AskAI(c *gin.Context) {
	// Create request struct
	var req AIRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	// Prepare context with timeout
	ctx, cancel := a.prepareContext()
	defer cancel()

	// Generate AI response
	result, err := a.GenerateText(ctx, req.Prompt)
	if err != nil {
		// Log the error for debugging purposes
		a.logger.Error("Failed to generate AI text",
			"error", err,
			"prompt", req.Prompt,
		)

		// Return a structured error response
		response.InternalServerError(c, "AI generation failed", err.Error())
		return
	}

	// Prepare response
	resp := AIResponse{
		Response: result,
	}
	resp.Metadata.Length = len(resp.Response)
	// Note: Actual token count would require additional implementation

	// Send successful response
	response.SuccessOK(c, resp, "AI text generated successfully")
}

// prepareContext creates a context with timeout for AI requests
func (a *AIClient) prepareContext() (context.Context, context.CancelFunc) {
	// Default timeout of 30 seconds
	return context.WithTimeout(context.Background(), 30*time.Second)
}
