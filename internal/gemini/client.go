package gemini

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/holycann/cultour-backend/internal/supabase"
	"google.golang.org/genai"
)

// SafetySettings defines content safety configurations
type SafetySettings struct {
	HarmCategory   string
	HarmBlockLevel int // Use int to represent safety threshold
}

// CacheConfig defines caching strategy for AI responses
type CacheConfig struct {
	Enabled    bool
	MaxSize    int
	Expiration time.Duration
}

// LoggingConfig defines detailed logging for AI interactions
type LoggingConfig struct {
	Enabled       bool
	LogPrompts    bool
	LogResponses  bool
	LogTokenUsage bool
	LogLatency    bool
}

// Config represents comprehensive configuration for Gemini AI client
type Config struct {
	// Core AI Configuration
	ApiKey          string
	AIModel         string
	Temperature     float32
	TopK            int32
	TopP            float32
	MaxOutputTokens int32

	// Advanced Safety and Filtering
	SafetySettings []SafetySettings

	// Caching Configuration
	CacheConfig CacheConfig

	// Logging Configuration
	LoggingConfig LoggingConfig

	// System Instruction for consistent AI behavior
	SystemInstruction string

	// Optional dependencies
	Logger         *slog.Logger
	SupabaseClient *supabase.SupabaseClient
}

// Safety threshold constants
const (
	SafetyThresholdBlockNone        = 0
	SafetyThresholdBlockLow         = 1
	SafetyThresholdBlockMedium      = 2
	SafetyThresholdBlockMediumAndUp = 3
	SafetyThresholdBlockHigh        = 4
)

// DefaultConfig provides a standard configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		AIModel:         "gemini-2.0-flash",
		Temperature:     0.7,
		TopK:            40,
		TopP:            0.95,
		MaxOutputTokens: 2048,
		SafetySettings: []SafetySettings{
			{
				HarmCategory:   "HARASSMENT",
				HarmBlockLevel: SafetyThresholdBlockMediumAndUp,
			},
			{
				HarmCategory:   "HATE_SPEECH",
				HarmBlockLevel: SafetyThresholdBlockMediumAndUp,
			},
		},
		CacheConfig: CacheConfig{
			Enabled:    true,
			MaxSize:    100,
			Expiration: 1 * time.Hour,
		},
		LoggingConfig: LoggingConfig{
			Enabled:       true,
			LogPrompts:    true,
			LogResponses:  true,
			LogTokenUsage: true,
			LogLatency:    true,
		},
		SystemInstruction: GetFullSystemPolicy(),
	}
}

// NewGeminiAIClient creates a more robust Gemini AI client with advanced configuration
func NewGeminiAIClient(cfg *Config) (*AIClient, error) {
	// Validate configuration
	if cfg == nil {
		cfg = DefaultConfig()
	}
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Create Google AI client with enhanced configuration
	clientCfg := &genai.ClientConfig{
		APIKey: cfg.ApiKey,
	}
	client, err := genai.NewClient(context.Background(), clientCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini AI client: %w", err)
	}

	// Create cache with default settings
	cache := NewAICache(
		cfg.CacheConfig.MaxSize,
		cfg.CacheConfig.Expiration,
	)

	return &AIClient{
		client:         client,
		config:         cfg,
		logger:         cfg.Logger,
		supabaseClient: cfg.SupabaseClient,
		cache:          cache,
	}, nil
}

// AIClient now includes more configuration and advanced features
type AIClient struct {
	client         *genai.Client
	config         *Config
	messageRepo    MessageRepo
	logger         *slog.Logger
	supabaseClient *supabase.SupabaseClient
	cache          *AICache
}
