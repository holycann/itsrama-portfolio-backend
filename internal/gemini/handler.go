package gemini

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/configs"
	culturalModels "github.com/holycann/cultour-backend/internal/cultural/models"
	culturalServices "github.com/holycann/cultour-backend/internal/cultural/services"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"google.golang.org/genai"
)

// ErrorResponse defines a standardized error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// GeminiHandler manages AI-powered interactions
type GeminiHandler struct {
	aiService    *AIService
	eventService culturalServices.EventService
	cityService  placeServices.CityService
	userService  userServices.UserProfileService
	config       *configs.Config
}

// NewGeminiHandler creates a new handler for Gemini AI interactions
func NewGeminiHandler(
	config *configs.Config,
	eventService culturalServices.EventService,
	cityService placeServices.CityService,
	userService userServices.UserProfileService,
) (*GeminiHandler, error) {
	// Initialize AI service
	aiService, err := NewAIService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI service: %v", err)
	}

	return &GeminiHandler{
		aiService:    aiService,
		eventService: eventService,
		cityService:  cityService,
		userService:  userService,
		config:       config,
	}, nil
}

// validateFeatureScope checks if the request is within supported application features
func (h *GeminiHandler) validateFeatureScope(feature string) error {
	supportedFeatures := map[string]bool{
		"event_exploration": true,
		"ai_assistant":      true,
		"discussion_forum":  true,
		"warlok_creation":   true,
	}

	if !supportedFeatures[feature] {
		return fmt.Errorf("feature %s is not currently supported", feature)
	}
	return nil
}

// CreateChatSession handles chat session creation with strict feature validation
// @Summary Create a new chat session
// @Description Creates a new AI chat session for a user, optionally with an event context
// @Tags AI
// @Accept json
// @Produce json
// @Param request body CreateChatSessionRequest true "Chat Session Creation Request"
// @Success 200 {object} CreateChatSessionResponse "Successfully created chat session"
// @Failure 400 {object} ErrorResponse "Invalid request or user not found"
// @Failure 403 {object} ErrorResponse "Feature not supported"
// @Failure 500 {object} ErrorResponse "Internal server error during session creation"
// @Router /ai/chat/session [post]
func (h *GeminiHandler) CreateChatSession(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("ai_assistant"); err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FEATURE_NOT_SUPPORTED",
			Message: "AI assistant feature is currently restricted",
			Details: err.Error(),
		})
		return
	}

	var req CreateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format. Please provide a valid user ID.",
			Details: fmt.Sprintf("Binding error: %v", err),
		})
		return
	}

	// Validate user
	ctx := context.Background()
	userProfile, err := h.userService.GetProfileByUserID(ctx, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "Pengguna tidak ditemukan. Silakan masuk atau daftar terlebih dahulu.",
			Details: fmt.Sprintf("User retrieval error: %v", err),
		})
		return
	}

	// Create chat session with event context if provided
	sessionID, err := h.aiService.CreateChatSession(userProfile.ID.String(), req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "SESSION_CREATE_FAILED",
			Message: "Gagal membuat sesi. Silakan coba lagi.",
			Details: fmt.Sprintf("Session creation error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, CreateChatSessionResponse{
		SessionID: sessionID,
	})
}

