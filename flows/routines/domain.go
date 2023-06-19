package routines

import (
	"bufio"
	"fmt"
	"github.com/omidxplimbo/mustache/base/parser"
	"github.com/omidxplimbo/mustache/base/runner/flows"
	"github.com/omidxplimbo/mustache/db"
	"github.com/omidxplimbo/mustache/logger"
	"os"
	"time"
)

func Domain(args ...interface{}) {

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
	data := parser.Parser{}.YamlParse("ip", true)

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

	setIpSubdomain(projectName)

}

func setDomain(projectName string) {

	var domains []db.Domain

	filePath := fmt.Sprintf("%s-root-domain.txt", projectName)
	var domain []string

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
		domain = append(domain, line)
	}

	if len(domain) > 0 {
		for _, line := range domain {
			if line == "" {
				continue
			}
			domainModel := db.Domain{
				Domain:      line,
				CreatedDate: time.Now(),
				UpdatedDate: time.Now(),
			}
			domains = append(domains, domainModel)
		}
	}

	if len(domains) > 0 {
		domainDb := db.Domain{}
		domainDb.InsertDomain(domains, projectName)
	}
}
