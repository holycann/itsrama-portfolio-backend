package gemini

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/configs"
	achievementModels "github.com/holycann/cultour-backend/internal/achievement/models"
	culturalModels "github.com/holycann/cultour-backend/internal/cultural/models"
	discussionModels "github.com/holycann/cultour-backend/internal/discussion/models"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
	"google.golang.org/genai"
)

// AIService provides specialized AI interactions for Cultour
type AIService struct {
	client *CultourAIClient
	config *configs.Config
}

// AIInteractionContext provides context for AI interactions
type AIInteractionContext struct {
	UserProfile  *userModels.UserProfile
	UserBadge    *userModels.UserBadge
	User         *userModels.User
	Event        *culturalModels.Event
	LocalStory   *culturalModels.LocalStory
	City         *placeModels.City
	Province     *placeModels.Province
	Location     *placeModels.Location
	Thread       *discussionModels.Thread
	Badges       []*achievementModels.Badge
	Conversation []ChatMessage
}

// NewAIService creates a new AI service for Cultour
func NewAIService(config *configs.Config) (*AIService, error) {
	aiClient, err := NewCultourAIClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI service: %v", err)
	}

	return &AIService{
		client: aiClient,
		config: config,
	}, nil
}

// PopulateKnowledgeBase adds comprehensive data to the AI's knowledge base
func (s *AIService) PopulateKnowledgeBase(context AIInteractionContext) {
	kb := s.client.knowledgeBase

	// Add user information
	if context.User != nil {
		kb.AddUser(context.User)
	}

	// Add user profile
	if context.UserProfile != nil {
		kb.AddUserProfile(context.UserProfile)
	}

	// Add user badge
	if context.UserBadge != nil {
		kb.AddUserBadge(context.User.ID, context.UserBadge)
	}

	// Add event details
	if context.Event != nil {
		kb.AddEvent(context.Event)
	}

	// Add local story
	if context.LocalStory != nil {
		kb.AddLocalStory(context.LocalStory)
	}

	// Add location details
	if context.City != nil {
		kb.AddCity(context.City)
	}

	if context.Province != nil {
		kb.AddProvince(context.Province)
	}

	if context.Location != nil {
		kb.AddLocation(context.Location)
	}

	// Add discussion context
	if context.Thread != nil {
		kb.AddThread(context.Thread)
	}

	// Add user achievements
	if context.Badges != nil {
		for _, badge := range context.Badges {
			kb.AddBadge(badge)
		}
	}

	// Add contextual facts
	kb.AddContextualFact("app_description", "Cultour is an innovative cultural tourism app focused on Indonesian experiences")
}

// CreateChatSession initializes a new chat session
func (s *AIService) CreateChatSession(userID string, eventID *string) (string, error) {
	session, err := s.client.CreateSession(userID, eventID)
	if err != nil {
		return "", err
	}
	return session.ID, nil
}

// SendMessage processes a user message and generates AI response
func (s *AIService) SendMessage(sessionID, message string) (string, error) {
	// Add user message to session
	err := s.client.AddMessage(sessionID, "user", message)
	if err != nil {
		return "", err
	}

	// Generate AI response
	response, err := s.client.GenerateResponse(sessionID, message)
	if err != nil {
		return "", err
	}

	return response, nil
}

// GenerateEventDescription creates an enhanced AI-generated event description
func (s *AIService) GenerateEventDescription(event *culturalModels.Event) (string, error) {
	// Prepare context-aware prompt
	prompt := fmt.Sprintf(`Generate a compelling and informative description for a cultural event:
- Event Name: %s
- Start Date: %s
- End Date: %s
- Initial Description: %s

Create a rich, engaging description that:
- Highlights the cultural significance
- Explains why tourists should attend
- Provides context about the event's history and importance
- Suggests potential activities or experiences
- Maintains an enthusiastic and inviting tone`,
		event.Name,
		event.StartDate.Format("2006-01-02"),
		event.EndDate.Format("2006-01-02"),
		event.Description)

	// Use Gemini to generate response
	resp, err := s.client.client.Models.GenerateContent(context.Background(), s.config.GeminiAI.AIModel, genai.Text(prompt), &genai.GenerateContentConfig{
		Temperature: s.config.GeminiAI.Temperature,
		TopP:        s.config.GeminiAI.TopP,
		TopK:        s.config.GeminiAI.TopK,
		SystemInstruction: genai.NewContentFromParts([]*genai.Part{
			{
				Text: GetFullSystemPolicy(),
			},
		}, "system"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate event description: %v", err)
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// GetChatSession retrieves a chat session by its ID
func (s *AIService) GetChatSession(sessionID string) (*ChatSession, error) {
	// Retrieve session from the client's session manager
	session, err := s.client.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chat session: %v", err)
	}
	return session, nil
}

// SendMessageWithContext sends a message with comprehensive user and event context
func (s *AIService) SendMessageWithContext(sessionID string, message string, context AIInteractionContext) (string, error) {
	// Populate knowledge base with context
	s.PopulateKnowledgeBase(context)

	// Prepare contextual prompt
	contextualPrompt := s.BuildContextAwarePrompt(context, message)

	// Send message
	response, err := s.client.GenerateResponse(sessionID, contextualPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate AI response: %v", err)
	}

	return response, nil
}

// BuildContextAwarePrompt creates a comprehensive context-specific AI interaction
func (s *AIService) BuildContextAwarePrompt(context AIInteractionContext, query string) string {
	// Prepare base context
	baseContext := fmt.Sprintf(
		"User Profile: %s\n",
		context.UserProfile.Fullname,
	)

	// Add event context if available
	if context.Event != nil {
		baseContext += fmt.Sprintf(
			"Event Context:\n"+
				"Event Name: %s\n"+
				"Event Description: %s\n"+
				"Event Location: %s\n",
			context.Event.Name,
			context.Event.Description,
			context.City.Name,
		)
	}

	// Combine base context with user query
	return fmt.Sprintf(
		"%s\n\n"+
			"User Query: %s\n\n"+
			"Please provide a helpful and contextually relevant response.",
		baseContext,
		query,
	)
}

// AIServiceInterface defines the contract for AI interactions
type AIServiceInterface interface {
	CreateChatSession(userString string, eventID *int) (string, error)
	SendMessage(sessionID, message string) (string, error)
	GenerateEventDescription(event *culturalModels.Event) (string, error)
	BuildContextAwarePrompt(context AIInteractionContext, query string) string
	Close()
}
