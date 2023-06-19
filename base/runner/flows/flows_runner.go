package flows

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

type RunnerFlow struct {
}

type AllSwitches struct {
	Target  string
	Project string
	WC      interface{}
	PR      interface{}
	PP      interface{}
	PN      interface{}
	NWC     interface{}
	AA      interface{}
	WCR     interface{}
}

type CommandParams struct {
	Project      string
	Target       string
	ConfigPath   string
	ResolverPath string
	WordListPath string
	ChaosApi     string
	BinaryConfig string
}

func (f RunnerFlow) ExecuteFlows(data map[string]interface{}, projectName string, target string, flagMap map[string]interface{}) []interface{} {

	// Set Flags
	allSwitches := setFlags(projectName, target, flagMap)
	allParams := setParams(projectName, target)

	// run flow
	return runFlow(data, allParams, allSwitches)

}

func runFlow(data map[string]interface{}, params CommandParams, switches AllSwitches) []interface{} {

	var postRun []interface{}
	projectConfig := config.ProjectConfig()

	logger.Info(fmt.Sprintf("Running %s routine at: %s", data["name"], time.Now().Format(projectConfig.TimeShow)))
	logger.Process(fmt.Sprintf("%s", data["description"]))

	steps := data["steps"].([]interface{})

	for _, step := range steps {
		stepMap := step.(map[string]interface{})

		for key, item := range stepMap {
			switch key {
			case "scripts":
				for _, script := range item.([]interface{}) {
					runCommand(script.(string), params)
				}
			case "switches":
				for _, sw := range item.([]interface{}) {
					switchMap := sw.(map[string]interface{})
					for switchKey, switchValue := range switchMap {
						fieldValue := reflect.ValueOf(switches).FieldByName(switchKey)
						if fieldValue.IsValid() {
							underlyingValue := fieldValue.Interface()
							if underlyingValue.(bool) == true {
								runCommand(switchValue.(string), params)
							}
						}
					}
				}
			case "postRun":
				for _, postRunData := range item.([]interface{}) {
					postRun = append(postRun, replaceVariables(postRunData.(string), params))
				}
			default:
				logger.Warning(fmt.Sprintf("Can't find any valid flow in yaml config."))

			}

		}
	}

	return postRun
}

func setFlags(projectName string, target string, flagsMap map[string]interface{}) AllSwitches {

	flags := AllSwitches{
		Target:  target,
		Project: projectName,
	}

	for key, value := range flagsMap {
		switch key {
		case "WC":
			if v, ok := value.(bool); ok {
				flags.WC = v
			}
		case "PR":
			if v, ok := value.(bool); ok {
				flags.PR = v
			}
		case "PP":
			if v, ok := value.(bool); ok {
				flags.PP = v
			}
		case "PN":
			if v, ok := value.(bool); ok {
				flags.PN = v
			}
		case "AA":
			if v, ok := value.(bool); ok {
				flags.AA = v
			}
		case "WCR":
			if v, ok := value.(bool); ok {
				flags.WCR = v
			}
		}

	}

	// Set NWC based on WC value
	if flags.WC.(bool) {
		flags.NWC = false
	} else {
		flags.NWC = true
	}

	return flags
}

func setParams(projectName string, target string) CommandParams {

	dir, _ := os.Getwd()
	params := CommandParams{
		Target:       target,
		Project:      projectName,
		ConfigPath:   strings.ReplaceAll(config.ProjectConfig().Addresses.Configs, "{{projectPath}}", dir),
		ChaosApi:     config.ProjectConfig().Providers.ChaosApi,
		WordListPath: strings.ReplaceAll(config.ProjectConfig().Addresses.Wordlist, "{{projectPath}}", dir),
		ResolverPath: strings.ReplaceAll(config.ProjectConfig().Addresses.Resolvers, "{{projectPath}}", dir),
		BinaryConfig: strings.ReplaceAll(config.ProjectConfig().Addresses.BinaryConfig, "{{projectPath}}", dir),
	}

	return params
}

func runCommand(command string, params CommandParams) {

	command = replaceVariables(command, params)
	logger.Process(fmt.Sprintf("Running Process "))
	logger.AddLog(fmt.Sprintf("Running Process: " + command))
	_, err := executeCommand(command)
	if err != nil {
		logger.WarningC(fmt.Sprintf("Command execution error: %v\n", err))
	}
	logger.Info(fmt.Sprintf("Run Process Successfull"))
}

func replaceVariables(command string, params CommandParams) string {
	command = strings.ReplaceAll(command, "$domain$", params.Target)
	command = strings.ReplaceAll(command, "$projectName$", params.Project)
	command = strings.ReplaceAll(command, "$configPath$", params.BinaryConfig)
	command = strings.ReplaceAll(command, "$chaosApi$", params.ChaosApi)
	command = strings.ReplaceAll(command, "$resolver$", params.ResolverPath)
	command = strings.ReplaceAll(command, "$wordlist$", params.WordListPath)
	command = strings.ReplaceAll(command, "$excludeDomain$", params.ConfigPath)
	return command
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
