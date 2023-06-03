package recon

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type subdomain struct {
	Subdomain string    `bson:"subdomain"`
	Created   time.Time `bson:"created_at"`
	Updated   time.Time `bson:"updated_at"`
}

func General(projectName string, target string) {
	// List YAML files in a directory
	yamlFiles, err := filepath.Glob("flows/recon/configs/general/*.yaml")
	if err != nil || len(yamlFiles) == 0 {
		logger.Fetal("Yaml Configuration not exist")
	}

	// Parse each YAML file
	for _, file := range yamlFiles {
		// Read the YAML file
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			logger.Fetal("Yaml Configuration not exist")
		}

		// Unmarshal the YAML into a map[string]interface{}
		var data map[string]interface{}
		err = yaml.Unmarshal(yamlFile, &data)
		if err != nil {
			logger.Fetal("Yaml Configuration not exist")
		}
		// Execute commands
		runCommands(data, projectName, target)

		fmt.Println("-------------------")
	}
}

// Execute commands in the YAML data
func runCommands(data map[string]interface{}, projectName string, domain string) {

	// Extract database and collection names from YAML data
	databaseName := data["database"].(string)
	collectionName := data["collection"].(string)
	checkDatabase := false
	var collection *mongo.Collection
	var ctx context.Context

	if databaseName != "" && collectionName != "" {
		checkDatabase = true
		// Create MongoDB session and collection
		client, ctxData := db.ConnectToDatabase()
		ctx = ctxData
		collection = client.Database(databaseName).Collection(collectionName)
	}

	commands := data["command"].([]interface{})

	for _, cmd := range commands {
		command, ok := cmd.(string)
		if !ok {
			fmt.Println("Invalid command format")
			continue
		}

		// Replace variables in the command
		command = replaceVariables(command, projectName, domain)

		// Run the command
		output, err := executeCommand(command)
		if err != nil {
			fmt.Printf("Command execution error: %v\n", err)
			continue
		}

		// Extract subdomains from command output
		var subdomains []subdomain
		for _, line := range strings.Split(output, "\n") {
			if line == "" {
				continue
			}
			subdomain := subdomain{
				Subdomain: strings.TrimSpace(line),
				Created:   time.Now(),
				Updated:   time.Now(),
			}
			subdomains = append(subdomains, subdomain)
		}

		if checkDatabase {
			// Insert subdomains into MongoDB collection
			for _, subdomain := range subdomains {
				_, err := collection.InsertOne(ctx, subdomain)
				if err != nil {
					logger.Fetal(fmt.Sprintf("MongoDB insert error: %v\n", err))
				}
			}
		}

	}
}

// Replace variables in the command
func replaceVariables(command string, projectName string, domain string) string {

	projectConfig := config.ProjectConfig()
	command = strings.ReplaceAll(command, "$domain$", domain)
	command = strings.ReplaceAll(command, "$projectName$", projectName)
	command = strings.ReplaceAll(command, "$database$", projectConfig.DbPrefix+projectName)

	return command
}

// Execute a command and return its output
func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
