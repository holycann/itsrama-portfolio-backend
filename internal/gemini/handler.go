package gemini

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/configs"
	achievementServices "github.com/holycann/cultour-backend/internal/achievement/services"
	culturalServices "github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/middleware"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	_ "github.com/holycann/cultour-backend/pkg/response"
)

// GeminiHandler manages AI-powered interactions and provides endpoints for AI chat sessions and event description generation
// @Tag AI Handler
type GeminiHandler struct {
	base.BaseHandler
	aiService          AIServiceInterface
	eventService       culturalServices.EventService
	cityService        placeServices.CityService
	userService        userServices.UserService
	provinceService    placeServices.ProvinceService
	locationService    placeServices.LocationService
	badgeService       achievementServices.BadgeService
	userProfileService userServices.UserProfileService
	userBadgeService   userServices.UserBadgeService
	config             *configs.Config
	logger             *logger.Logger
}

func NewGeminiHandler(
	config *configs.Config,
	logger *logger.Logger,
	eventService culturalServices.EventService,
	cityService placeServices.CityService,
	provinceService placeServices.ProvinceService,
	locationService placeServices.LocationService,
	userService userServices.UserService,
	badgeService achievementServices.BadgeService,
	userProfileService userServices.UserProfileService,
	userBadgeService userServices.UserBadgeService,
) (*GeminiHandler, error) {
	// Initialize AI service
	aiService, err := NewAIService(
		config,
		logger,
		eventService,
		userService,
		cityService,
		locationService,
		provinceService,
		badgeService,
		userProfileService,
		userBadgeService,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to initialize AI service")
	}

	baseHandler := base.NewBaseHandler(logger)

	return &GeminiHandler{
		BaseHandler:        *baseHandler,
		aiService:          aiService,
		eventService:       eventService,
		cityService:        cityService,
		userService:        userService,
		provinceService:    provinceService,
		locationService:    locationService,
		badgeService:       badgeService,
		userProfileService: userProfileService,
		userBadgeService:   userBadgeService,
		config:             config,
		logger:             logger,
	}, nil
}

// validateFeatureScope checks if the requested feature is currently supported in the application
func (h *GeminiHandler) validateFeatureScope(feature string) error {
	supportedFeatures := map[string]bool{
		"event_exploration": true,
		"ai_assistant":      true,
		"discussion_forum":  true,
		"warlok_creation":   true,
	}

	if !supportedFeatures[feature] {
		return errors.New(errors.ErrUnauthorized, "feature not currently supported", nil)
	}
	return nil
}

// CreateChatSession handles chat session creation with strict feature validation and user authentication
// @Summary Create a new AI chat session
// @Description Creates a new AI chat session for a user, with optional event context for personalized interactions
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body CreateChatSessionRequest true "Chat Session Creation Request with User ID and Optional Event Context"
// @Success 200 {object} CreateChatSessionResponse "Successfully created chat session with unique session identifier"
// @Failure 400 {object} response.APIResponse "Invalid request format or missing user ID"
// @Failure 403 {object} response.APIResponse "AI assistant feature is currently restricted"
// @Failure 404 {object} response.APIResponse "User not found in the system"
// @Failure 500 {object} response.APIResponse "Internal server error during session creation"
// @Router /ai/chat/session [post]
func (h *GeminiHandler) CreateChatSession(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("ai_assistant"); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrUnauthorized, "AI assistant feature is currently restricted"))
		return
	}

	// Get user ID from middleware context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrAuthentication, "Pengguna tidak terautentikasi"))
		return
	}

	var req CreateChatSessionRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrBadRequest, "Invalid request format. Please provide a valid event ID."))
		return
	}

	// Create chat session with event context if provided
	sessionID, err := h.aiService.CreateChatSession(userID, req.EventID)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrInternal, "Gagal membuat sesi. Silakan coba lagi."))
		return
	}

	h.HandleSuccess(c, CreateChatSessionResponse{
		SessionID: sessionID,
	}, "Chat session created successfully")
}