// SendMessage processes user messages with flexible topic handling
// @Summary Send a message in an AI chat session
// @Description Sends a user message to the AI and retrieves the AI's response
// @Tags AI
// @Accept json
// @Produce json
// @Param sessionID path string true "Chat Session ID"
// @Param request body SendMessageRequest true "Message Request"
// @Success 200 {object} SendMessageResponse "Successfully processed message"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 403 {object} ErrorResponse "Feature not supported"
// @Failure 500 {object} ErrorResponse "Error processing message"
// @Router /ai/chat/{sessionID}/message [post]
func (h *GeminiHandler) SendMessage(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("ai_assistant"); err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FEATURE_NOT_SUPPORTED",
			Message: "AI assistant feature is currently restricted",
			Details: err.Error(),
		})
		return
	}

	// Extract session ID from path
	sessionID := c.Param("sessionID")

	// Parse request body
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Format pesan tidak valid. Mohon periksa kembali.",
			Details: fmt.Sprintf("Binding error: %v", err),
		})
		return
	}

	// Retrieve session details to get user and event context
	session, err := h.aiService.GetChatSession(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "SESSION_NOT_FOUND",
			Message: "Sesi chat tidak ditemukan. Silakan buat sesi baru.",
			Details: fmt.Sprintf("Session retrieval error: %v", err),
		})
		return
	}

	// Fetch user details
	ctx := context.Background()
	user, err := h.userService.GetProfileByID(ctx, session.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "Pengguna tidak ditemukan. Silakan masuk atau daftar terlebih dahulu.",
			Details: fmt.Sprintf("User retrieval error: %v", err),
		})
		return
	}

	// Fetch user profile
	userProfile, err := h.userService.GetProfileByID(ctx, user.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "PROFILE_NOT_FOUND",
			Message: "Profil pengguna tidak ditemukan.",
			Details: fmt.Sprintf("User profile retrieval error: %v", err),
		})
		return
	}

	// Prepare context for AI interaction
	var event *culturalModels.ResponseEvent
	var eventContext AIInteractionContext

	// If session has an event ID, fetch event details
	if session.EventID != nil {
		event, err = h.eventService.GetEventByID(ctx, *session.EventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "EVENT_NOT_FOUND",
				Message: "Event tidak ditemukan. Silakan periksa kembali.",
				Details: fmt.Sprintf("Event retrieval error: %v", err),
			})
			return
		}

		// Fetch event-related details like city, province
		city, err := h.cityService.GetCityByID(ctx, event.Event.CityID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "CITY_NOT_FOUND",
				Message: "Kota tidak ditemukan.",
				Details: fmt.Sprintf("City retrieval error: %v", err),
			})
			return
		}

		// Prepare full event context
		eventContext = AIInteractionContext{
			UserProfile:  userProfile,
			Event:        &event.Event,
			City:         &city.City,
			Conversation: session.Messages,
		}
	} else {
		eventContext = AIInteractionContext{
			UserProfile:  userProfile,
			Conversation: session.Messages,
		}
	}

	// Send message with comprehensive context
	response, err := h.aiService.SendMessageWithContext(sessionID, req.Message, eventContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "MESSAGE_SEND_FAILED",
			Message: "Gagal memproses pesan. Silakan coba lagi.",
			Details: fmt.Sprintf("Message processing error: %v", err),
		})
		return
	}

	// Split and clean response
	var responseData struct {
		Response []string `json:"response"`
	}

	// Remove asterisks, backslashes, and quotes
	cleanedResponse := strings.ReplaceAll(response, "*", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "\\", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "\"", "")
	responseData.Response = strings.Split(cleanedResponse, "\n")

	// Trim any empty lines from the start or end
	var trimmedLines []string
	for _, line := range responseData.Response {
		if strings.TrimSpace(line) != "" {
			trimmedLines = append(trimmedLines, line)
		}
	}
	responseData.Response = trimmedLines

	c.JSON(http.StatusOK, SendMessageResponse{
		Response: responseData.Response,
	})
}

