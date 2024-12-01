package userdb

import (
	"clanplan/server/bus/domain/userbus"
	"context"

	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store implements the Storer interface to manage the interactions between the Business layer and the database.
type Store struct {
	coll *mongo.Collection
	log  *logger.Logger
}

func NewStore(col *mongo.Collection, log *logger.Logger) *Store {
	return &Store{coll: col, log: log}
}

func (s *Store) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	bsonFilter := createQueryFilter(filter)
	count, err := s.coll.CountDocuments(ctx, bsonFilter)
	if err != nil {
		s.log.Info(ctx, "userdb.Count", "Filter", bsonFilter, "Error", err)
		return 0, err
	}
	s.log.Info(ctx, "userdb.Count", "Filter", bsonFilter, "Result", count)
	return int(count), nil
}

func (s *Store) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, page page.Page) ([]userbus.User, error) {
	bsonFilter := createQueryFilter(filter)

	sortDir := 1
	if orderBy.Direction == "DESC" {
		sortDir = -1
	}
	sort := bson.D{
		{
			Key:   orderBy.Field,
			Value: sortDir,
		},
	}

	opts := options.Find()
	opts.SetSort(sort)
	opts.SetLimit(int64(page.RowsPerPage()))
	opts.SetSkip(int64(page.RowsPerPage()) * int64(page.Number()-1))

	cursor, err := s.coll.Find(ctx, bsonFilter, opts)
	if err != nil {
		s.log.Info(ctx, "userdb.Query", "Filter", bsonFilter, "Error", err)
		return nil, err
	}

	var res []user
	if err = cursor.All(ctx, &res); err != nil {
		s.log.Info(ctx, "userdb.Query", "Filter", bsonFilter, "Error", err)
		return nil, err
	}
	s.log.Info(ctx, "userdb.Query", "Filter", bsonFilter, "Results", res)
	return toBusUsers(res), nil
}

func (s *Store) QueryById(ctx context.Context, id uuid.UUID) (userbus.User, error) {
	filter := bson.D{{"_id", id}}

	var result user
	err := s.coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		s.log.Info(ctx, "userdb.QueryById", "Filter", filter, "Error", err)
		return userbus.User{}, err
	}

	s.log.Info(ctx, "userdb.QueryById", "Filter", filter, "Result", result)
	return toBusUser(result), nil
}

func (s *Store) QueryOne(ctx context.Context, filter userbus.QueryFilter) (userbus.User, error) {
	bsonFilter := createQueryFilter(filter)

	var result user
	err := s.coll.FindOne(ctx, bsonFilter).Decode(&result)
	if err != nil {
		s.log.Info(ctx, "userdb.QueryOne", "Filter", bsonFilter, "Error", err)
		return userbus.User{}, err
	}

	s.log.Info(ctx, "userdb.QueryOne", "Filter", bsonFilter, "Result", result)
	return toBusUser(result), nil
}

func (s *Store) DeleteById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objectID}}

	_, err = s.coll.DeleteOne(ctx, filter)
	if err != nil {
		s.log.Info(ctx, "userdb.DeleteById", "Filter", filter, "Error", err)
		return err
	}
	s.log.Info(ctx, "userdb.DeleteById", "Filter", filter, "Results", "Success")
	return nil
}

func (s *Store) Delete(ctx context.Context, filter userbus.QueryFilter) error {
	bsonFilter := createQueryFilter(filter)

	_, err := s.coll.DeleteMany(ctx, bsonFilter)
	if err != nil {
		s.log.Info(ctx, "userdb.Delete", "Filter", bsonFilter, "Error", err)
		return err
	}
	s.log.Info(ctx, "userdb.Delete", "Filter", bsonFilter, "Result", "Success")
	return nil
}

func (s *Store) DeleteOne(ctx context.Context, filter userbus.QueryFilter) error {
	bsonFilter := createQueryFilter(filter)

	_, err := s.coll.DeleteOne(ctx, bsonFilter)
	if err != nil {
		s.log.Info(ctx, "userdb.DeleteOne", "Filter", bsonFilter, "Error", err)
		return err
	}
	s.log.Info(ctx, "userdb.DeleteOne", "Filter", bsonFilter, "Result", "Success")
	return nil
}

func (s *Store) Insert(ctx context.Context, u userbus.User) error {
	dbUser := toDbUser(u)
	_, err := s.coll.InsertOne(ctx, dbUser)
	if err != nil {
		s.log.Info(ctx, "userdb.Insert", "User", dbUser, "Error", err)
		return err
	}
	s.log.Info(ctx, "userdb.Insert", "User", dbUser, "Result", "Success")
	return nil
}

func (s *Store) Update(ctx context.Context, u userbus.User) error {
	dbUser := toDbUser(u)

	filter := bson.D{{"_id", dbUser.ID}}
	update := bson.D{{"$set", dbUser}}

	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		s.log.Info(ctx, "userdb.Update", "Filter", filter, "Update", update, "Error", err)
		return err
	}

	s.log.Info(ctx, "userdb.Update", "Filter", filter, "Update", update)
	return nil
}

func createQueryFilter(filter userbus.QueryFilter) *bson.D {
	bsonFilter := bson.D{}

	if filter.ID != nil {
		bsonFilter = append(bsonFilter, bson.E{"_id", *filter.ID})
	}

	if filter.Username != nil {
		bsonFilter = append(bsonFilter, bson.E{"name", filter.Username})
	}

	if filter.Email != nil {
		bsonFilter = append(bsonFilter, bson.E{"email", filter.Email.String()})
	}

	if filter.StartCreatedDate != nil {
		bsonFilter = append(bsonFilter, bson.E{"date_created", bson.D{{"$gte", filter.StartCreatedDate}}})
	}

	if filter.EndCreatedDate != nil {
		bsonFilter = append(bsonFilter, bson.E{"date_created", bson.D{{"$lte", filter.EndCreatedDate}}})
	}

	return &bsonFilter
}
