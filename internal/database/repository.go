package database

import "context"

type Repository[T any] interface {
	// Create
	Save(ctx context.Context, entity *T) error
	SaveAll(ctx context.Context, entities []*T) error
	
	//Update
	Update(ctx context.Context, id interface{}, entity *T) error
	UpdateByFilter(ctx context.Context, filter interface{}, entity *T) error
	BulkUpdate(ctx context.Context, filters []interface{}, updates []*T) error
	Upsert(ctx context.Context, filter interface{}, entity *T) error

	// Read
	Find(ctx context.Context, filter interface{}, opts *QueryOptions) ([]T, error)
	FindOne(ctx context.Context, filter interface{}) (*T, error)
	FindByID(ctx context.Context, id interface{}) (*T, error)
	FindByIDs(ctx context.Context, ids []interface{}) ([]T, error)
	FindByFilter(ctx context.Context, filter interface{}) ([]T, error)
	FindAll(ctx context.Context) ([]T, error)
	FindPaginated(ctx context.Context, filter interface{}, opts *QueryOptions) ([]T, int64, error)
	Distinct(ctx context.Context, field string, filter interface{}) ([]interface{}, error)

	// Delete
	Delete(ctx context.Context, id interface{}) error
	DeleteByFilter(ctx context.Context, filter interface{}) error
	DeleteAllRecords(ctx context.Context) error
	BulkDelete(ctx context.Context, filters []interface{}) error

	// Stats
	Count(ctx context.Context) (int64, error)
	CountByFilter(ctx context.Context, filter interface{}) (int64, error)
	CountWithOptions(ctx context.Context, filter interface{}, opts *QueryOptions) (int64, error)
	ExistsByID(ctx context.Context, id interface{}) (bool, error)
	ExistsByFilter(ctx context.Context, filter interface{}) (bool, error)

	// Advanced Queries
	Project(ctx context.Context, filter interface{}, projection interface{}, result interface{}) error
	Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error
	RawQuery(ctx context.Context, query interface{}, result interface{}) error

	// Transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// Utility
	SetQueryTimeout(ctx context.Context, timeout int) context.Context
}
