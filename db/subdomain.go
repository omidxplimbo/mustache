package db

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"time"
)

type Subdomain struct {
	Domain       string    `bson:"domain"`
	Subdomain    string    `bson:"subdomain"`
	CreatedDate  time.Time `bson:"created_date"`
	IP           string    `bson:"ip"`
	CIDR         string    `bson:"cidr"`
	HTTPXService bool      `bson:"httpx_service"`
	CDN          string    `bson:"cdn"`
}

func (s Subdomain) InsertSubdomain(data []Subdomain, projectName string) {

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
	collection := client.Database(dbName).Collection(configProject.Collections[0])
	_, err := collection.InsertMany(ctx, dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Subdomains added to %s project", projectName))

}
