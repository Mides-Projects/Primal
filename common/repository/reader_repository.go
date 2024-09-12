package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReaderRepository[T any] struct {
	col *mongo.Collection

	wrapper func(result *mongo.SingleResult) T
}

func (r *ReaderRepository[T]) FindById(id string) T {
	result := r.col.FindOne(context.Background(), bson.M{"_id": id})
	if result.Err() != nil {
		return r.wrapper(nil)
	}

	return r.wrapper(result)
}

func (r *ReaderRepository[T]) Find(filter bson.M) T {
	result := r.col.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return r.wrapper(nil)
	}

	return r.wrapper(result)
}
