package nosqldb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Db struct {
    db *mongo.Database
}

func NewDb(db *mongo.Database) *Db {
    return &Db{db: db}
}

// Collection returns a mongo collection
func (d *Db) Collection(name string) *mongo.Collection {
    return d.db.Collection(name)
}




