package recon

import (
	"bufio"
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CommandParams struct {
	ProjectName  string
	Domain       string
	ConfigPath   string
	ChaosApi     string
	WordListPath string
	ResolverPath string
}

type ResolvesSubdomain struct {
	Subdomain string
	Ip        string
}

type CheckAnySwitch struct {
	WC bool
	PR bool
	PP bool
	PN bool
}

type LiveHosts struct {
	Subdomain string
	CDN       string
	Http      bool
}

func SubDomain(projectName string, target string, pr bool, pn bool, pp bool, wc bool) {
	projectConfig := config.ProjectConfig()

	yamlFile, err := ioutil.ReadFile("config/flows/subdomain.yaml")
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		logger.Fetal("Yaml Configuration not exist")
	}

	logger.Info(fmt.Sprintf("Running routin flow %s at: %s", data["name"].(string), time.Now().Format(projectConfig.TimeShow)))

	// set switch

	anySwitch := CheckAnySwitch{
		WC: wc,
		PR: pr,
		PP: pp,
		PN: pn,
	}

	runYaml(data, projectName, target, anySwitch)

	logger.Info(fmt.Sprintf("General flow done at: %s", time.Now().Format(projectConfig.TimeShow)))
}

func runYaml(data map[string]interface{}, projectName string, domain string, anySwitch CheckAnySwitch) {

	configPath := config.ProjectConfig().Addresses.Configs
	chaosApi := config.ProjectConfig().Providers.ChaosApi
	wordlistPath := config.ProjectConfig().Addresses.Wordlist
	resolverPath := config.ProjectConfig().Addresses.Resolvers
	dir, _ := os.Getwd()
	configPath = strings.ReplaceAll(configPath, "{{projectPath}}", dir)
	wordlistPath = strings.ReplaceAll(wordlistPath, "{{projectPath}}", dir)
	resolverPath = strings.ReplaceAll(resolverPath, "{{projectPath}}", dir)

	params := CommandParams{
		ProjectName:  projectName,
		Domain:       domain,
		ConfigPath:   configPath,
		ChaosApi:     chaosApi,
		WordListPath: wordlistPath,
		ResolverPath: resolverPath,
	}

	subdomains, resolveSubdomain, liveSubdomains, err := executeCommands(data, params, anySwitch, projectName)
	if err != nil {
		logger.Warning(err.Error())
		return
	}

	if len(subdomains) > 0 {
		subdomainDb := db.Subdomain{}
		subdomainDb.InsertSubdomain(subdomains, projectName)
	}
	if len(resolveSubdomain) > 0 {
		resolveSubdomainDb := db.Subdomain{}
		resolveSubdomainDb.InsertSubdomain(resolveSubdomain, projectName)
	}
	if len(liveSubdomains) > 0 {
		liveSubdomainsDb := db.Subdomain{}
		liveSubdomainsDb.InsertSubdomain(liveSubdomains, projectName)
	}
}

func executeCommands(data map[string]interface{}, params CommandParams, anySwitch CheckAnySwitch, projectName string) ([]db.Subdomain, []db.Subdomain, []db.Subdomain, error) {

	// Set Commands and init data
	var subdomains []db.Subdomain
	var rs []db.Subdomain
	var lives []db.Subdomain
	var passiveSubdomain []string
	var resolves []ResolvesSubdomain
	var liveHosts []LiveHosts

	cp, _ := data["Commands-passive"].([]interface{})
	cw, _ := data["Commands-wc"].([]interface{})
	cr, _ := data["Commands-resolver"].([]interface{})
	cb, _ := data["Commands-bruteForce"].([]interface{})
	cpr, _ := data["Commands-pr"].([]interface{})
	cpp, _ := data["Commands-pp"].([]interface{})
	cpn, _ := data["Commands-pn"].([]interface{})
	cl, _ := data["Commands-lives"].([]interface{})
	cf, _ := data["Commands-final"].([]interface{})

	// run commands
	filePassiveStatus := fmt.Sprintf("%s-passive.txt", projectName)
	fileResolveStatus := fmt.Sprintf("%s-final.txt", projectName)
	fileLivesStatus := fmt.Sprintf("%s-lives.txt", projectName)
	if !fileExists(fileResolveStatus) && !fileExists(filePassiveStatus) && !fileExists(fileLivesStatus) {
		runCommand(cp, params)
		if anySwitch.WC {
			runCommand(cw, params)
		} else {
			runCommand(cr, params)
		}
		runCommand(cb, params)
		if anySwitch.PR {
			runCommand(cpr, params)
		}
		if anySwitch.PP {
			runCommand(cpp, params)
		}
		if anySwitch.PN {
			runCommand(cpn, params)
		}
		runCommand(cl, params)
	}

	//set passive and resolves
	passiveSubdomain = setPassiveSubdomain(projectName)
	resolves = setResolveSubdomain(projectName)
	liveHosts = setLiveHosts(projectName)

	// save passive domain into database
	if len(passiveSubdomain) > 0 {
		for _, line := range passiveSubdomain {
			if line == "" {
				continue
			}
			subdomain := db.Subdomain{
				Domain:      "",
				Subdomain:   strings.TrimSpace(line),
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
				IP:          "",
				CIDR:        "",
				Http:        false,
				CDN:         "",
			}
			subdomains = append(subdomains, subdomain)
		}
	}
	if len(resolves) > 0 {
		for _, line := range resolves {
			resolves := db.Subdomain{
				Domain:      projectName,
				Subdomain:   strings.TrimSpace(line.Subdomain),
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
				IP:          line.Ip,
				CIDR:        "",
				Http:        false,
				CDN:         "",
			}
			rs = append(rs, resolves)
		}
	}
	if len(liveHosts) > 0 {
		for _, line := range liveHosts {
			resolves := db.Subdomain{
				Domain:      projectName,
				Subdomain:   strings.TrimSpace(line.Subdomain),
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
				Http:        line.Http,
				CDN:         line.CDN,
			}
			lives = append(lives, resolves)
		}
	}

	// remove temp file
	runCommand(cf, params)
	return subdomains, rs, lives, nil
}

