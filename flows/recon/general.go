package recon

import (
	"fmt"
	"github.com/omidxplimbo/mustache/base/parser"
	"github.com/omidxplimbo/mustache/base/runner/routines"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"os/exec"
	"time"
)

func General(projectName string, target string, wc bool, pr bool, pp bool, pn bool, aa bool) {

	projectConfig := config.ProjectConfig()

	// Check project exist
	db.CheckProject(projectName)

	logger.Info(fmt.Sprintf("Running general flow at: %s", time.Now().Format(projectConfig.TimeShow)))
	logger.AddLog(fmt.Sprintf("Running general flow at: %s", time.Now().Format(projectConfig.TimeShow)))

	// parse yaml files
	parserObj := parser.Parser{}
	runnerObj := routines.RoutineRunner{}

	data := parserObj.YamlParse("general", false)

	// execute routine
	runnerObj.ExecuteRoutines(data, projectName, target, wc, pr, pp, pn, aa)

	// remove data
	cmd := exec.Command("bash", "-c", fmt.Sprintf("rm %s*", projectName))
	_ = cmd.Run()
	logger.Info(fmt.Sprintf("General flow done at: %s", time.Now().Format(projectConfig.TimeShow)))
	logger.AddLog(fmt.Sprintf("General flow done at: %s", time.Now().Format(projectConfig.TimeShow)))
}
