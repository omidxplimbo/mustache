package db

import (
	"context"
	"fmt"
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

func (s Domain) InsertDomain(data []Domain, projectName string) {
	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName

	// Insert the array of data into the MongoDB collection
	collection := client.Database(dbName).Collection(configProject.Collections[0])

	var updates []mongo.WriteModel

	for _, domain := range data {
		filter := bson.M{"domain": domain.Domain}
		update := bson.M{"$set": bson.M{}}

		if domain.Domain != "" {
			update["$set"].(bson.M)["domain"] = domain.Domain
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
