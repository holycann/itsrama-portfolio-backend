package gemini

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/holycann/cultour-backend/internal/supabase"
	"google.golang.org/genai"
)

// GeminiAIConfig represents the comprehensive configuration options for the Gemini AI client
// It allows fine-tuning of AI generation parameters, model selection, and safety controls
type Config struct {
	// ApiKey is the authentication key required to access Google AI services
	ApiKey string

	// AIModel specifies the specific Gemini AI model to use, such as "gemini-pro" or "gemini-ultra"
	AIModel string

	Tuning *genai.GenerateContentConfig

	Logger         *slog.Logger
	SupabaseClient supabase.SupabaseClient
}

// geminiAiClient implements the GeminiAIClient interface
type AIClient struct {
	client         *genai.Client
	messageRepo    MessageRepo
	model          string
	tuning         *genai.GenerateContentConfig
	logger         *slog.Logger
	supabaseClient *supabase.SupabaseClient
}

// NewGeminiAiClient creates a new Gemini AI client
func NewGeminiAIClient(cfg *Config) (*AIClient, error) {
	// Validate configuration
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Create Google AI client
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: cfg.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini AI client: %w", err)
	}

	return &AIClient{
		client:         client,
		model:          cfg.AIModel,
		tuning:         cfg.Tuning,
		logger:         cfg.Logger,
		supabaseClient: &cfg.SupabaseClient,
	}, nil
}
