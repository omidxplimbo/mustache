package module

import (
	"github.com/omidxplimbo/mustache/flows/watch"
	"github.com/omidxplimbo/mustache/logger"
	"os"
)

type Watch struct{}

func (w Watch) Subdomain(projectName string) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.GetAllSubdomain(projectName)
}