// GenerateEventDescription creates an AI-generated event description
func (h *GeminiHandler) GenerateEventDescription(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("event_exploration"); err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FEATURE_NOT_SUPPORTED",
			Message: "Event exploration feature is currently restricted",
			Details: err.Error(),
		})
		return
	}

	ctx := context.Background()
	eventID := c.Param("eventID")

	// Fetch event details with more robust error handling
	event, err := h.eventService.GetEventByID(ctx, eventID)
	if err != nil {
		// Log the full error for debugging
		log.Printf("Event retrieval error: %v", err)

		// Check for specific Supabase/PostgreSQL row retrieval errors
		errorMessage := err.Error()
		switch {
		case strings.Contains(errorMessage, "PGRST116") ||
			strings.Contains(errorMessage, "multiple (or no) rows returned") ||
			strings.Contains(errorMessage, "no rows in result set"):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "EVENT_NOT_FOUND",
				Message: "Event tidak ditemukan atau memiliki data ganda. Silakan periksa kembali ID event.",
				Details: fmt.Sprintf("Row retrieval error: %v", errorMessage),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "EVENT_RETRIEVAL_ERROR",
				Message: "Terjadi kesalahan saat mengambil detail event.",
				Details: fmt.Sprintf("Unexpected error: %v", err),
			})
			return
		}
	}

	// Additional null checks
	if event == nil || event.Event.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "EVENT_NOT_FOUND",
			Message: "Event tidak ditemukan. Silakan periksa kembali ID event.",
			Details: "Nil or invalid event returned from service",
		})
		return
	}

	// Generate AI description
	description, err := h.generateAIEventDescription(&event.Event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "DESCRIPTION_GENERATION_FAILED",
			Message: "Gagal menghasilkan deskripsi event. Silakan coba lagi.",
			Details: fmt.Sprintf("Description generation error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, EventDescriptionResponse{
		Description: description,
	})
}

// generateAIEventDescription creates an AI-powered event description
func (h *GeminiHandler) generateAIEventDescription(event *culturalModels.Event) (string, error) {
	// Prepare context-aware prompt using system policies
	prompt := h.buildEventDescriptionPrompt(event)

	// Initialize Gemini client
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: h.config.GeminiAI.ApiKey,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}

	// Generate description with system policy
	result, err := client.Models.GenerateContent(context.Background(), h.config.GeminiAI.AIModel, genai.Text(prompt), &genai.GenerateContentConfig{
		Temperature: h.config.GeminiAI.Temperature,
		TopP:        h.config.GeminiAI.TopP,
		TopK:        h.config.GeminiAI.TopK,
		SystemInstruction: genai.NewContentFromParts([]*genai.Part{
			{
				Text: GetSystemPolicies(Feature, Response, Strictness),
			},
		}, "system"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate event description: %v", err)
	}

	// Validate and return response
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no description generated: empty response from AI")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

// buildEventDescriptionPrompt creates a comprehensive prompt for event description
func (h *GeminiHandler) buildEventDescriptionPrompt(event *culturalModels.Event) string {
	return fmt.Sprintf(`Generate a compelling and informative description for a cultural event in Indonesia:
- Event Name: %s
- Start Date: %s
- End Date: %s
- Initial Description: %s

Create a rich, engaging description that:
- Highlights the cultural significance
- Explains why tourists should attend
- Provides context about the event's history and importance
- Suggests potential activities or experiences
- Maintains an enthusiastic and inviting tone
- Strictly focuses on local Indonesian cultural context`,
		event.Name,
		event.StartDate.Format("2006-01-02"),
		event.EndDate.Format("2006-01-02"),
		event.Description)
}

// Request and Response Structures

// CreateChatSessionRequest represents the request payload for creating a chat session
// @Description Request to create a new AI chat session
type CreateChatSessionRequest struct {
	// UserID is the unique identifier of the user creating the session
	// @Required true
	UserID string `json:"user_id" binding:"required"`

	// Optional EventID to provide context for the chat session
	EventID *string `json:"event_id,omitempty"`
}

// CreateChatSessionResponse represents the response after creating a chat session
// @Description Response containing the created session ID
type CreateChatSessionResponse struct {
	// Unique identifier for the created chat session
	SessionID string `json:"session_id"`
}

// SendMessageRequest represents the request payload for sending a message in a chat session
// @Description Request to send a message to the AI
type SendMessageRequest struct {
	// Message content to be sent to the AI
	// @Required true
	// @Max length 500
	Message string `json:"message" binding:"required,max=500"`
}

// SendMessageResponse represents the AI's response to a message
// @Description Response from the AI containing multiple lines of text
type SendMessageResponse struct {
	// Multiple lines of the AI's response
	Response []string `json:"response"`
}

// EventDescriptionResponse represents the AI-generated event description
// @Description Response containing an AI-generated description for an event
type EventDescriptionResponse struct {
	// Comprehensive description of the event
	Description string `json:"description"`
}
