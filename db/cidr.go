package db

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"time"
)

type Cidr struct {
	Cidr        string    `bson:"cidr"`
	CreatedDate time.Time `bson:"created_date"`
	Company     string    `bson:"company"`
	Description string    `bson:"description"`
}

func (d Domain) InsertCidr(data []Cidr, projectName string) {

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
	collection := client.Database(dbName).Collection(configProject.Collections[2])
	_, err := collection.InsertMany(ctx, dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Cidrs added to %s project", projectName))

}
