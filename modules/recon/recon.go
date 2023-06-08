package recon

import (
	"github.com/omidxplimbo/mustache/flows/recon"
	"github.com/omidxplimbo/mustache/logger"
	"log"
	"os"
)

type Recon struct{}

type CheckAnySwitch struct {
	WC bool
	PR bool
	PP bool
	PN bool
}

func (r Recon) General(projectName string, target string, wc bool, pr bool, pp bool, pn bool) {
	if projectName == "" {
		log.Fatal("Project flag is required. Please use -h for help.")
		os.Exit(0)
	}
	if target == "" {
		log.Fatal("Domain flag is required. Please use -h for help.")
		os.Exit(0)
	}

	if (pr && pp && pn) || (pr && pp) || (pn && pp) || (pr && pn) {
		logger.Warning("You can use just one of the pr,pn,pp flag. Pleas use -h to show help")
	}
	if !pr && !pp && !pn {
		pn = true
	}

	anySwitch := CheckAnySwitch{
		WC: wc,
		PR: pr,
		PP: pp,
		PN: pn,
	}
	// send request to flow if exists
	recon.General(projectName, target, anySwitch)
}
