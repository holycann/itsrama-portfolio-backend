package gemini

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/response"
)

// MessageRepo is the interface for message/chat repository
type MessageRepo interface {
	// ListByThreadID retrieves a list of messages by thread ID, limit, and offset
	ListByThreadID(ctx context.Context, threadID string, limit int, offset int) ([]*Message, error)
}

// Message is a struct representing a chat message
type Message struct {
	ID        string
	ThreadID  string
	UserID    string
	Content   string
	CreatedAt time.Time
}

// EventDetail and EventSummary are helper structs for context
// (implementation of getEventDetail and getRelatedEvents can use event and user services)

type EventDetail struct {
	ID          string
	Name        string
	Description string
	StartDate   time.Time
	City        string
	Province    string
	UserEmail   string
}

type EventSummary struct {
	ID        string
	Name      string
	StartDate time.Time
	City      string
}

// Logger is a simple interface for logging
type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

// AIRequest represents the AI request structure
type AIRequest struct {
	// Prompt is the text input from the user for the AI
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

	// Metadata about the generated AI content
	Metadata struct {
		// Length of the generated text
		Length int `json:"length" example:"250"`

		// Tokens used in the generation process
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
		// Log error for debugging
		a.logger.Error("Failed to generate AI text",
			"error", err,
			"prompt", req.Prompt,
		)

		// Return structured error response
		response.InternalServerError(c, "AI generation failed", err.Error())
		return
	}

	// Prepare response
	resp := AIResponse{
		Response: result,
	}
	resp.Metadata.Length = len(resp.Response)
	// Note: Actual token count requires additional implementation

	// Send successful response
	response.SuccessOK(c, resp, "AI text generated successfully")
}

var aiCache = make(map[string]string) // key: hash(context+prompt), value: response

// AskEventAI handles AI text generation requests with event context
// @Summary Generate AI text response for a specific event
// @Description Generate text using Gemini AI with event context and chat history
// @Tags AI
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param thread_id query string false "Thread ID for chat history"
// @Param request body AIRequest true "AI Text Generation Request"
// @Success 200 {object} AIResponse "Successful AI text generated"
// @Failure 400 {object} response.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} response.ErrorResponse "Event not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error during AI generation"
// @Router /ask/event/{id} [post]
func (a *AIClient) AskEventAI(c *gin.Context) {
	var req AIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}
	eventID := c.Param("id")
	if eventID == "" {
		response.BadRequest(c, "Event ID is required", nil)
		return
	}
	threadID := c.Query("thread_id")

	// Get event details and other events (summary)
	eventRepo := repositories.NewEventRepository(a.supabaseClient.GetClient(), *repositories.DefaultEventConfig())
	eventService := services.NewEventService(eventRepo)
	event, err := eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil || event == nil {
		response.NotFound(c, "Event not found", err)
		return
	}
	relatedEvents, err := eventService.ListRelatedEvents(c.Request.Context(), eventID, 3)
	if err != nil {
		a.logger.Error("Failed to fetch related events", "error", err, "eventID", eventID)
		// Continue without related events instead of returning
	}

	// Get chat history if available
	chatHistory := ""
	if threadID != "" {
		msgs, _ := a.getChatHistory(threadID, 10)
		if len(msgs) > 0 {
			chatHistory = "[CHAT HISTORY]\n"
			for _, m := range msgs {
				chatHistory += m + "\n"
			}
		}
	}

	// Build context string
	contextStr := a.buildEventContext(&EventDetail{
		Name:        event.Name,
		Description: event.Description,
		StartDate:   event.StartDate,
	}, func(events []*models.Event) []EventSummary {
		summaries := make([]EventSummary, len(events))
		for i, e := range events {
			summaries[i] = EventSummary{
				Name:      e.Name,
				StartDate: e.StartDate,
			}
		}
		return summaries
	}(relatedEvents))
	finalPrompt := contextStr + "\n" + chatHistory + "\n[PROMPT USER]\n" + req.Prompt

	// Check cache
	cacheKey := a.hashPrompt(finalPrompt)
	if cached, ok := aiCache[cacheKey]; ok {
		resp := AIResponse{Response: cached}
		resp.Metadata.Length = len(resp.Response)
		response.SuccessOK(c, resp, "AI text generated successfully (cache)")
		return
	}

	ctx, cancel := a.prepareContext()
	defer cancel()
	result, err := a.GenerateText(ctx, finalPrompt)
	if err != nil {
		a.logger.Error("Failed to generate AI text", "error", err, "prompt", req.Prompt)
		response.InternalServerError(c, "AI generation failed", err.Error())
		return
	}
	aiCache[cacheKey] = result
	resp := AIResponse{Response: result}
	resp.Metadata.Length = len(resp.Response)
	response.SuccessOK(c, resp, "AI text generated successfully")
}

// buildEventContext builds a context string for the AI prompt
func (a *AIClient) buildEventContext(event *EventDetail, related []EventSummary) string {
	ctx := "[CONTEXT EVENT]\n"
	ctx += "Event Name: " + event.Name + "\n"
	ctx += "Description: " + event.Description + "\n"
	ctx += "Date: " + event.StartDate.Format("2006-01-02") + "\n"
	ctx += "Location: " + event.City + ", " + event.Province + "\n"
	ctx += "Creator User: " + event.UserEmail + "\n"
	ctx += "\n[OTHER EVENTS CONTEXT]\n"
	for _, e := range related {
		ctx += "- " + e.Name + ", " + e.StartDate.Format("2006-01-02") + ", " + e.City + "\n"
	}
	return ctx
}

// getChatHistory fetches last N messages as string slice (stub, implementation can use message repo)
func (a *AIClient) getChatHistory(threadID string, limit int) ([]string, error) {
	if a.messageRepo == nil {
		return nil, fmt.Errorf("messageRepo not set")
	}
	ctx := context.Background()
	msgs, err := a.messageRepo.ListByThreadID(ctx, threadID, limit, 0)
	if err != nil {
		return nil, err
	}
	var history []string
	for _, m := range msgs {
		role := "User"
		if m.UserID == "ai" { // assume AI message uses user_id "ai"
			role = "AI"
		}
		history = append(history, fmt.Sprintf("%s: %s", role, m.Content))
	}
	return history, nil
}

// hashPrompt creates a simple hash from the prompt for cache key
func (a *AIClient) hashPrompt(s string) string {
	// Simple: can be replaced with a stronger hash if needed
	return fmt.Sprintf("%x", len(s)) + "-" + s[:min(16, len(s))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// prepareContext creates a context with timeout for AI requests
func (a *AIClient) prepareContext() (context.Context, context.CancelFunc) {
	// Default timeout of 30 seconds
	return context.WithTimeout(context.Background(), 30*time.Second)
}
