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
		update := bson.M{"$set": bson.M{}}

		if subdomain.Domain != "" {
			update["$set"].(bson.M)["domain"] = subdomain.Domain
		}
		if subdomain.Http != s.Http {
			update["$set"].(bson.M)["http"] = subdomain.Http
		}
		if subdomain.CDN != s.CDN {
			update["$set"].(bson.M)["cdn"] = subdomain.CDN
		}
		if subdomain.IP != "" {
			update["$set"].(bson.M)["ip"] = subdomain.IP
		}
		if subdomain.CIDR != "" {
			update["$set"].(bson.M)["cidr"] = subdomain.CIDR
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
	count, _ := client.Database(dbName).Collection(configProject.Collections[0]).CountDocuments(context.Background(), bson.M{})
	logger.Info(fmt.Sprintf("Subdomains for %s project:", projectName))
	logger.Info(fmt.Sprintf("Count of subdomain for %s project is: %d", projectName, count))
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

func (s Subdomain) LatestSubdomain(projectName string, count int) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	optionsGet := options.Find()
	optionsGet.SetSort(bson.D{{"_id", -1}})
	optionsGet.SetLimit(int64(count))

	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), bson.D{}, optionsGet)
	if err != nil {
		// Handle error
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	countItem, _ := client.Database(dbName).Collection(configProject.Collections[0]).CountDocuments(context.Background(), bson.M{})
	logger.Info(fmt.Sprintf("%d Latest subdomains for %s project", count, projectName))
	logger.Info(fmt.Sprintf("Count of subdomain for %s project is: %d", projectName, countItem))
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

func (s Subdomain) LivesSubdomain(projectName string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	filter := bson.M{"http": true}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		logger.Fetal(err.Error())
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	countItem, _ := client.Database(dbName).Collection(configProject.Collections[0]).CountDocuments(context.Background(), filter)
	logger.Info(fmt.Sprintf("All Live Subdomain for %s", projectName))
	logger.Info(fmt.Sprintf("Count of live subdomains for %s project is: %d", projectName, countItem))
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

func (s Subdomain) LatestLivesSubdomain(projectName string, count int) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	optionsGet := options.Find()
	optionsGet.SetSort(bson.D{{"_id", -1}})
	optionsGet.SetLimit(int64(count))
	filter := bson.M{"http": true}

	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), filter, optionsGet)
	if err != nil {
		// Handle error
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	countItem, _ := client.Database(dbName).Collection(configProject.Collections[0]).CountDocuments(context.Background(), filter)
	logger.Info(fmt.Sprintf("%d Latest live subdomains for %s project", count, projectName))
	logger.Info(fmt.Sprintf("Count of latest live subdomains for %s project is: %d", projectName, countItem))
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

func (s Subdomain) GetSub(projectName string, target string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	filter := bson.M{"subdomain": target}

	// Retrieve all documents from the collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		logger.Fetal(err.Error())
	}
	defer cursor.Close(context.Background())

	// Check if cursor has any documents

	logger.Info(fmt.Sprintf("Get information of %s subdomains for %s project: ", target, projectName))
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

func (s Subdomain) ResolvedSubdomain(projectName string, justCount *bool) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	filter := bson.M{"ip": bson.M{"$exists": true, "$ne": ""}}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		logger.Fetal(err.Error())
	}

	defer cursor.Close(context.Background())

	// Iterate through the cursor and print the documents
	countItem, _ := client.Database(dbName).Collection(configProject.Collections[0]).CountDocuments(context.Background(), filter)
	logger.Info(fmt.Sprintf("All Resolved Subdomain for %s", projectName))
	logger.Info(fmt.Sprintf("Count of resolved subdomains for %s project is: %d", projectName, countItem))
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
