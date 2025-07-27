package configs

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	PoolMode     string
}

type SupabaseConfig struct {
	ApiPublicKey    string
	ApiSecretKey    string
	JwtApiKeySecret string
	ProjectID       string
	StorageBucketID string
}

func loadSupabaseConfig() SupabaseConfig {
	return SupabaseConfig{
		ApiPublicKey:    getEnv("SUPABASE_API_PUBLIC_KEY", ""),
		ApiSecretKey:    getEnv("SUPABASE_API_SECRET_KEY", ""),
		JwtApiKeySecret: getEnv("SUPABASE_JWT_API_SECRET_KEY", ""),
		ProjectID:       getEnv("SUPABASE_PROJECT_ID", ""),
		StorageBucketID: getEnv("SUPABASE_STORAGE_BUCKET_ID", ""),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvAsInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", ""),
		Password:     getEnv("DB_PASSWORD", ""),
		DatabaseName: getEnv("DB_NAME", ""),
		PoolMode:     getEnv("DB_POOL_MODE", "transaction"),
	}
}
