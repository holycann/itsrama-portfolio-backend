package base

import (
	"context"

	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
)

// BaseService defines a generic interface for service operations
type BaseService[T any, R any] interface {
	FindByID(ctx context.Context, id string) (*R, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, opts ListOptions) ([]R, error)
	Search(ctx context.Context, opts ListOptions) ([]R, int, error)
	BulkCreate(ctx context.Context, values []*T) ([]T, error)
	BulkUpdate(ctx context.Context, values []*T) ([]T, error)
	BulkDelete(ctx context.Context, ids []string) error
}

// BaseServiceImpl provides a base implementation for common service operations
type BaseServiceImpl[T any, R any] struct {
	logger *logger.Logger
}

// NewBaseService creates a new base service
func NewBaseService[T any, R any](logger *logger.Logger) *BaseServiceImpl[T, R] {
	return &BaseServiceImpl[T, R]{
		logger: logger,
	}
}
