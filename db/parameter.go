package db

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"time"
)

type Parameter struct {
	Subdomain   string    `bson:"subdomain"`
	CreatedDate time.Time `bson:"created_date"`
	Url         string    `bson:"url"`
	Parameter   string    `bson:"parameter"`
	Live        bool      `bson:"live"`
}

func (s Subdomain) InsertParameter(data []Parameter, projectName string) {

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, ctx := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	var dataInterface []interface{}
	for _, d := range data {
		dataInterface = append(dataInterface, d)
	}

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[5])
	_, err := collection.InsertMany(ctx, dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Parameters added to %s project", projectName))

}
