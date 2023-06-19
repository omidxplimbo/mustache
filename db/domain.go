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

type Domain struct {
	Domain      string    `bson:"domain"`
	CreatedDate time.Time `bson:"created_date"`
	UpdatedDate time.Time `bson:"updated_date"`
}

func (d Domain) InsertDomain(data []Domain, projectName string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[1])

	var updates []mongo.WriteModel

	for _, domain := range data {
		filter := bson.M{"domain": domain.Domain}
		update := bson.M{"$set": bson.M{}}

		if domain.Domain != "" {
			update["$set"].(bson.M)["domain"] = domain.Domain
			update["$set"].(bson.M)["updated_date"] = time.Now()
		}
		update["$setOnInsert"] = bson.M{"created_date": time.Now()}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		updates = append(updates, updateModel)
	}

	opts := options.BulkWrite().SetOrdered(false)

	_, err := collection.BulkWrite(context.Background(), updates, opts)
	if err != nil {
		logger.Fetal(err.Error())
	}

	logger.Info(fmt.Sprintf("Domains added to %s project at: %s", projectName, time.Now().Format(configProject.TimeShow)))
}

func (d Domain) GetAllDomains(projectName string, justCount *bool) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[1])

	// Access the collection
	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		logger.Fetal(err.Error())
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	count, _ := client.Database(dbName).Collection(configProject.Collections[1]).CountDocuments(context.Background(), bson.M{})
	logger.Info(fmt.Sprintf("Domains for %s project:", projectName))
	logger.Info(fmt.Sprintf("Count of domain for %s project is: %d", projectName, count))

	if !*justCount {
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
			color.Yellow(string(prettyJSON))
		}
		if err := cursor.Err(); err != nil {
			logger.Fetal(err.Error())
		}
	}

}

func (d Domain) GetDomain(projectName string, target string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[1])

	filter := bson.M{"domain": target}

	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		logger.Fetal(err.Error())
	}
	defer cursor.Close(context.Background())

	// Check if cursor has any documents

	logger.Info(fmt.Sprintf("Get information of %s root domain for %s project: ", target, projectName))
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
		color.Yellow(string(prettyJSON))
	}

	if err := cursor.Err(); err != nil {
		logger.Fetal(err.Error())
	}
}
