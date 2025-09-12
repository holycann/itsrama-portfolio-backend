package supabase

import (
	"fmt"

	supabaseAuth "github.com/supabase-community/auth-go"
)

type SupabaseAuth struct {
	auth supabaseAuth.Client
}

type SupabaseAuthConfig struct {
	ApiKey    string
	ProjectID string
}

func NewSupabaseAuth(cfg SupabaseAuthConfig) (*SupabaseAuth, error) {
	client := supabaseAuth.New(
		cfg.ProjectID,
		cfg.ApiKey,
	)

	if client == nil {
		return nil, fmt.Errorf("failed to initialize Supabase auth client")
	}

	return &SupabaseAuth{
		auth: client,
	}, nil
}

func (s *SupabaseAuth) GetClient() supabaseAuth.Client {
	return s.auth
}
