package supabase

import (
	"fmt"
	"sync"

	"github.com/supabase-community/supabase-go"
)

type SupabaseClientConfig struct {
	ApiUrl    string
	ApiSecret string
}

type SupabaseClient struct {
	mu     sync.RWMutex
	client *supabase.Client
}

func NewSupabaseClient(cfg SupabaseClientConfig) (*SupabaseClient, error) {
	if cfg.ApiSecret == "" || cfg.ApiUrl == "" {
		return nil, fmt.Errorf("Supabase Api Key & Api Url Cannot Be Empty")
	}

	client, err := supabase.NewClient(cfg.ApiUrl, cfg.ApiSecret, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed To Initialize Supabase Client: %v", err)
	}

	return &SupabaseClient{
		client: client,
	}, nil
}

func (s *SupabaseClient) GetClient() *supabase.Client {
	return s.client
}
