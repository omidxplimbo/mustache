package project

import (
	"github.com/omidxplimbo/mustache/flows/project"
	"log"
	"os"
)

type Project struct{}

func (p Project) Init(projectName string) {
	if projectName == "" {
		log.Fatal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}
	// send request to flow if exists
	project.InitiateProject(projectName)
}

func (p Project) Report(ProjectName string) {

	if ProjectName == "" {
		log.Fatal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}
	// send request to flow if exists
	project.GetReport(ProjectName)
}

func (p Project) Backup(projectName string, collection string, path string) {

	if path == "" {
		log.Fatal("Path flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// check collection
	if collection == "all" {
		project.BackupAllProject(projectName, path)
	} else {
		project.BackupCollectionProject(projectName, collection, path)
	}

}

func (p Project) Remove(projectName string) {
	// send request to flow if exists
	project.RemoveProject(projectName)
}
