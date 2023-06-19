package routines

import (
	"bufio"
	"fmt"
	"github.com/omidxplimbo/mustache/base/parser"
	"github.com/omidxplimbo/mustache/base/runner/flows"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
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
	WC  bool
	PR  bool
	PN  bool
	PP  bool
	AA  bool
	WCR bool
}

func SubDomain(args ...interface{}) {

	//Set arguments
	projectName := args[0].(string)
	target := args[1].(string)
	wc := args[2].(bool)
	pr := args[3].(bool)
	pp := args[4].(bool)
	pn := args[5].(bool)
	aa := args[6].(bool)
	wcr := args[7].(bool)

	// parse yaml files
	data := parser.Parser{}.YamlParse("subdomain", true)

	// set flags
	allFlags := checkAnySwitches{
		WC:  wc,
		PR:  pr,
		PN:  pn,
		PP:  pp,
		AA:  aa,
		WCR: wcr,
	}
	flagsMap := make(map[string]interface{})
	flagsMap["WC"] = allFlags.WC
	flagsMap["PR"] = allFlags.PR
	flagsMap["PN"] = allFlags.PN
	flagsMap["PP"] = allFlags.PP
	flagsMap["AA"] = allFlags.AA
	flagsMap["WCR"] = allFlags.WCR

	// run flow
	flows.RunnerFlow{}.ExecuteFlows(data, projectName, target, flagsMap)

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
