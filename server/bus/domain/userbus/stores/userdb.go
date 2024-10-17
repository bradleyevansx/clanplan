package userdb

import (
	"clanplan/server/bus/domain/userbus"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Store implements the Storer interface to manage the interactions between the Business layer and the database.
type Store struct {
	coll *mongo.Collection
}

func NewStore(col *mongo.Collection) *Store {
	return &Store{coll: col}
}

func (s *Store) Query(filter userbus.QueryFilter) {
	var user user
	err := s.coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {

	}
}

func (s *Store) QueryById(){}
func (s *Store) Delete(){}
func (s *Store) DeleteById(){}
func (s *Store) Insert(){}
func (s *Store) Update(){}