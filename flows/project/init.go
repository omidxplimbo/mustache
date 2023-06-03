package project

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func InitiateProject(projectName string) {

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := db.ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName
	filter := bson.D{}
	databases, err := client.ListDatabaseNames(context.Background(), filter)
	if err != nil {
		log.Fatal("Get list of database with error")
	}
	dbExists := false
	for _, database := range databases {
		if database == dbName {
			dbExists = true
			break
		}
	}

	// Create database and collections if it doesn't exist
	if !dbExists {
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[0])
		if err != nil {
			log.Fatal(err)
		}

		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[1])
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[2])
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[3])
		if err != nil {
			log.Fatal(err)
		}

		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[4])
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[5])
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[6])
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), configProject.Collections[7])
		if err != nil {
			log.Fatal(err)
		}
		logger.Info(fmt.Sprintf("Project %s created successfully\n", projectName))
	} else {
		logger.Warning(fmt.Sprintf("Project %s already exists\n", projectName))
	}
}
