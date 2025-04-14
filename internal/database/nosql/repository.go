package nosql

import (
	"context"
	"errors"
	"fmt"
	"go-service-template/internal/database"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoSQLRepository[T any] struct {
	collection *mongo.Collection
}

// var _ database.Repository[T] = *(NoSQLRepository[T])(nil)

// NewNoSQLRepository creates a new NoSQLRepository instance
func NewNoSQLRepository[T any](client *mongo.Client, dbName, collectionName string) (*NoSQLRepository[T], error) {
	if client == nil {
		return nil, &database.DatabaseError{
			Operation: "NewNoSQLRepository",
			Err:       errors.New("MongoDB client cannot be nil")}
	}

	if len(strings.TrimSpace(dbName)) == 0 || len(strings.TrimSpace(collectionName)) == 0 {
		return nil, &database.DatabaseError{
			Operation: "NewNoSQLRepository",
			Err:       errors.New("database name and collection name must not be empty or whitespace")}
	}

	return &NoSQLRepository[T]{collection: client.Database(dbName).Collection(collectionName)}, nil
}

func (r *NoSQLRepository[T]) Save(ctx context.Context, entity *T) error {
	if entity == nil {
		return &database.DatabaseError{Operation: "Save", Err: errors.New("entity cannot be nil")}
	}
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return &database.DatabaseError{
			Operation: "Save",
			Err:       fmt.Errorf("failed to insert entity: %v, error: %w", entity, err)}
	}

	return nil
}

// Atomic Save: Fails on the first error
// need to return the idx that fail on it
func (r *NoSQLRepository[T]) SaveAtomic(ctx context.Context, entities []*T) (*T, error) {
	// Return error if the input slice is empty
	if len(entities) == 0 {
		return nil, &database.DatabaseError{
			Operation: "saveAtomic in SaveAll",
			Err:       errors.New("entities cannot be empty")}
	}

	// Convert []*T to []interface{} for InsertMany
	entitiesInterface := make([]interface{}, len(entities))
	for idx, entity := range entities {
		entitiesInterface[idx] = entity
	}

	// InsertMany in ordered mode: stops at the first error (atomic-like behavior)
	opts := options.InsertMany().SetOrdered(true)

	// Perform the insertion
	_, err := r.collection.InsertMany(ctx, entitiesInterface, opts)
	if err == nil {
		// All documents inserted successfully
		return nil, nil
	}

	// Check if it's a BulkWriteException to extract which insert failed
	bulkErr, ok := err.(mongo.BulkWriteException)
	if !ok {
		// Not a bulk write error
		return nil, &database.DatabaseError{
			Operation: "saveAtomic in SaveAll",
			Err:       fmt.Errorf("failed to insert: %w", err)}
	}

	// If there are write errors, return the first failed entity
	if len(bulkErr.WriteErrors) > 0 {
		firstEntityFailIndex := bulkErr.WriteErrors[0].Index
		firstEntityFail := entities[firstEntityFailIndex]

		return firstEntityFail, &database.DatabaseError{
			Operation: fmt.Sprintf("saveAtomic failed at index %v", firstEntityFail),
			Err:       bulkErr.WriteErrors[0],
		}
	}
	// No specific write error details found, fallback
	return nil, nil
}

