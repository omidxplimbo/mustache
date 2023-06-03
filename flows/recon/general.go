package recon

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	databaseName, ok := data["database"].(string)
	if !ok {
		fmt.Println("Invalid database name")
		return
	}
	collectionName, ok := data["collection"].(string)
	if !ok {
		fmt.Println("Invalid collection name")
		return
	}

	// Create MongoDB session and collection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("MongoDB client error: %v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("MongoDB connection error: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)
	collection := client.Database(databaseName).Collection(collectionName)

	var subdomains []subdomain // Declare subdomains variable outside of loop

	commands, ok := data["command"].([]interface{})
	if !ok {
		fmt.Println("No commands found")
		return
	}

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
	}

	// Insert subdomains into MongoDB collection
	for _, subdomain := range subdomains {
		_, err := collection.InsertOne(ctx, subdomain)
		if err != nil {
			logger.Fetal(fmt.Sprintf("MongoDB insert error: %v\n", err))
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
