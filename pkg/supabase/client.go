package supabase

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
)

type SupabaseClientConfig struct {
	ProjectID string
	ApiSecret string
	Schema    string
}

type SupabaseClient struct {
	client *supabase.Client
}

func NewSupabaseClient(cfg SupabaseClientConfig) (*SupabaseClient, error) {
	if cfg.ApiSecret == "" || cfg.ProjectID == "" {
		return nil, fmt.Errorf("supabase API key & project ID cannot be empty")
	}

	client, err := supabase.NewClient(fmt.Sprintf("https://%s.supabase.co", cfg.ProjectID), cfg.ApiSecret, &supabase.ClientOptions{
		Schema: cfg.Schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Supabase client: %v", err)
	}

	return &SupabaseClient{
		client: client,
	}, nil
}

func (s *SupabaseClient) GetClient() *supabase.Client {
	return s.client
}
