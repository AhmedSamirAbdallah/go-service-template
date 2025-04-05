package nosql

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type NoSQLRepository[Model any] struct {
	collection *mongo.Collection
}

// NewNoSQLRepository creates a new NoSQLRepository instance
func NewNoSQLRepository[Model any](client *mongo.Client, dbName, collectionName string) (*NoSQLRepository[Model], error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client cannot be nil")
	}
	if dbName == "" || collectionName == "" {
		return nil, fmt.Errorf("database name and collection name must not be empty")
	}

	return &NoSQLRepository[Model]{collection: client.Database(dbName).Collection(collectionName)}, nil
}

// Save inserts a new entity into MongoDB
func (r *NoSQLRepository[Model]) Save(ctx context.Context, entity Model) error {
	_, err := r.collection.InsertOne(ctx, entity)
	return err
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
			return nil, fmt.Errorf("document with id %v not found", id)
		}
		return nil, err
	}
	return &result, nil
}

func (r *NoSQLRepository[Model]) FindByFilter(ctx context.Context, filter interface{}) ([]Model, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []Model
	for cursor.Next(ctx) {
		var entity Model
		if err := cursor.Decode(&entity); err != nil {
			return nil, err
		}
		results = append(results, entity)
	}

	return results, nil
}
