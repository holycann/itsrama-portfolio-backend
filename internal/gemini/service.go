package gemini

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"google.golang.org/genai"
)

// AICache provides a thread-safe caching mechanism for AI responses
type AICache struct {
	mu       sync.RWMutex
	cache    map[string]cachedResponse
	maxSize  int
	duration time.Duration
}

// cachedResponse represents a cached AI response with metadata
type cachedResponse struct {
	content   string
	timestamp time.Time
}

// NewAICache creates a new thread-safe AI response cache
func NewAICache(maxSize int, duration time.Duration) *AICache {
	return &AICache{
		cache:    make(map[string]cachedResponse),
		maxSize:  maxSize,
		duration: duration,
	}
}

// generateCacheKey creates a unique key for caching based on input parameters
func generateCacheKey(model string, temperature float32, topK int32, topP float32, text string) string {
	// Create a hash of all input parameters
	input := fmt.Sprintf("%s|%f|%d|%f|%s", model, temperature, topK, topP, text)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// get retrieves a cached response if available and not expired
func (c *AICache) get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return "", false
	}

	// Check if entry is expired
	if time.Since(entry.timestamp) > c.duration {
		return "", false
	}

	return entry.content, true
}

// set adds a response to the cache, managing cache size
func (c *AICache) set(key string, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove expired entries
	now := time.Now()
	for k, v := range c.cache {
		if now.Sub(v.timestamp) > c.duration {
			delete(c.cache, k)
		}
	}

	// Manage cache size
	if len(c.cache) >= c.maxSize {
		// Remove oldest entry
		var oldestKey string
		var oldestTime time.Time
		for k, v := range c.cache {
			if oldestTime.IsZero() || v.timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.timestamp
			}
		}
		delete(c.cache, oldestKey)
	}

	// Add new entry
	c.cache[key] = cachedResponse{
		content:   content,
		timestamp: now,
	}
}

// AIGenerationOptions provides flexible configuration for text generation
type AIGenerationOptions struct {
	// Optional override for model configuration
	Model       string
	Temperature *float32
	TopK        *int32
	TopP        *float32

	// Caching options
	DisableCache bool

	// Logging and tracing
	TraceID string
}

// GenerateTextWithOptions generates text with advanced configuration and caching
func (a *AIClient) GenerateTextWithOptions(ctx context.Context, text string, opts *AIGenerationOptions) (string, error) {
	// Use default options if not provided
	if opts == nil {
		opts = &AIGenerationOptions{}
	}

	// Determine model and configuration
	model := a.config.AIModel
	if opts.Model != "" {
		model = opts.Model
	}

	temperature := a.config.Temperature
	if opts.Temperature != nil {
		temperature = *opts.Temperature
	}

	topK := a.config.TopK
	if opts.TopK != nil {
		topK = *opts.TopK
	}

	topP := a.config.TopP
	if opts.TopP != nil {
		topP = *opts.TopP
	}

	// Generate cache key
	cacheKey := generateCacheKey(model, temperature, topK, topP, text)

	// Check cache if not disabled
	if !opts.DisableCache {
		if cachedResponse, found := a.cache.get(cacheKey); found {
			a.logger.Info("AI response retrieved from cache",
				"trace_id", opts.TraceID,
				"cache_key", cacheKey)
			return cachedResponse, nil
		}
	}

	// Prepare the model configuration
	topKFloat := float32(topK)
	modelCfg := &genai.GenerateContentConfig{
		Temperature: &temperature,
		TopK:        &topKFloat,
		TopP:        &topP,
	}

	// Prepare content
	content := []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			genai.NewPartFromText(text),
		},
	}}

	// Log generation attempt
	startTime := time.Now()
	a.logger.Info("Generating AI text",
		"trace_id", opts.TraceID,
		"model", model,
		"prompt_length", len(text))

	// Attempt to generate content
	var result *genai.GenerateContentResponse
	var err error

	result, err = a.client.Models.GenerateContent(ctx, model, content, modelCfg)

	// Log generation result
	duration := time.Since(startTime)
	if err != nil {
		a.logger.Error("AI text generation failed",
			"trace_id", opts.TraceID,
			"error", err,
			"duration", duration)
		return "", fmt.Errorf("content generation failed: %w", err)
	}

	// Validate result
	if result == nil || len(result.Candidates) == 0 {
		a.logger.Warn("No AI candidates generated",
			"trace_id", opts.TraceID,
			"duration", duration)
		return "", fmt.Errorf("no candidates generated")
	}

	// Extract text from the first candidate
	candidate := result.Candidates[0]
	if candidate == nil || len(candidate.Content.Parts) == 0 {
		a.logger.Warn("No text in AI candidate",
			"trace_id", opts.TraceID,
			"duration", duration)
		return "", fmt.Errorf("no text in candidate")
	}

	generatedText := candidate.Content.Parts[0].Text

	// Cache the response if not disabled
	if !opts.DisableCache {
		a.cache.set(cacheKey, generatedText)
	}

	// Log successful generation
	a.logger.Info("AI text generated successfully",
		"trace_id", opts.TraceID,
		"response_length", len(generatedText),
		"duration", duration)

	return generatedText, nil
}

// Wrapper for backward compatibility
func (a *AIClient) GenerateText(ctx context.Context, text string) (string, error) {
	return a.GenerateTextWithOptions(ctx, text, nil)
}
