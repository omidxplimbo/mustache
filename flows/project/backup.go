package project

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os/exec"
)

func BackupAllProject(projectName string, path string) {

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := db.ConnectToDatabase()

	// Set up MongoDB export options
	dbName := configProject.DbPrefix + projectName

	// Export collections to JSON file
	filter := bson.M{}
	collections, err := client.Database(dbName).ListCollectionNames(context.Background(), filter)
	if err != nil {
		logger.Fetal(err.Error())
	}
	if len(collections) == 0 {
		logger.Warning(fmt.Sprintf("Project %s dose not exist or dose not has any collection", projectName))
	}
	for _, collection := range collections {
		pathName := path + "/" + collection + ".json"
		cmd := exec.Command("mongoexport", "-d", dbName, "-c", collection, "-o", pathName)
		err := cmd.Run()
		if err != nil {
			logger.Fetal(err.Error())
		}
		logger.Info(fmt.Sprintf("Collection %s in project '%s' exported to "+pathName+" successfully", collection, projectName))
	}
}

func BackupCollectionProject(projectName string, collection string, path string) {
	// Get project config
	configProject := config.ProjectConfig()

	dbName := configProject.DbPrefix + projectName
	cmd := exec.Command("mongoexport", "-d", dbName, "-c", collection, "-o", path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info(fmt.Sprintf("Collection %s in project '%s' exported to "+path+" successfully", collection, projectName))
}
