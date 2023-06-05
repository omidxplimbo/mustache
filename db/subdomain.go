package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	var updates []mongo.WriteModel

	for _, subdomain := range data {
		filter := bson.M{"subdomain": subdomain.Subdomain}
		update := bson.M{
			"$set": bson.M{
				"domain":       subdomain.Domain,
				"updated_date": time.Now(),
				"ip":           subdomain.IP,
				"cidr":         subdomain.CIDR,
				"http":         subdomain.Http,
				"cdn":          subdomain.CDN,
			},
			"$setOnInsert": bson.M{"created_date": time.Now()},
		}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		updates = append(updates, updateModel)
	}

	opts := options.BulkWrite().SetOrdered(false)

	_, err := collection.BulkWrite(context.Background(), updates, opts)
	if err != nil {
		logger.Fetal(err.Error())
	}

	logger.Info(fmt.Sprintf("Subdomains added to %s project at: %s", projectName, time.Now().Format(configProject.TimeShow)))

}

func (s Subdomain) GetAllSubdomain(projectName string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	// Access the collection
	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		// Handle error
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	logger.Info(fmt.Sprintf("Subdomains for %s project:", projectName))
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			logger.Fetal(err.Error())
		}
		prettyJSON, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			logger.Fetal(err.Error())
		}
		color.Cyan(string(prettyJSON))
	}

	if err := cursor.Err(); err != nil {
		logger.Fetal(err.Error())
	}
}
