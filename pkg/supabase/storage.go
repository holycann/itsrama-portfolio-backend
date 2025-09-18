package supabase

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

// StorageConfig provides comprehensive configuration for Supabase storage
type StorageConfig struct {
	// Supabase project details
	ProjectID    string
	JwtApiSecret string

	// Storage bucket configuration
	BucketID      string
	DefaultFolder string

	// Optional custom headers for client initialization
	Headers map[string]string

	// Advanced configuration options
	MaxFileSize         int64
	AllowedFileTypes    []string
	DefaultCacheControl string
}

// SupabaseStorage provides enhanced storage management capabilities
type SupabaseStorage struct {
	client *storage_go.Client
	Config StorageConfig
}

// NewSupabaseStorage creates an enhanced Supabase storage client with robust configuration
func NewSupabaseStorage(cfg StorageConfig) (*SupabaseStorage, error) {
	// Validate required configuration parameters
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if cfg.JwtApiSecret == "" {
		return nil, fmt.Errorf("JWT API secret is required")
	}

	// Use default headers if not provided
	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}

	// Set default max file size if not specified
	if cfg.MaxFileSize <= 0 {
		cfg.MaxFileSize = 10 * 1024 * 1024 // 10MB default
	}

	// Set default allowed file types if not specified
	if len(cfg.AllowedFileTypes) == 0 {
		cfg.AllowedFileTypes = []string{
			"image/jpeg", "image/png", "image/webp", "application/pdf",
		}
	}

	// Set default cache control if not specified
	if cfg.DefaultCacheControl == "" {
		cfg.DefaultCacheControl = "public, max-age=3600, must-revalidate"
	}

	// Construct Supabase storage client URL
	storageURL := fmt.Sprintf("https://%s.supabase.co/storage/v1", cfg.ProjectID)

	// Create storage client with comprehensive configuration
	storageClient := storage_go.NewClient(
		storageURL,
		cfg.JwtApiSecret,
		cfg.Headers,
	)

	// Return configured Supabase storage instance
	return &SupabaseStorage{
		client: storageClient,
		Config: cfg,
	}, nil
}

// Upload handles file upload with comprehensive validation and storage
func (s *SupabaseStorage) Upload(
	ctx context.Context,
	file *multipart.FileHeader,
	path string,
	opts ...storage_go.FileOptions,
) (string, error) {
	// Validate file size
	if file.Size > s.Config.MaxFileSize {
		return "", fmt.Errorf("file size %d bytes exceeds maximum limit of %d",
			file.Size, s.Config.MaxFileSize)
	}

	// Validate file type
	fileType := file.Header.Get("Content-Type")
	if !s.isAllowedFileType(fileType) {
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}

	// Explicitly check and handle Close error
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			// Log the close error or handle it as needed
			log.Printf("Failed to close file: %v", closeErr)
		}
	}()

	// Prepare file options
	fileOpts := storage_go.FileOptions{
		Upsert:       boolPtr(true),
		CacheControl: stringPtr(s.Config.DefaultCacheControl),
		ContentType:  stringPtr(fileType),
	}

	// Merge with any provided options
	if len(opts) > 0 {
		fileOpts = mergeFileOptions(fileOpts, opts[0])
	}

	if !strings.HasPrefix(path, s.Config.DefaultFolder) {
		path = filepath.Clean(filepath.Join(s.Config.DefaultFolder, path))
	}
	path = filepath.ToSlash(path)

	// Upload file
	_, err = s.client.UploadFile(
		s.Config.BucketID,
		path,
		src,
		fileOpts,
	)
	if err != nil {
		return "", err
	}
	return path, nil
}

// mergeFileOptions combines default and custom file options
func mergeFileOptions(defaultOpts, customOpts storage_go.FileOptions) storage_go.FileOptions {
	if customOpts.Upsert != nil {
		defaultOpts.Upsert = customOpts.Upsert
	}
	if customOpts.CacheControl != nil {
		defaultOpts.CacheControl = customOpts.CacheControl
	}
	if customOpts.ContentType != nil {
		defaultOpts.ContentType = customOpts.ContentType
	}
	return defaultOpts
}

// isAllowedFileType checks if the file type is in the allowed list
func (s *SupabaseStorage) isAllowedFileType(fileType string) bool {
	// If no allowed types specified, allow all
	if len(s.Config.AllowedFileTypes) == 0 {
		return true
	}

	// Normalize file type
	fileType = strings.ToLower(strings.TrimSpace(fileType))

	// Check if file type matches any allowed type
	for _, allowedType := range s.Config.AllowedFileTypes {
		if strings.ToLower(allowedType) == fileType {
			return true
		}
	}

	return false
}

// GetFileType determines the general category of a file based on its MIME type
func (s *SupabaseStorage) GetFileType(fileType string) string {
	// Normalize file type
	fileType = strings.ToLower(strings.TrimSpace(fileType))

	// Image types
	imagePrefixes := []string{"image/"}
	for _, prefix := range imagePrefixes {
		if strings.HasPrefix(fileType, prefix) {
			return "image"
		}
	}

	// Document types
	documentTypes := []string{
		"application/pdf",
	}
	for _, docType := range documentTypes {
		if fileType == docType {
			return "document"
		}
	}

	return "other"
}

// Delete removes a file from storage
func (s *SupabaseStorage) Delete(
	ctx context.Context,
	filepath string,
) (string, error) {
	_, err := s.client.RemoveFile(
		s.Config.BucketID,
		[]string{filepath},
	)
	if err != nil {
		return "", err
	}

	// Return success message as a string
	return "File deleted successfully", nil
}

// ListFiles retrieves files in a specific path
func (s *SupabaseStorage) ListFiles(
	ctx context.Context,
	path string,
	opts ...storage_go.FileSearchOptions,
) ([]storage_go.FileObject, error) {
	searchOpts := storage_go.FileSearchOptions{}
	if len(opts) > 0 {
		searchOpts = opts[0]
	}

	return s.client.ListFiles(
		s.Config.BucketID,
		path,
		searchOpts,
	)
}

// GetPublicURL generates a public URL for a file
func (s *SupabaseStorage) GetPublicURL(
	path string,
	transformOpts ...storage_go.UrlOptions,
) (string, error) {
	if !strings.HasPrefix(path, s.Config.DefaultFolder) {
		path = filepath.Clean(filepath.Join(s.Config.DefaultFolder, path))
	}
	path = filepath.ToSlash(path)

	resp := s.client.GetPublicUrl(
		s.Config.BucketID,
		path,
		transformOpts...,
	)
	return resp.SignedURL, nil
}

// Helper functions to create pointers for optional values
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