func setLiveHosts(projectName string) []LiveHosts {
	var liveHosts []LiveHosts

	filePath := fmt.Sprintf("%s-lives.txt", projectName)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fetal(fmt.Sprintf("Failed to open file: %v\n", err))
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		// Ensure there are exactly 2 parts (subdomain and IP)
		if len(parts) > 1 {
			subdomain := strings.Trim(parts[0], "https://")
			subdomain = strings.Trim(subdomain, "http://")
			CDN := strings.Trim(parts[1], "[]")

			liveHost := LiveHosts{
				Subdomain: subdomain,
				Http:      true,
				CDN:       CDN,
			}

			liveHosts = append(liveHosts, liveHost)
		} else {
			subdomain := parts[0]
			liveHost := LiveHosts{
				Subdomain: subdomain,
				Http:      true,
			}

			liveHosts = append(liveHosts, liveHost)
		}
	}

	return liveHosts
}

func setPassiveSubdomain(projectName string) []string {

	filePath := fmt.Sprintf("%s-passive.txt", projectName)
	var passive []string

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fetal(fmt.Sprintf("Failed to open file: %v\n", err))
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		passive = append(passive, line)
	}

	return passive
}

func setResolveSubdomain(projectName string) []ResolvesSubdomain {

	var resolves []ResolvesSubdomain

	filePath := fmt.Sprintf("%s-final.txt", projectName)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fetal(fmt.Sprintf("Failed to open file: %v\n", err))
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		// Ensure there are exactly 2 parts (subdomain and IP)
		if len(parts) > 2 {
			subdomain := parts[0]
			ip := strings.Trim(parts[1], "[]")

			resolve := ResolvesSubdomain{
				Subdomain: subdomain,
				Ip:        ip,
			}

			resolves = append(resolves, resolve)
		} else {
			subdomain := parts[0]
			resolve := ResolvesSubdomain{
				Subdomain: subdomain,
				Ip:        "",
			}

			resolves = append(resolves, resolve)
		}
	}

	return resolves
}

func runCommand(cp []interface{}, params CommandParams) {
	for _, cmd := range cp {
		command, ok := cmd.(string)
		if !ok {
			logger.Fetal("invalid command format")
		}

		command = replaceVariables(command, params)

		logger.Info(fmt.Sprintf("Running Process: " + command))
		_, err := executeCommand(command)
		if err != nil {
			logger.WarningC(fmt.Sprintf("Command execution error: %v\n", err))
			continue
		}
		logger.Info(fmt.Sprintf("Run Process Successfully: " + command))
	}
}

func replaceVariables(command string, params CommandParams) string {
	command = strings.ReplaceAll(command, "$domain$", params.Domain)
	command = strings.ReplaceAll(command, "$projectName$", params.ProjectName)
	command = strings.ReplaceAll(command, "$configPath$", params.ConfigPath)
	command = strings.ReplaceAll(command, "$chaosApi$", params.ChaosApi)
	command = strings.ReplaceAll(command, "$resolver$", params.ResolverPath)
	command = strings.ReplaceAll(command, "$wordlist$", params.WordListPath)
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

func fileExists(filename string) bool {
	// Use os.Stat to get file information
	_, err := os.Stat(filename)

	// Check if the error indicates that the file doesn't exist
	if os.IsNotExist(err) {
		return false
	}

	// Return true if no error occurred (file exists or other error occurred)
	return err == nil
}
