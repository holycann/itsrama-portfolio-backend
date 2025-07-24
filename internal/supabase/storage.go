package supabase

import (
	"fmt"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorage struct {
	client   *storage_go.Client
	bucketID string
}

type SupabaseStorageConfig struct {
	JwtApiSecret string
	ProjectID    string
	BucketID     string
}

func NewSupabaseStorage(cfg SupabaseStorageConfig) *SupabaseStorage {
	storageClient := storage_go.NewClient(fmt.Sprintf("https://%s.supabase.co/storage/v1", cfg.ProjectID), cfg.JwtApiSecret, nil)

	return &SupabaseStorage{
		client:   storageClient,
		bucketID: cfg.BucketID,
	}
}

func (s *SupabaseStorage) GetClient() *storage_go.Client {
	return s.client
}

func (s *SupabaseStorage) GetBucketID() string {
	return s.bucketID
}
