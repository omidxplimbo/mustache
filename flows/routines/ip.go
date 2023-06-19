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

type ResolvedIp struct {
	Subdomain string
	Ip        string
}

func Ip(args ...interface{}) {

	//Set arguments
	projectName := args[0].(string)
	target := args[1].(string)

	// parse yaml files
	data := parser.Parser{}.YamlParse("ip", true)

	// no switch
	var flagsMap map[string]interface{}

	// run flow
	flows.RunnerFlow{}.ExecuteFlows(data, projectName, target, flagsMap)

	setIpSubdomain(projectName)

}

func setIpSubdomain(projectName string) {

	var rs []db.Subdomain
	var resolves []ResolvedIp

	filePath := fmt.Sprintf("%s-sub-ip.txt", projectName)

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

			resolve := ResolvedIp{
				Subdomain: subdomain,
				Ip:        ip,
			}

			resolves = append(resolves, resolve)
		} else {
			subdomain := parts[0]
			resolve := ResolvedIp{
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
