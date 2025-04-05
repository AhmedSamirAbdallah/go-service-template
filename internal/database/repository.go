package database

import "context"

type Repository[T any] interface {
	Save(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	FindByFilter(ctx context.Context, filter interface{}) ([]T, error)
	FindAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, id interface{}, entity *T) error
	Delete(ctx context.Context, id interface{}) error
	Count(ctx context.Context) (int64, error)
	ExistsByID(ctx context.Context, id interface{}) (bool, error)
}
