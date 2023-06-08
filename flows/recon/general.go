package recon

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	recon "github.com/omidxplimbo/mustache/flows/recon/routines"
	"github.com/omidxplimbo/mustache/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

func General(projectName string, target string, wc bool, pr bool, pp bool, pn bool) {

	projectConfig := config.ProjectConfig()
	logger.Info(fmt.Sprintf("Running general flow at: %s", time.Now().Format(projectConfig.TimeShow)))

	// Check yaml files for get all routines in the general flow
	generalFlow, err := ioutil.ReadFile("config/base-routine/general.yaml")
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	// Unmarshal the YAML into a map[string]interface{}
	var data map[string]interface{}
	err = yaml.Unmarshal(generalFlow, &data)
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	logger.Info(fmt.Sprintf("Running routin flow %s at: %s", data["name"].(string), time.Now().Format(projectConfig.TimeShow)))

	for _, routine := range data["routine"].([]interface{}) {
		switch routine {
		case "subdomain":
			recon.SubDomain(projectName, target, pr, pn, pp, wc)
		default:
			logger.Warning(fmt.Sprintf("There isn't any routins in the general flow"))
		}
	}

	logger.Info(fmt.Sprintf("General flow done at: %s", time.Now().Format(projectConfig.TimeShow)))
}
