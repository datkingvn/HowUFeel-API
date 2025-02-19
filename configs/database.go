package configs

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	log.Println("Connecting to MongoDB...")

	var mongodbURI string
	if mongodbURI = os.Getenv("MONGODB_URI"); mongodbURI == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}

	clientOptions := options.Client().ApplyURI(mongodbURI)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client
}

var Client *mongo.Client = ConnectDB()

func GetCollection(collectionName string) *mongo.Collection {
	var dbName string
	if dbName = os.Getenv("MONGODB_DATABASE"); dbName == "" {
		log.Fatal("You must set your 'MONGODB_DATABASE' environment variable.")
	}
	collection := Client.Database(dbName).Collection(collectionName)
	return collection
}
