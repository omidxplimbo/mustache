package db

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func CheckProject(projectName string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

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
	if !dbExists {
		logger.Warning(fmt.Sprintf("Project %s dose not exist. Please use -h to get help", projectName))
	}
}
