package configs

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Supabase    SupabaseConfig
	Database    DatabaseConfig
	GeminiAI    GeminiAIConfig
	Logging     LoggingConfig
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Environment: getEnv("APP_ENV", "development"),
		Server:      loadServerConfig(),
		Supabase:    loadSupabaseConfig(),
		Database:    loadDatabaseConfig(),
		GeminiAI:    loadGeminiAIConfig(),
		Logging:     loadLoggingConfig(),
	}

	createDirIfNotExists(config.Logging.FilePath)

	return config, nil
}
