package db

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"time"
)

type Url struct {
	Domain      string    `bson:"domain"`
	Subdomain   string    `bson:"subdomain"`
	Url         string    `bson:"url"`
	CreatedDate time.Time `bson:"created_date"`
	Live        bool      `bson:"live"`
}

func (s Subdomain) InsertUrl(data []Url, projectName string) {

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
	collection := client.Database(dbName).Collection(configProject.Collections[4])
	_, err := collection.InsertMany(ctx, dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Urls added to %s project", projectName))

}
