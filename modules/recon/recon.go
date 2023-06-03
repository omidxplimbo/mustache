package recon

import (
	"github.com/omidxplimbo/mustache/flows/recon"
	"log"
	"os"
)

type Project struct{}

func (p Project) General(projectName string) {
	if projectName == "" {
		log.Fatal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}
	// send request to flow if exists
	recon.General(projectName)
}
