package db

import (
	"context"
	"github.com/omidxplimbo/mustache/config"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDatabase() (*mongo.Client, context.Context) {

	// Get project config
	configProject := config.ProjectConfig()

	// Replace the following with your MongoDB connection string
	mongoURI := configProject.MongoUrl

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configProject.DbTimeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Connect to database with error. See logs")
	}

	return client, ctx
}
