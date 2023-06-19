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

func (w Watch) LatestSubdomain(projectName string, count int) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	if count == 0 {
		count = 10
	}

	// send request to flow if exists
	watch.LatestSubdomain(projectName, count)
}

func (w Watch) LivesSubdomain(projectName string, justCount *bool) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.LivesSubdomain(projectName, justCount)
}

func (w Watch) LatestLivesSubdomain(projectName string, count int) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	if count == 0 {
		count = 10
	}

	// send request to flow if exists
	watch.LatestLivesSubdomain(projectName, count)
}

func (w Watch) GetSub(projectName string, target string) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	if target == "" {
		logger.Fetal("Target flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.GetSub(projectName, target)
}

func (w Watch) Resolved(projectName string, justCount *bool) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.ResolvedSubdomain(projectName, justCount)
}

func (w Watch) Domain(projectName string, justCount *bool) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.Domain(projectName, justCount)
}

func (w Watch) GetDomain(projectName string, target string) {
	if projectName == "" {
		logger.Fetal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}

	if target == "" {
		logger.Fetal("Target flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	watch.GetDomain(projectName, target)
}
