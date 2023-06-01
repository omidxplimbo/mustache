package db

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDatabase() *mongo.Client {

	// Replace the following with your MongoDB connection string
	mongoURI := "mongodb://localhost:27017"

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	return client
}
