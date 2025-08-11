package gemini

import (
	"context"
	"time"

	"github.com/holycann/cultour-backend/configs"
	achievementServices "github.com/holycann/cultour-backend/internal/achievement/services"
	culturalModels "github.com/holycann/cultour-backend/internal/cultural/models"
	culturalServices "github.com/holycann/cultour-backend/internal/cultural/services"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// AIGenerationInterface defines methods for AI content generation
type AIGenerationInterface interface {
	// Generate content based on a query
	GenerateContent(ctx context.Context, query string, params map[string]interface{}) (string, error)

	// Generate specific content types
	GenerateEventDescription(ctx context.Context, userID string, params map[string]interface{}) (string, error)
}

// AISessionInterface defines methods for managing chat sessions
type AISessionInterface interface {
	// Session management
	CreateChatSession(payload CreateChatSessionRequest) (string, error)
	GetChatSession(sessionID string) (*ChatSession, error)
	SendMessage(payload SendMessageRequest) (string, error)
	SendMessageWithContext(sessionID string, message string, context map[string]interface{}) (string, error)
}

// AIKnowledgeInterface defines methods for knowledge base operations
type AIKnowledgeInterface interface {
	// Knowledge base initialization
	PopulateKnowledgeBase(ctx context.Context) error

	// Domain-specific knowledge operations
	LoadUserKnowledge(ctx context.Context, userID string) error
	LoadEventKnowledge(ctx context.Context, eventID string) error
	LoadPlaceKnowledge(ctx context.Context, placeID string) error

	// Context building
	BuildContextAwarePrompt(ctx context.Context, params map[string]interface{}, query string) (string, error)
}

// AIServiceInterface combines all AI service functionalities
type AIServiceInterface interface {
	AIGenerationInterface
	AISessionInterface
	AIKnowledgeInterface

	// Resource management
	Close()
}

// AIService implements AIServiceInterface
type AIService struct {
	client         *CultourAIClient
	config         *configs.Config
	logger         *logger.Logger
	adapterManager *DomainAdapterManager

	// Domain services
	eventService       culturalServices.EventService
	userService        userServices.UserService
	cityService        placeServices.CityService
	locationService    placeServices.LocationService
	provinceService    placeServices.ProvinceService
	badgeService       achievementServices.BadgeService
	userProfileService userServices.UserProfileService
	userBadgeService   userServices.UserBadgeService
}

// NewAIService creates a new AI service with all dependencies
func NewAIService(
	config *configs.Config,
	logger *logger.Logger,
	eventService culturalServices.EventService,
	userService userServices.UserService,
	cityService placeServices.CityService,
	locationService placeServices.LocationService,
	provinceService placeServices.ProvinceService,
	badgeService achievementServices.BadgeService,
	userProfileService userServices.UserProfileService,
	userBadgeService userServices.UserBadgeService,
) (*AIService, error) {
	// Create AI client
	client, err := NewCultourAIClient(config, logger)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to create AI client")
	}

	service := &AIService{
		client:             client,
		config:             config,
		logger:             logger,
		eventService:       eventService,
		userService:        userService,
		cityService:        cityService,
		locationService:    locationService,
		provinceService:    provinceService,
		badgeService:       badgeService,
		userProfileService: userProfileService,
		userBadgeService:   userBadgeService,
	}

	// Create domain adapter manager
	adapterManager := NewDomainAdapterManager(client.knowledgeBase, logger)

	// Create and register domain adapters
	culturalAdapter := NewCulturalDomainAdapter(eventService, logger)
	placeAdapter := NewPlaceDomainAdapter(locationService, cityService, provinceService, logger)
	userAdapter := NewUserDomainAdapter(
		userService,
		userProfileService,
		userBadgeService,
		badgeService,
		logger,
	)

	// Register adapters
	_ = adapterManager.RegisterAdapter("cultural", culturalAdapter)
	_ = adapterManager.RegisterAdapter("place", placeAdapter)
	_ = adapterManager.RegisterAdapter("user", userAdapter)

	service.adapterManager = adapterManager

	return service, nil
}

// PopulateKnowledgeBase initializes the knowledge base with essential data
func (s *AIService) PopulateKnowledgeBase(ctx context.Context) error {
	if s.adapterManager == nil {
		return errors.New(errors.ErrInternal, "adapter manager not initialized", nil)
	}

	// Use adapter manager to load data from all domains
	return s.adapterManager.LoadAllData(ctx)
}

// LoadUserKnowledge loads specific user knowledge into the knowledge base
func (s *AIService) LoadUserKnowledge(ctx context.Context, userID string) error {
	// Find user adapter
	userAdapter, ok := s.adapterManager.adapters["user"]
	if !ok {
		return errors.New(errors.ErrInternal, "user adapter not found", nil)
	}

	// Try to cast to UserDomainAdapter to access LoadUserData
	if userDomainAdapter, ok := userAdapter.(*UserDomainAdapter); ok {
		return userDomainAdapter.LoadUserData(ctx, userID)
	}

	return errors.New(errors.ErrInternal, "failed to cast to user domain adapter", nil)
}

// LoadEventKnowledge loads specific event knowledge into the knowledge base
func (s *AIService) LoadEventKnowledge(ctx context.Context, eventID string) error {
	event, err := s.eventService.GetEventByID(ctx, eventID)
	if err != nil {
		return errors.Wrap(err, errors.ErrNotFound, "event not found")
	}

	// Create Event model and add to knowledge base
	eventModel := culturalModels.Event{
		ID:            event.ID,
		Name:          event.Name,
		Description:   event.Description,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		IsKidFriendly: event.IsKidFriendly,
	}

	// Add to knowledge base
	s.client.knowledgeBase.AddEvent(&eventModel)
	return nil
}

