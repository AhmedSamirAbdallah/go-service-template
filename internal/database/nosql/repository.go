package nosql

import (
	"context"
	"fmt"
	"go-service-template/internal/database"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)
type NoSQLRepository[Model any] struct {
	collection *mongo.Collection
}
var _ database.Repository[Model] = *(NoSQLRepository[Model])(nil)

// NewNoSQLRepository creates a new NoSQLRepository instance
func NewNoSQLRepository[Model any](client *mongo.Client, dbName, collectionName string) (*NoSQLRepository[Model], error) {
	if client == nil {
		return nil, &database.DatabaseError{Operation: "NewNoSQLRepository", Err: fmt.Errorf("MongoDB client cannot be nil")}
	}
	if len(strings.TrimSpace(dbName)) == 0 || len(strings.TrimSpace(collectionName)) == 0 {
		return nil, &database.DatabaseError{Operation: "NewNoSQLRepository", Err: fmt.Errorf("database name and collection name must not be empty or whitespace")}

	}

	return &NoSQLRepository[Model]{collection: client.Database(dbName).Collection(collectionName)}, nil
}

// Save inserts a new entity into MongoDB
func (r *NoSQLRepository[Model]) Save(ctx context.Context, entity Model) error {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return &database.DatabaseError{Operation: "Save", Err: err}
	}
	return nil
}

// FindByID retrieves an entity by ID
func (r *NoSQLRepository[Model]) FindByID(ctx context.Context, id interface{}) (*Model, error) {
	var result Model
	filter := map[string]interface{}{
		"_id": id,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &database.DatabaseError{Operation: "FindByID", Err: fmt.Errorf("document with id %v not found", id)}
		}
		return nil, &database.DatabaseError{Operation: "FindByID", Err: err}

	}
	return &result, nil
}

// FindByFilter retrieves entities matching a filter
func (r *NoSQLRepository[Model]) FindByFilter(ctx context.Context, filter interface{}) ([]Model, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, &database.DatabaseError{Operation: "FindByFilter", Err: err}
	}
	defer cursor.Close(ctx)

	var results []Model
	for cursor.Next(ctx) {
		var entity Model
		if err := cursor.Decode(&entity); err != nil {
			return nil, &database.DatabaseError{Operation: "FindByFilter", Err: err}
		}
		results = append(results, entity)
	}

	if cursor.Err() != nil {
		return nil, &database.DatabaseError{Operation: "FindByFilter", Err: cursor.Err()}
	}

	return results, nil
}

// FindAll retrieves all documents
func (r *NoSQLRepository[Model]) FindAll(ctx context.Context) ([]Model, error) {
	return r.FindByFilter(ctx, nil)
}

// Update modifies an existing entity
func (r *NoSQLRepository[Model]) Update(ctx context.Context, id interface{}, entity *Model) error {
	if id == nil {
		return &database.DatabaseError{Operation: "Update", Err: fmt.Errorf("id is nil")}
	}

	if entity == nil {
		return &database.DatabaseError{Operation: "Update", Err: fmt.Errorf("entity is nil")}
	}

	filter := map[string]interface{}{
		"_id": id,
	}

	update := map[string]interface{}{
		"$set": entity,
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &database.DatabaseError{Operation: "Update", Err: err}
	}
	return nil
}

// Delete removes a document by ID
func (r *NoSQLRepository[Model]) Delete(ctx context.Context, id interface{}) error {
	if id == nil {
		return &database.DatabaseError{Operation: "Delete", Err: fmt.Errorf("id is nil")}
	}

	filter := map[string]interface{}{
		"_id": id,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &database.DatabaseError{Operation: "Delete", Err: err}
	}
	return nil
}

// Count returns the number of documents
func (r *NoSQLRepository[Model]) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, nil)
}

// CountByFilter returns the count of entities matching the given filter.
func (r *NoSQLRepository[Model]) CountByFilter(ctx context.Context, filter interface{}) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, &database.DatabaseError{Operation: "CountByFilter", Err: err}
	}
	return count, nil
}
		return 0, &database.DatabaseError{Operation: "CountByFilter", Err: err}
	}
	return count, nil
}

// ExistsByID checks if a document exists
func (r *NoSQLRepository[Model]) ExistsByID(ctx context.Context, id interface{}) (bool, error) {
	filter := map[string]interface{}{
		"_id": id,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, &database.DatabaseError{Operation: "ExistsByID", Err: err}
	}
	return count > 0, nil
}

// ExistsByFilter checks if any entity matches the given filter criteria.
func (r *NoSQLRepository[Model]) ExistsByFilter(ctx context.Context, filter interface{}) (bool, error) {
	if filter == nil {
		return false, &database.DatabaseError{Operation: "ExistsByFilter", Err: fmt.Errorf("filter is nil")}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, &database.DatabaseError{Operation: "ExistsByFilter", Err: err}
	}

	return count > 0, nil
}
