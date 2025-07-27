package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
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

// AIEventResponse represents the structured response from the AI for event context
type AIEventResponse struct {
	// Generated text from the AI
	Response string `json:"response" example:"Quantum computing is a type of computing that uses quantum-mechanical phenomena..."`

	// Related events for context
	RelatedEvents []EventSummary `json:"related_events"`

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
// @Failure 400 {object} response.APIResponse "Invalid request parameters"
// @Failure 500 {object} response.APIResponse "Internal server error during AI generation"
// @Router /ask [post]
func (a *AIClient) AskAI(c *gin.Context) {
	// Create request struct
	var req AIRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"validation_error": err.Error(),
		})
		response.BadRequest(c, "Invalid request", string(details), "")
		return
	}

	// Validate prompt length
	if len(req.Prompt) < 1 || len(req.Prompt) > 2000 {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"prompt_length": fmt.Sprintf("must be between 1 and 2000 characters, current length: %d", len(req.Prompt)),
		})
		response.BadRequest(c, "Invalid prompt length", string(details), "")
		return
	}

	// Prepare context with timeout
	ctx, cancel := a.prepareContext()
	defer cancel()

	// Generate AI response with advanced options
	generationOpts := &AIGenerationOptions{
		TraceID: c.GetString("request_id"), // Assuming request_id is set in middleware
		// Optional: Add more specific configuration if needed
		// Model: "gemini-2.0-flash", // Optional model override
		// Temperature: &customTemp,  // Optional temperature override
	}

	result, err := a.GenerateTextWithOptions(ctx, req.Prompt, generationOpts)
	if err != nil {
		// Log error for debugging
		a.logger.Error("Failed to generate AI text",
			"error", err,
			"prompt", req.Prompt,
		)

		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"generation_error": err.Error(),
			"prompt":           req.Prompt,
		})

		// Return structured error response
		response.InternalServerError(c, "AI generation failed", string(details), "")
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
// @Failure 400 {object} response.APIResponse "Invalid request parameters"
// @Failure 404 {object} response.APIResponse "Event not found"
// @Failure 500 {object} response.APIResponse "Internal server error during AI generation"
// @Router /ask/event/{id} [post]
func (a *AIClient) AskEventAI(c *gin.Context) {
	// Create request struct
	var req AIRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"validation_error": err.Error(),
		})
		response.BadRequest(c, "Invalid request", string(details), "")
		return
	}

	// Validate prompt length
	if len(req.Prompt) < 1 || len(req.Prompt) > 2000 {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"prompt_length": fmt.Sprintf("must be between 1 and 2000 characters, current length: %d", len(req.Prompt)),
		})
		response.BadRequest(c, "Invalid prompt length", string(details), "")
		return
	}

	// Extract event ID from path parameter
	eventID := c.Param("id")
	if eventID == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"event_id": "Event ID is required",
		})
		response.BadRequest(c, "Invalid event ID", string(details), "")
		return
	}

	// Prepare context with timeout
	ctx, cancel := a.prepareContext()
	defer cancel()

	// Fetch event details
	eventRepo := repositories.NewEventRepository(a.supabaseClient.GetClient())
	eventService := services.NewEventService(eventRepo)

	event, err := eventService.GetEventByID(ctx, eventID)
	if err != nil {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"event_fetch_error": err.Error(),
			"event_id":          eventID,
		})
		response.NotFound(c, "Event not found", string(details), "")
		return
	}

	// Prepare event context for AI prompt
	eventContext := fmt.Sprintf("Event Name: %s\nDescription: %s\nStart Date: %s",
		event.Name, event.Description, event.StartDate.Format(time.RFC3339))

	// Combine event context with user prompt
	fullPrompt := fmt.Sprintf("%s\n\nUser Query: %s", eventContext, req.Prompt)

	// Generate AI response with advanced options
	generationOpts := &AIGenerationOptions{
		TraceID: c.GetString("request_id"), // Assuming request_id is set in middleware
		// Optional: Add more specific configuration if needed
		// Model: "gemini-2.0-flash", // Optional model override
		// Temperature: &customTemp,  // Optional temperature override
	}

	result, err := a.GenerateTextWithOptions(ctx, fullPrompt, generationOpts)
	if err != nil {
		// Log error for debugging
		a.logger.Error("Failed to generate event AI text",
			"error", err,
			"event_id", eventID,
			"prompt", req.Prompt,
		)

		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"generation_error": err.Error(),
			"event_id":         eventID,
			"prompt":           req.Prompt,
		})

		// Return structured error response
		response.InternalServerError(c, "Event AI generation failed", string(details), "")
		return
	}

	// Fetch related events for additional context
	relatedEvents, err := eventService.ListEvents(ctx, repository.ListOptions{
		Filters: []repository.FilterOption{
			{Field: "city_id", Value: event.CityID},
		},
		Limit: 3,
	})
	if err != nil {
		a.logger.Warn("Failed to fetch related events",
			"error", err,
			"event_id", eventID)
	}

	// Prepare response
	resp := AIEventResponse{
		Response: result,
		RelatedEvents: func(events []models.Event) []EventSummary {
			summaries := make([]EventSummary, 0, len(events))
			for _, event := range events {
				summaries = append(summaries, EventSummary{
					ID:        event.ID.String(),
					Name:      event.Name,
					StartDate: event.StartDate,
				})
			}
			return summaries
		}(relatedEvents),
	}
	resp.Metadata.Length = len(resp.Response)

	// Send successful response
	response.SuccessOK(c, resp, "Event AI text generated successfully")
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