// LoadPlaceKnowledge loads specific place knowledge into the knowledge base
func (s *AIService) LoadPlaceKnowledge(ctx context.Context, placeID string) error {
	// Find place adapter
	placeAdapter, ok := s.adapterManager.adapters["place"]
	if !ok {
		return errors.New(errors.ErrInternal, "place adapter not found", nil)
	}

	// Add location context to params
	params := map[string]interface{}{
		"location_id": placeID,
	}

	// Build context (this will load necessary data)
	_, err := placeAdapter.BuildContext(ctx, params)
	return err
}

// BuildContextAwarePrompt creates a prompt with rich context for better AI understanding
func (s *AIService) BuildContextAwarePrompt(ctx context.Context, params map[string]interface{}, query string) (string, error) {
	if s.adapterManager == nil {
		return query, nil // Return query as is if no adapter manager
	}

	// Build comprehensive context using all domain adapters
	context, err := s.adapterManager.BuildComprehensiveContext(ctx, params)
	if err != nil {
		s.logger.Warn("Error building comprehensive context", map[string]interface{}{
			"error": err.Error(),
		})
		// Continue with just the query if context building fails
	}

	if context != "" {
		// Combine context and query
		return context + "\n\n" + query, nil
	}

	return query, nil
}

// CreateChatSession creates a new chat session
func (s *AIService) CreateChatSession(payload CreateChatSessionRequest) (string, error) {
	session, err := s.client.CreateSession(payload.UserID, payload.EventID)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to create chat session")
	}

	// Preload user and event knowledge if available
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("Creating chat session", map[string]interface{}{
		"user_id":  payload.UserID,
		"event_id": payload.EventID,
	})

	// Load user knowledge
	_ = s.LoadUserKnowledge(ctx, payload.UserID)

	// Load event knowledge if provided
	if payload.EventID != nil && *payload.EventID != "" {
		_ = s.LoadEventKnowledge(ctx, *payload.EventID)
	}

	return session.ID, nil
}

// GetChatSession retrieves an existing chat session
func (s *AIService) GetChatSession(sessionID string) (*ChatSession, error) {
	return s.client.GetSession(sessionID)
}

// SendMessage sends a message in an existing chat session
func (s *AIService) SendMessage(payload SendMessageRequest) (string, error) {
	// Add message to session
	if err := s.client.AddMessage(payload.SessionID, "user", payload.Message); err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to add message to session")
	}

	// Generate response
	return s.client.GenerateResponse(payload.SessionID, payload.Message)
}

// SendMessageWithContext sends a message with additional context
func (s *AIService) SendMessageWithContext(sessionID string, message string, params map[string]interface{}) (string, error) {
	ctx := context.Background()

	// Build context-aware prompt
	contextualPrompt, err := s.BuildContextAwarePrompt(ctx, params, message)
	if err != nil {
		s.logger.Warn("Error building contextual prompt", map[string]interface{}{
			"error": err.Error(),
		})
		// Continue with original message if context building fails
		contextualPrompt = message
	}

	// Add original message to session (not the enriched prompt)
	if err := s.client.AddMessage(sessionID, "user", message); err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to add message to session")
	}

	// Generate response with contextual prompt
	return s.client.GenerateResponse(sessionID, contextualPrompt)
}

// GenerateContent generates AI content based on a query and parameters
func (s *AIService) GenerateContent(ctx context.Context, query string, params map[string]interface{}) (string, error) {
	// Create a temporary session for this query
	userID, ok := params["user_id"].(string)
	if !ok {
		return "", errors.New(errors.ErrBadRequest, "user ID is required", nil)
	}

	eventID, ok := params["event_id"].(string)
	if !ok {
		eventID = ""
	}

	if userID == "" {
		return "", errors.New(errors.ErrBadRequest, "user ID is required", nil)
	}

	sessionID, err := s.CreateChatSession(CreateChatSessionRequest{
		UserID:  userID,
		EventID: &eventID,
	})
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to create temporary session")
	}

	// Send message with context
	return s.SendMessageWithContext(sessionID, query, params)
}

// GenerateEventDescription generates a description for an event
func (s *AIService) GenerateEventDescription(ctx context.Context, userID string, params map[string]interface{}) (string, error) {
	// Build prompt for event description generation using title and optional additional context
	prompt := s.buildEventDescriptionPrompt(params)
	if prompt == "" {
		return "", errors.New(errors.ErrBadRequest, "title is required to generate event description", nil)
	}

	// Generate content
	return s.GenerateContent(ctx, prompt, params)
}

// buildEventDescriptionPrompt creates a prompt for generating an event description
func (s *AIService) buildEventDescriptionPrompt(params map[string]interface{}) string {
	title, _ := params["title"].(string)
	additional, _ := params["additional_context"].(string)

	if title == "" {
		return ""
	}

	prompt := "Please generate a comprehensive, culturally rich, and engaging description for the following proposed cultural event based on its title.\n" +
		"Event Title: " + title + "\n"
	if additional != "" {
		prompt += "Additional Context: " + additional + "\n"
	}
	prompt += "\nFocus on cultural significance, unique highlights, audience appeal, and why it's special."
	prompt += "\nMax 1000 characters with Indonesian language."
	return prompt
}

// Close releases resources used by the service
func (s *AIService) Close() {
	// Any cleanup needed
	s.logger.Info("Closing AI service")
}
