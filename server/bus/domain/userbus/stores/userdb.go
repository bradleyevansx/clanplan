package db

import "go.mongodb.org/mongo-driver/mongo"


type Store struct {
	col *mongo.Collection
}

func NewStore(col *mongo.Collection) *Store {
	return &Store{col: col}
}