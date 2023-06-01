package project

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func InitiateProject(projectName string) {

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client := db.ConnectToDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configProject.DbTimeout)*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		log.Fatal("Connect to database with error. See logs")
	}

	// Check if database exists
	dbName := "mustache_" + projectName
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
		err = client.Database(dbName).CreateCollection(context.Background(), "domain")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "cidr")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "asn")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "subdomain")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "urls")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "parameters")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Database(dbName).CreateCollection(context.Background(), "wordlist")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Project %s created successfully\n", projectName)
	} else {
		fmt.Printf("Project %s already exists\n", projectName)
	}
}
