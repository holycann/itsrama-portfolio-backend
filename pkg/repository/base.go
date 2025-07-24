package repository

import (
	"context"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, model *T) error
	FindByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]T, error)
	Count(ctx context.Context) (int, error)
}
