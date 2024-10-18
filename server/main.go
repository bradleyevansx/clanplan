package main

import (
	"clanplan/server/bus/domain/userbus"
	userdb "clanplan/server/bus/domain/userbus/stores"
	"clanplan/server/bus/sdk/nosqldb"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	dbClient := client.Database("clanplan")
	mongoDb := nosqldb.NewDb(dbClient)
	userColl := mongoDb.Collection("users")
	userRepo := userdb.NewStore(userColl)
	// id := uuid.New()
	// log.Printf("ID: %s\n", id)
	// newUser := userbus.User{
	// 	ID:       id,
	// 	Username: "Test Create",
	// 	Email: mail.Address{
	// 		Address: "Test",
	// 	},
	// 	PasswordHash: []byte("Test"),
	// 	Enabled:      true,
	// 	DateCreated:  time.Now(),
	// 	DateUpdated:  time.Now(),
	// }
	// err = userRepo.Insert(newUser)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	queryId, err := uuid.Parse("463507a9-11a7-4a01-8a9f-778c28c9609e")
	if err != nil {
		log.Fatal(err)
	}
	userFilter := userbus.QueryFilter{
		ID: &queryId,
	}
	res, err := userRepo.Query(userFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

}