// Partial Success Save: Inserts as many documents as possible,
func (r *NoSQLRepository[T]) SavePartialSuccess(ctx context.Context, entities []*T) ([]*T, error) {
	// Return error if the input slice is empty
	if len(entities) == 0 {
		return nil, &database.DatabaseError{
			Operation: "savePartialSuccess in SaveAll",
			Err:       errors.New("entities cannot be empty")}
	}

	// Convert []*T to []interface{} for InsertMany
	entitiesInterface := make([]interface{}, len(entities))
	for idx, entity := range entities {
		entitiesInterface[idx] = entity
	}

	// Allow partial success
	opt := options.InsertMany().SetOrdered(false)

	_, err := r.collection.InsertMany(ctx, entitiesInterface, opt)
	if err == nil {
		// All documents inserted successfully
		return nil, nil
	}

	// Handle bulk write errors to determine which inserts failed
	bulkErr, ok := err.(mongo.BulkWriteException)
	if !ok {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	// Track failed indexes
	failedIndex := make(map[int]bool)
	for _, we := range bulkErr.WriteErrors {
		failedIndex[we.Index] = true
	}

	// Collect failed entities
	var failedEntities []*T
	for idx, entity := range entities {
		if _, ok := failedIndex[idx]; ok {
			failedEntities = append(failedEntities, entity)
		}
	}

	if len(failedEntities) > 0 {
		return failedEntities, &database.DatabaseError{
			Operation: "SavePartialSuccess",
			Err:       fmt.Errorf("partial insert failure: %d documents failed", len(failedEntities))}
	}

	return nil, nil
}

// // Transaction Save: Full atomicity with a transaction
// func (r *NoSQLRepository[T]) saveWithTransaction(ctx context.Context, entities []*T) error {
// 	if len(entities) == 0 {
// 		return &database.DatabaseError{Operation: "saveWithTransaction in SaveAll", Err: errors.New("entities cannot be empty")}
// 	}

// }

func (r *NoSQLRepository[T]) SaveAll(ctx context.Context, entities []*T) error {
	// Check if the input slice is empty
	if len(entities) == 0 {
		return &database.DatabaseError{
			Operation: "SaveAll",
			Err:       fmt.Errorf("entities slice is empty (length: %d)", len(entities))}
	}

	// Convert []*T to []interface{}
	interfaceEntities := make([]interface{}, len(entities))
	for idx, entity := range entities {
		interfaceEntities[idx] = entity
	}

	// Prepare options for InsertMany operation (ordered: true for atomic insertions)
	opts := options.InsertMany().SetOrdered(true)

	_, err := r.collection.InsertMany(ctx, interfaceEntities, opts)
	if err != nil {
		return &database.DatabaseError{Operation: "SaveAll", Err: fmt.Errorf("failed to insert documents: %v", err)}
	}
	return nil
}

// // FindByID retrieves an entity by ID
// func (r *NoSQLRepository[Model]) FindByID(ctx context.Context, id interface{}) (*Model, error) {
// 	var result Model
// 	filter := map[string]interface{}{
// 		"_id": id,
// 	}
// 	err := r.collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, &database.DatabaseError{Operation: "FindByID", Err: fmt.Errorf("document with id %v not found", id)}
// 		}
// 		return nil, &database.DatabaseError{Operation: "FindByID", Err: err}

// 	}
// 	return &result, nil
// }

// // FindByFilter retrieves entities matching a filter
// func (r *NoSQLRepository[Model]) FindByFilter(ctx context.Context, filter interface{}) ([]Model, error) {
// 	cursor, err := r.collection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, &database.DatabaseError{Operation: "FindByFilter", Err: err}
// 	}
// 	defer cursor.Close(ctx)

// 	var results []Model
// 	for cursor.Next(ctx) {
// 		var entity Model
// 		if err := cursor.Decode(&entity); err != nil {
// 			return nil, &database.DatabaseError{Operation: "FindByFilter", Err: err}
// 		}
// 		results = append(results, entity)
// 	}

// 	if cursor.Err() != nil {
// 		return nil, &database.DatabaseError{Operation: "FindByFilter", Err: cursor.Err()}
// 	}

// 	return results, nil
// }

// // FindAll retrieves all documents
// func (r *NoSQLRepository[Model]) FindAll(ctx context.Context) ([]Model, error) {
// 	return r.FindByFilter(ctx, nil)
// }

// // Update modifies an existing entity
// func (r *NoSQLRepository[Model]) Update(ctx context.Context, id interface{}, entity *Model) error {
// 	if id == nil {
// 		return &database.DatabaseError{Operation: "Update", Err: fmt.Errorf("id is nil")}
// 	}

// 	if entity == nil {
// 		return &database.DatabaseError{Operation: "Update", Err: fmt.Errorf("entity is nil")}
// 	}

// 	filter := map[string]interface{}{
// 		"_id": id,
// 	}

// 	update := map[string]interface{}{
// 		"$set": entity,
// 	}

// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return &database.DatabaseError{Operation: "Update", Err: err}
// 	}
// 	return nil
// }

// // Delete removes a document by ID
// func (r *NoSQLRepository[Model]) Delete(ctx context.Context, id interface{}) error {
// 	if id == nil {
// 		return &database.DatabaseError{Operation: "Delete", Err: fmt.Errorf("id is nil")}
// 	}

// 	filter := map[string]interface{}{
// 		"_id": id,
// 	}
// 	_, err := r.collection.DeleteOne(ctx, filter)
// 	if err != nil {
// 		return &database.DatabaseError{Operation: "Delete", Err: err}
// 	}
// 	return nil
// }

// // Count returns the number of documents
// func (r *NoSQLRepository[Model]) Count(ctx context.Context) (int64, error) {
// 	return r.collection.CountDocuments(ctx, nil)
// }

// // CountByFilter returns the count of entities matching the given filter.
// func (r *NoSQLRepository[Model]) CountByFilter(ctx context.Context, filter interface{}) (int64, error) {
// 	count, err := r.collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		return 0, &database.DatabaseError{Operation: "CountByFilter", Err: err}
// 	}
// 	return count, nil
// }
// 		return 0, &database.DatabaseError{Operation: "CountByFilter", Err: err}
// 	}
// 	return count, nil
// }

// // ExistsByID checks if a document exists
// func (r *NoSQLRepository[Model]) ExistsByID(ctx context.Context, id interface{}) (bool, error) {
// 	filter := map[string]interface{}{
// 		"_id": id,
// 	}

// 	count, err := r.collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		return false, &database.DatabaseError{Operation: "ExistsByID", Err: err}
// 	}
// 	return count > 0, nil
// }

// // ExistsByFilter checks if any entity matches the given filter criteria.
// func (r *NoSQLRepository[Model]) ExistsByFilter(ctx context.Context, filter interface{}) (bool, error) {
// 	if filter == nil {
// 		return false, &database.DatabaseError{Operation: "ExistsByFilter", Err: fmt.Errorf("filter is nil")}
// 	}

// 	count, err := r.collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		return false, &database.DatabaseError{Operation: "ExistsByFilter", Err: err}
// 	}

// 	return count > 0, nil
// }
