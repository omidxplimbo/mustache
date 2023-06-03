package recon

import (
	"github.com/omidxplimbo/mustache/flows/recon"
	"log"
	"os"
)

type Recon struct{}

func (r Recon) General(projectName string, target string) {
	if projectName == "" {
		log.Fatal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}
	if target == "" {
		log.Fatal("Domain flag is required. Please use -h for help.")
		os.Exit(0)
	}

	// send request to flow if exists
	recon.General(projectName, target)
}
