package db

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"time"
)

type Subdomain struct {
	Domain      string    `bson:"domain"`
	Subdomain   string    `bson:"subdomain"`
	CreatedDate time.Time `bson:"created_date"`
	UpdatedDate time.Time `bson:"updated_date"`
	IP          string    `bson:"ip"`
	CIDR        string    `bson:"cidr"`
	Http        bool      `bson:"http"`
	CDN         string    `bson:"cdn"`
}

func (s Subdomain) InsertSubdomain(data []Subdomain, projectName string) {

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	var dataInterface []interface{}
	for _, d := range data {
		dataInterface = append(dataInterface, d)
	}

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configProject.DbTimeout)*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, dataInterface)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Subdomains added to %s project at: %s", projectName, time.Now().Format(configProject.TimeShow)))

}
