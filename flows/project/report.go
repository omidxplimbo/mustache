package project

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
)

func GetReport(projectName string) {
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
		logger.Warning(fmt.Sprintf("Project %s dose not exist or dose not has any collection\n", projectName))
	}

	// Print the top border
	color.HiYellow("+" + "------------------------------------------------------------" + "+")

	// Print the collection counts

	for _, collection := range collections {
		count, _ := client.Database(dbName).Collection(collection).CountDocuments(context.Background(), filter)
		color.HiYellow("| Collection %-10s in project %-10s has %5d items |\n", collection, projectName, count)
	}

	// Print the bottom border
	color.HiYellow("+" + "------------------------------------------------------------" + "+")
}
