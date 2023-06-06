package recon

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func SubDomain(projectName string, target string) {

	projectConfig := config.ProjectConfig()

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile("config/flows/subdomain.yaml")
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	// Unmarshal the YAML into a map[string]interface{}
	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	logger.Info(fmt.Sprintf("Running routin flow %s at: %s", data["name"].(string), time.Now().Format(projectConfig.TimeShow)))

	// Execute commands
	runYaml(data, projectName, target)

	logger.Info(fmt.Sprintf("General flow done at: %s", time.Now().Format(projectConfig.TimeShow)))

}

// Execute commands in the YAML data
func runYaml(data map[string]interface{}, projectName string, domain string) {

	// Set address for tools config
	configPath := config.ProjectConfig().Addresses.Configs
	chaosApi := config.ProjectConfig().Providers.ChaosApi
	dir, _ := os.Getwd()
	configPath = strings.ReplaceAll(configPath, "{{projectPath}}", dir)

	commands, ok := data["command"].([]interface{})
	if !ok {
		return
	}

	var subdomains []db.Subdomain

	for _, cmd := range commands {
		command, ok := cmd.(string)
		if !ok {
			fmt.Println("Invalid command format")
			continue
		}

		// Replace variables in the command
		command = replaceVariables(command, projectName, domain, configPath, chaosApi)

		// Run the command
		logger.Info(fmt.Sprintf("Running Process: " + command))
		output, err := executeCommand(command)
		if err != nil {
			fmt.Printf("Command execution error: %v\n", err)
			continue
		}
		logger.Info(fmt.Sprintf("Run Process Successfully: " + command))

		for _, line := range strings.Split(output, "\n") {
			if line == "" {
				continue
			}
			subdomain := db.Subdomain{
				Domain:      "",
				Subdomain:   strings.TrimSpace(line),
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
				IP:          "",
				CIDR:        "",
				Http:        false,
				CDN:         "",
			}
			subdomains = append(subdomains, subdomain)
		}

	}

	// Add data to database
	if len(subdomains) > 0 {
		subdomainDb := db.Subdomain{}
		subdomainDb.InsertSubdomain(subdomains, projectName)
	}
}

// Replace variables in the command
func replaceVariables(command string, projectName string, domain string, configPath string, chaosApi string) string {
	command = strings.ReplaceAll(command, "$domain$", domain)
	command = strings.ReplaceAll(command, "$projectName$", projectName)
	command = strings.ReplaceAll(command, "$configPath$", configPath)
	command = strings.ReplaceAll(command, "$chaosApi$", chaosApi)
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
