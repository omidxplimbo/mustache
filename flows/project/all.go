package project

import (
	"context"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"strings"
	"time"
)

func AllProject() {
	// Get project config
	configProject := config.ProjectConfig()

	client, _ := db.ConnectToDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configProject.DbTimeout)*time.Second)
	defer cancel()

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Filter databases with names starting with "mustache_*"
	if len(databases) > 0 {
		for _, dbName := range databases {
			if strings.HasPrefix(dbName, configProject.DbPrefix) {
				logger.Info(dbName)
			}
		}
		os.Exit(0)
	}

	logger.Warning("There isn't any project.")

}
