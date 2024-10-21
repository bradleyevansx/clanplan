package main

import (
	userdb "clanplan/server/bus/domain/userbus/stores"
	"clanplan/server/bus/sdk/nosqldb"
	"context"
	"log"
	"os"

	"github.com/ardanlabs/service/foundation/logger"
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
	ardanLog := logger.New(os.Stdout, logger.LevelInfo, "Store", nil)
	userRepo := userdb.NewStore(userColl, ardanLog)
	// id := uuid.New()
	// log.Printf("ID: %s\n", id)
	// newUser := userbus.User{
	// 	Username: "to receive id",
	// 	Email: mail.Address{
	// 		Address: "Test",
	// 	},
	// 	PasswordHash: []byte("Test"),
	// 	Enabled:      true,
	// 	DateCreated:  time.Now(),
	// 	DateUpdated:  time.Now(),
	// }
	// err = userRepo.Insert(context.TODO(), newUser)
	idString := "2d850ba3-6bce-4b23-86a1-c53df3ec1901"
	id := uuid.MustParse(idString)

	_, err = userRepo.QueryById(context.TODO(), id)
	if err != nil {
		log.Fatal(err)
	}

}
