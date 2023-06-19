package parser

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Parser struct{}

func (p Parser) YamlParse(yamlName string, flow bool) map[string]interface{} {

	projectConfig := config.ProjectConfig()

	// Check yaml files for get all routines in the general flow
	var yamlFile []byte
	var err error

	if flow == true {
		yamlFile, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.yaml", projectConfig.Addresses.Flows, yamlName))
	} else {
		yamlFile, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.yaml", projectConfig.Addresses.Routines, yamlName))
	}

	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	// Unmarshal the YAML into a map[string]interface{}
	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	return data
}
