package routines

import (
	"fmt"
	"github.com/omidxplimbo/mustache/flows/routines"
	"github.com/omidxplimbo/mustache/logger"
)

type RoutineRunner struct {
}

func (r RoutineRunner) ExecuteRoutines(data map[string]interface{}, args ...interface{}) {
	routineMap := map[string]func(...interface{}){
		"subdomain": routines.SubDomain,
		"ip":        routines.Ip,
	}

	routinesGet, ok := data["routine"].([]interface{})
	if !ok {
		logger.Fetal("Failed to parse routines from data")
		return
	}

	for _, routine := range routinesGet {
		routineName, ok := routine.(string)
		if !ok {
			logger.Warning("Failed to parse routine name")
			continue
		}

		if fn, ok := routineMap[routineName]; ok {
			fn(args...)
		} else {
			logger.Warning(fmt.Sprintf("Routine '%s' not found in the general flow", routineName))
		}
	}
}
