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

var MongoConnection *mongo.Client
var UserCollection *mongo.Collection
var EventsCollection *mongo.Collection

func Connect() {

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

	MongoConnection = client
	UserCollection = MongoConnection.Database("calendar_db").Collection("users")
	EventsCollection = MongoConnection.Database("calendar_db").Collection("events")
}
