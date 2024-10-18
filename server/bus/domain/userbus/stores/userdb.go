package userdb

import (
	"clanplan/server/bus/domain/userbus"
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store implements the Storer interface to manage the interactions between the Business layer and the database.
type Store struct {
	coll *mongo.Collection
}

func NewStore(col *mongo.Collection) *Store {
	return &Store{coll: col}
}

func (s *Store) Query(filter userbus.QueryFilter) ([]userbus.User, error) {
	bsonFilter := createQueryFilter(filter)

	cursor, err := s.coll.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}

	var results []user
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return toBusUsers(results), nil
}

func (s *Store) QueryById(id string) (*user, error) {

	filter := bson.D{{"_id", id}}

	var result user
	err := s.coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Store) QueryOne(filter userbus.QueryFilter) (*user, error) {
	bsonFilter := createQueryFilter(filter)

	var result user
	err := s.coll.FindOne(context.TODO(), bsonFilter).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Store) DeleteById(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objectID}}

	_, err = s.coll.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Delete(filter userbus.QueryFilter) error {
	bsonFilter := createQueryFilter(filter)

	_, err := s.coll.DeleteMany(context.TODO(), bsonFilter)

	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteOne(filter userbus.QueryFilter) error {
	bsonFilter := createQueryFilter(filter)

	_, err := s.coll.DeleteOne(context.TODO(), bsonFilter)

	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Insert(u userbus.User) error {
	dbUser := toDbUser(u)
	_, err := s.coll.InsertOne(context.TODO(), dbUser)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Update(u userbus.User) error {
	dbUser := toDbUser(u)

	filter := bson.D{{"_id", dbUser.ID}}
	update := bson.D{{"$set", dbUser}}

	_, err := s.coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func createQueryFilter(filter userbus.QueryFilter) *bson.D {
	bsonFilter := bson.D{}

	val := reflect.ValueOf(filter)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsNil() {
			continue
		}

		fieldName := typ.Field(i).Name
		fieldValue := field.Elem().Interface()

		bsonTag := typ.Field(i).Tag.Get("bson")
		if bsonTag == "" || bsonTag == "-" {
			bsonTag = fieldName
		}
		bsonFilter = append(bsonFilter, bson.E{Key: bsonTag, Value: fieldValue})
	}
	log.Printf("Filter: %v", bsonFilter)
	return &bsonFilter
}
