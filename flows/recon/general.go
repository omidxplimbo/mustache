package recon

import (
	"context"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
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

	logger.Info(fmt.Sprintf("Running general flow at: %s", time.DateTime))

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

		logger.Info(fmt.Sprintf("Running routin flow %s at: %s", data["name"].(string), time.DateTime))

		// Execute commands
		runYaml(data, projectName, target)

		logger.Info(fmt.Sprintf("General flow done at: %s", time.DateTime))
	}
}

// Execute commands in the YAML data
func runYaml(data map[string]interface{}, projectName string, domain string) {

	// Extract database and collection names from YAML data
	collectionName := data["collection"].(string)

	// Get project config
	configProject := config.ProjectConfig()

	// Connect to database
	client, _ := db.ConnectToDatabase()

	// Check if database exists
	dbName := configProject.DbPrefix + projectName
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configProject.DbTimeout)*time.Second)
	defer cancel()

	// Get database
	collection := client.Database(dbName).Collection(collectionName)

	var subdomains []subdomain // Declare subdomains variable outside of loop

	commands, ok := data["command"].([]interface{})
	if !ok {
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
		logger.Info(fmt.Sprintf("Running Process: " + command))
		output, err := executeCommand(command)
		if err != nil {
			fmt.Printf("Command execution error: %v\n", err)
			continue
		}
		logger.Info(fmt.Sprintf("Run Process Successfully: " + command))

		switch data["name"].(string) {
		case "subdomain":
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
		default:
			return
		}
	}

	if len(subdomains) > 0 {
		// Insert subdomains into MongoDB collection
		for _, subdomain := range subdomains {
			_, err := collection.InsertOne(ctx, subdomain)
			if err != nil {
				logger.Fetal(fmt.Sprintf("MongoDB insert error: %v\n", err))
			}
		}
	}
	logger.Info("Add data to database done.")

}

// Replace variables in the command
func replaceVariables(command string, projectName string, domain string) string {
	command = strings.ReplaceAll(command, "$domain$", domain)
	command = strings.ReplaceAll(command, "$projectName$", projectName)
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