// SendMessage processes user messages with flexible topic handling and context preservation
// @Summary Send a message in an AI chat session
// @Description Sends a user message to the AI, retrieves the AI's response, and maintains conversation context
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param sessionID path string true "Unique Chat Session Identifier"
// @Param request body SendMessageRequest true "User Message Request with Content Validation"
// @Success 200 {object} SendMessageResponse "Successfully processed message with multi-line AI response"
// @Failure 400 {object} response.APIResponse "Invalid request format or message exceeds length limit"
// @Failure 403 {object} response.APIResponse "AI assistant feature is currently restricted"
// @Failure 404 {object} response.APIResponse "Chat session not found"
// @Failure 500 {object} response.APIResponse "Internal error processing message"
// @Router /ai/chat/{sessionID}/message [post]
func (h *GeminiHandler) SendMessage(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("ai_assistant"); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrUnauthorized, "AI assistant feature is currently restricted"))
		return
	}

	// Extract session ID from path
	sessionID := c.Param("sessionID")
	if sessionID == "" {
		h.HandleError(c, errors.New(errors.ErrBadRequest, "Session ID is required", nil))
		return
	}

	// Get user ID from middleware context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrAuthentication, "Pengguna tidak terautentikasi"))
		return
	}

	// Parse request body
	var req SendMessageRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrBadRequest, "Format pesan tidak valid. Mohon periksa kembali."))
		return
	}

	// Create context parameters
	params := map[string]interface{}{
		"user_id": userID,
	}

	// Get session info to retrieve context
	session, err := h.aiService.GetChatSession(sessionID)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrNotFound, "Sesi chat tidak ditemukan. Silakan buat sesi baru."))
		return
	}

	// Add event ID to context params if available
	if session.EventID != nil && *session.EventID != "" {
		params["event_id"] = *session.EventID
	}

	// Send message with context parameters
	responseText, err := h.aiService.SendMessageWithContext(sessionID, req.Message, params)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrInternal, "Gagal memproses pesan. Silakan coba lagi."))
		return
	}

	// Clean and format response
	cleanedResponse := strings.ReplaceAll(responseText, "*", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "\\", "")
	cleanedResponse = strings.ReplaceAll(cleanedResponse, "\"", "")

	// Split into lines and remove empty lines
	responseLines := strings.Split(cleanedResponse, "\n")
	var trimmedLines []string
	for _, line := range responseLines {
		if strings.TrimSpace(line) != "" {
			trimmedLines = append(trimmedLines, line)
		}
	}

	h.HandleSuccess(c, SendMessageResponse{
		Response: trimmedLines,
	}, "Message processed successfully")
}

// GenerateEventDescription creates an AI-generated event description for a new event
// @Summary Generate an AI-powered event description for a new event
// @Description Generates a rich, contextual description based on provided title and optional additional context.
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body GenerateEventDescriptionRequest true "Event Title and Optional Details"
// @Success 200 {object} EventDescriptionResponse "Successfully generated comprehensive event description"
// @Failure 400 {object} response.APIResponse "Invalid or missing event title"
// @Failure 403 {object} response.APIResponse "Event exploration feature is currently restricted"
// @Failure 500 {object} response.APIResponse "Internal error generating event description"
// @Router /ai/events/description [post]
func (h *GeminiHandler) GenerateEventDescription(c *gin.Context) {
	// Validate feature scope
	if err := h.validateFeatureScope("event_exploration"); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrUnauthorized, "Event exploration feature is currently restricted"))
		return
	}

	// Get user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrAuthentication, "Failed to retrieve user context"))
		return
	}

	// Parse request body
	var req GenerateEventDescriptionRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrBadRequest, "Invalid event title. Please provide a valid title."))
		return
	}

	ctx := context.Background()

	// Prepare context parameters (no event ID; use provided title and additional context)
	params := map[string]interface{}{
		"user_id":            userID,
		"title":              req.Title,
		"additional_context": req.AdditionalContext,
	}

	// Generate AI description using title/context
	description, err := h.aiService.GenerateEventDescription(ctx, userID, params)
	if err != nil {
		h.HandleError(c, errors.Wrap(err, errors.ErrInternal, "Failed to generate event description. Please try again."))
		return
	}

	h.HandleSuccess(c, EventDescriptionResponse{
		Description: description,
	}, "Event description generated successfully")
}
