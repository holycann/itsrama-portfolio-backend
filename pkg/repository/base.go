package repository

import (
	"context"
)

// SortOrder defines the sorting direction
type SortOrder string

// Sorting order constants
const (
	SortAscending  SortOrder = "asc"
	SortDescending SortOrder = "desc"
)

// FilterOption represents a generic filter for repositories
type FilterOption struct {
	Field    string
	Operator string
	Value    interface{}
}

// ListOptions provides flexible querying parameters
type ListOptions struct {
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   SortOrder
	Filters     []FilterOption
	SearchQuery string
}

// BaseRepository provides a generic interface for CRUD and list operations
type BaseRepository[T any, R any] interface {
	// Basic CRUD operations
	Create(ctx context.Context, model *T) error
	FindByID(ctx context.Context, id string) (*R, error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, id string) error

	// Enhanced list and search methods
	List(ctx context.Context, opts ListOptions) ([]R, error)
	Search(ctx context.Context, opts ListOptions) ([]R, int, error)
	Count(ctx context.Context, filters []FilterOption) (int, error)

	// Additional utility methods
	Exists(ctx context.Context, id string) (bool, error)
	FindByField(ctx context.Context, field string, value interface{}) ([]R, error)
}
