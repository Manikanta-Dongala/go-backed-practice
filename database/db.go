package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}
// generic repository
type Repository[T any] struct {
	Collection *mongo.Collection
}

// new database creation
func NewDatabase(uri, dbName string) (*Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	_ = client.Ping(ctx, nil)
	log.Print("connected to database succesfully!..")

	return &Database{
		Client: client,
		DB:     client.Database(dbName),
	}, nil

}

// close function for database
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}

// new repository creation
func NewRepository[T any](db *Database, collectionName string) *Repository[T] {
	return &Repository[T]{
		Collection: db.DB.Collection(collectionName),
	}
}

// creating new document
func (r *Repository[T]) Create(ctx context.Context, document T) error {
	_, err := r.Collection.InsertOne(ctx, document)
	return err
}

// find documents based on id
func (r *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var result T
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// find documents based on filter or all
func (r *Repository[T]) Find(ctx context.Context, filter interface{}) ([]*T, error) {
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T

	for cursor.Next(ctx) {
		var document T
		err := cursor.Decode(*&document)
		if err != nil {
			return nil, err
		}
		results = append(results, &document)
	}
	return results, nil
}

func (r *Repository[T]) Update(ctx context.Context, id string, update interface{}) error {
    _, err := r.Collection.UpdateOne(ctx, map[string]string{"_id": id}, update)
    return err
}

func (r *Repository[T]) Delete(ctx context.Context, id string) error {
    _, err := r.Collection.DeleteOne(ctx, map[string]string{"_id": id})
    return err
}

func (r *Repository[T]) Aggregate(ctx context.Context, pipeline []interface{}) ([]*T, error) {
    cursor, err := r.Collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []*T
    for cursor.Next(ctx) {
        var elem T
        err := cursor.Decode(&elem)
        if err != nil {
            return nil, err
        }
        results = append(results, &elem)
    }

    return results, nil
}