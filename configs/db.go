package configs

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
