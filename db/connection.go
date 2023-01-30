package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {

	err := godotenv.Load()
	if err != nil {
		print("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION")))

	if err != nil {
		log.Fatal("Error here:", err)
	}

	print("Connection established")

	return client
}

var UserCollection = Connect().Database("calendar_db").Collection("users")
var EventsCollection = Connect().Database("calendar_db").Collection("events")
