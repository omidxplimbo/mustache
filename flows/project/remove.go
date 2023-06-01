package project

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/db"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
)

func RemoveProject(projectName string) {

	// Connect to database
	client, _ := db.ConnectToDatabase()

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

	if !dbExists {
		fmt.Printf("Project %s dose not exists\n", projectName)
	} else {
		fmt.Print("Are you sure you want to delete the database? With this action all project data will remove from database (y/N): ")
		var confirmation string
		_, _ = fmt.Scanln(&confirmation)
		if confirmation != "y" {
			fmt.Println("Aborting...")
			os.Exit(0)
		}

		// Delete database
		err = client.Database(dbName).Drop(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Project %s deleted successfully\n", projectName)
	}

}
