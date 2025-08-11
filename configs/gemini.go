package configs

type GeminiAIConfig struct {
	ApiKey      string
	AIModel     string
	Temperature *float32
	TopK        *float32
	TopP        *float32
	MaxTokens   int
}

func loadGeminiAIConfig() GeminiAIConfig {
	return GeminiAIConfig{
		ApiKey:  getEnv("GEMINI_API_KEY", ""),
		AIModel: getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
		Temperature: func() *float32 {
			val := getEnvAsFloat32("GEMINI_TEMPERATURE", 0.7)
			return &val
		}(),
		TopK: func() *float32 {
			val := getEnvAsFloat32("GEMINI_TOP_K", 40)
			return &val
		}(),
		TopP: func() *float32 {
			val := getEnvAsFloat32("GEMINI_TOP_P", 0.85)
			return &val
		}(),
		MaxTokens: getEnvAsInt("GEMINI_MAX_TOKEN", 256),
	}
}
