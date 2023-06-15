package routines

import (
	"bufio"
	"fmt"
	"github.com/omidxplimbo/mustache/base/parser"
	"github.com/omidxplimbo/mustache/base/runner/flows"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"io"
	"os"
	"strings"
	"time"
)

type ResolvesSubdomain struct {
	Subdomain string
	Ip        string
}

type LiveHosts struct {
	Subdomain string
	CDN       string
	Http      bool
}

type checkAnySwitches struct {
	WC bool
	PR bool
	PN bool
	PP bool
}

func SubDomain(args ...interface{}) {

	//Set arguments
	projectName := args[0].(string)
	target := args[1].(string)
	wc := args[2].(bool)
	pr := args[3].(bool)
	pp := args[4].(bool)
	pn := args[5].(bool)

	// parse yaml files
	data := parser.Parser{}.YamlParse("subdomain", true)

	// set flags
	allFlags := checkAnySwitches{
		WC: wc,
		PR: pr,
		PN: pn,
		PP: pp,
	}
	flagsMap := make(map[string]interface{})
	flagsMap["WC"] = allFlags.WC
	flagsMap["PR"] = allFlags.PR
	flagsMap["PN"] = allFlags.PN
	flagsMap["PP"] = allFlags.PP

	// execute flow if not data exist
	filePassiveStatus := fmt.Sprintf("%s-passive.txt", projectName)
	fileResolveStatus := fmt.Sprintf("%s-final.txt", projectName)
	fileLivesStatus := fmt.Sprintf("%s-lives.txt", projectName)
	fileResolveForLiveStatus := fmt.Sprintf("%s-finalLive.txt", projectName)
	if !fileExists(fileResolveStatus) && !fileExists(filePassiveStatus) && !fileExists(fileLivesStatus) {

		flows.RunnerFlow{}.ExecuteFlows(data, projectName, target, flagsMap)

		if !fileExists(fileResolveForLiveStatus) {
			modifyFile(fileResolveStatus, fileResolveForLiveStatus)
		}
	}

	//set passive and resolves into database
	setPassiveSubdomain(projectName)
	setResolveSubdomain(projectName)
	setLiveHosts(projectName)
}

func setLiveHosts(projectName string) {
	var lives []db.Subdomain
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
			subdomain := strings.ReplaceAll(parts[0], "https://", "")
			subdomain = strings.ReplaceAll(subdomain, "http://", "")
			CDN := strings.Trim(parts[1], "[]")

			liveHost := LiveHosts{
				Subdomain: subdomain,
				Http:      true,
				CDN:       CDN,
			}

			liveHosts = append(liveHosts, liveHost)
		} else {
			subdomain := parts[0]
			subdomain = strings.ReplaceAll(subdomain, "https://", "")
			subdomain = strings.ReplaceAll(subdomain, "http://", "")
			liveHost := LiveHosts{
				Subdomain: subdomain,
				Http:      true,
			}

			liveHosts = append(liveHosts, liveHost)
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

	if len(lives) > 0 {
		liveSubdomainsDb := db.Subdomain{}
		liveSubdomainsDb.InsertSubdomain(lives, projectName)
	}
}

func setPassiveSubdomain(projectName string) {

	var subdomains []db.Subdomain

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

	if len(passive) > 0 {
		for _, line := range passive {
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

	if len(subdomains) > 0 {
		subdomainDb := db.Subdomain{}
		subdomainDb.InsertSubdomain(subdomains, projectName)
	}
}

func setResolveSubdomain(projectName string) {

	var rs []db.Subdomain
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

	if len(rs) > 0 {
		resolveSubdomainDb := db.Subdomain{}
		resolveSubdomainDb.InsertSubdomain(rs, projectName)
	}
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

func modifyFile(fileResolveStatus string, fileResolveForLiveStatus string) {
	sourceFile, err := os.Open(fileResolveStatus)
	if err != nil {
		logger.Warning(err.Error())
	}
	defer sourceFile.Close()

	destFile, err := os.Create(fileResolveForLiveStatus)
	if err != nil {
		logger.Warning(err.Error())
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		logger.Warning(err.Error())
	}
	// Open the text file for reading and writing
	file, err := os.OpenFile(fileResolveForLiveStatus, os.O_RDWR, 0644)
	if err != nil {
		logger.Warning(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// Create a buffer to store the modified lines
	var lines []string

	// Loop over each line and remove [ip] if it exists
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, " [") {
			line = strings.Split(line, " [")[0]
		}
		lines = append(lines, line)
	}
	// Write the modified lines back to the file
	_, err = file.Seek(0, 0)
	if err != nil {
		logger.Warning(err.Error())
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	err = writer.Flush()
	if err != nil {
		logger.Warning(err.Error())
	}
}
