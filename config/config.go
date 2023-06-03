package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	AppName     string   `json:"appName"`
	DbPrefix    string   `json:"dbPrefix"`
	DbTimeout   int      `json:"dbTimeout"`
	MongoUrl    string   `json:"mongoUrl"`
	TimeShow    string   `json:"timeShow"`
	Collections []string `json:"collections"`
	PublicError struct {
		ErrorConnectDB   string `json:"errorConnectDb"`
		SubdomainAddToDB string `json:"subdomainAddToDatabase"`
	} `json:"publicErrors"`
}

func ProjectConfig() Config {
	// Read the JSON file
	filePath, _ := os.Getwd()
	jsonData, err := ioutil.ReadFile(filePath + "/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON data
	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
