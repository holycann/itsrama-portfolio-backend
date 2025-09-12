package configs

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	PoolMode     string
	Schema       string
}

type SupabaseConfig struct {
	ApiPublicKey         string
	ApiSecretKey         string
	ProjectID            string
	JWTSecret            string
	StorageBucketID      string
	DefaultStorageFolder string
	MaxFileSize          int64
	AllowedFileTypes     []string
	CacheControl         string
}

func loadSupabaseConfig() SupabaseConfig {
	return SupabaseConfig{
		ApiPublicKey:         getEnv("SUPABASE_API_PUBLIC_KEY", ""),
		ApiSecretKey:         getEnv("SUPABASE_API_SECRET_KEY", ""),
		ProjectID:            getEnv("SUPABASE_PROJECT_ID", ""),
		JWTSecret:            getEnv("SUPABASE_JWT_API_SECRET_KEY", ""),
		StorageBucketID:      getEnv("SUPABASE_STORAGE_BUCKET_ID", ""),
		DefaultStorageFolder: getEnv("SUPABASE_DEFAULT_STORAGE_FOLDER", "documents"),
		MaxFileSize:          int64(getEnvAsInt("SUPABASE_MAX_FILE_SIZE", 10*1024*1024)), // 10MB default
		AllowedFileTypes:     getEnvAsStringSlice("SUPABASE_ALLOWED_FILE_TYPES", []string{"image/jpeg", "image/png", "application/pdf"}),
		CacheControl:         getEnv("SUPABASE_CACHE_CONTROL", "public, max-age=3600, must-revalidate"),
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
		Schema:       getEnv("DB_SCHEMA", "itsrama"),
	}
}
