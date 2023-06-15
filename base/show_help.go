package base

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func ShowHelp() {

	color.Green("Command Line Usage: ")
	usagePublic := `	
	Public flags for all modules:
		-module [string]	"Select module that you want to run"
		-flow [string]	"Select flow that you want to run. Each module has some flow"
		-project [string]	"Project that you want use for modules"
		-target [string]	"The target that we want run processes"
	`
	usageProject := `
	Flags and usage for project module:
		-collection [string] "For backup from special collection in backup flow"
		-path [string]	"Path for save backup file"

		mustache -module project -flow init -project [projectName] 
		mustache -module project -flow backup -project [projectName] -path [pathFile]
		mustache -module project -flow backup -project [projectName] -collection [collectionName] -path [pathFile]
		mustache -module project -flow remove -project [projectName]
		mustache -module project -flow report -project [projectName]
		mustache -module project -flow all
		`
	usageRecon := `
	Flags and usage for Recon module:
		-wc 	"If set use shuffledns for resolve subdomain and if not set use dnsx"
		-pn	"If set we dont use permutation in subdomains"
		-pr	"If set we use permutation just in resolved subdomains"
		-pp	"If set we use permutation for all subdomains. Should use for small target"
		-aa	"If set we use active amass for subdomains. It's take long time'"
		-target string	"The target that we want run processes"

		mustache -module recon -flow [flowName] -project [projectName] -target [targetName] [-wc] [-pn] [-pr] [-pp]`
	usageCr := `
	Flags and usage for Continues recon module:
		mustache -module continues -flow [flowName] -target [targetName] -path [pathFile] -db [databaseName]`
	usageWatch := `
	Flags and usage for Watch module:
		-jc "For get just count of items in report flow"
		-count [int] "For get more count of items if exists in report flow"

		mustache -module watch -flow subdomain -project [projectName] [-jc]
		mustache -module watch -flow domain -project [projectName] [-jc]
		mustache -module watch -flow tech -project [projectName] [-jc]
		mustache -module watch -flow latest-lives -project [projectName] [-jc] -count [count]
`
	color.HiYellow(usagePublic)
	color.HiBlue(usageProject)
	color.HiWhite(usageRecon)
	color.HiGreen(usageCr)
	color.HiRed(usageWatch)
	fmt.Println()

	// Define the table headers and data
	data := [][]string{
		{"init-http", "Initiate subdomains with http service for special target into database for continues recon", "CR"},
		{"init-dns", "Initiate subdomains with dns resolve service for special target into database for continues recon", "CR"},
		{"init-ip", "Initiate resolved IP's for special target into database for continues recon", "CR"},
		{"http-probe", "Run http subdomain discovery with and update database if some change exist", "CR"},
	}

	recon := [][]string{
		{"general", "Run general flow for recon", "Recon"},
		{"extensive", "Run extensive flow for recon", "Recon"},
		{"cidr", "Run extensive flow for recon", "Recon"},
	}

	projects := [][]string{
		{"init", "Create project database for recon", "Project"},
		{"backup", "Download backup data from project", "Project"},
		{"report", "Create report from project", "Project"},
		{"remove", "Remove all data and database for special project", "Project"},
		{"all", "Get all projects", "Project"},
	}

	watch := [][]string{
		{"subdomain", "Get all subdomain of project", "Watch"},
		{"domain", "Get all domain of project", "Watch"},
		{"tech", "Get all technologies of project", "Watch"},
		{"cidr", "Get all cidr of project", "Watch"},
		{"asn", "Get all asn of project", "Watch"},
		{"urls", "Get all urls of project", "Watch"},
		{"lives", "Get all lives subdomain of project", "Watch"},
		{"resolved", "Get all resolved subdomain of project", "Watch"},
		{"latest-sub", "Get latest added subdomain of project (default 10 record)", "Watch"},
		{"latest-lives", "Get latest added live subdomain of project (default 10 record)", "Watch"},
		{"latest-urls", "Get latest added urls of project (default 10 record)", "Watch"},
		{"latest-domain", "Get latest added domain of project (default 10 record)", "Watch"},
		{"get-sub", "Get Subdomain information of project", "Watch"},
	}

	update := [][]string{
		{"cdn", "Update CDN ranges", "Update"},
		{"wordlist", "Update static wordList", "Update"},
		{"resolvers", "Update resolvers", "Update"},
	}

	// Define the table format
	format := "| %-18s | %-98s | %-17s|\n"
	headers := []string{"Flow Name", "Description", "Module Name"}

	// banner maker
	color.Blue("Flows For Continues Recon. We Are Mustache Team")
	bannerMaker(data, format, headers)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	color.Green("Flows For Recon. We Are Mustache Team")
	bannerMaker(recon, format, headers)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	color.Red("Work With Projects. We Are Mustache Team")
	bannerMaker(projects, format, headers)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	color.Magenta("Watch Tower Module. We Are Mustache Team")
	bannerMaker(watch, format, headers)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	color.Red("Update Tools Module. We Are Mustache Team")
	bannerMaker(update, format, headers)

	os.Exit(0)
}
func bannerMaker(data [][]string, format string, headers []string) {
	// Print the table header
	color.Cyan("+----------------+-------------------------------------------------------------------------------------------------------+----------------")
	color.Cyan(fmt.Sprintf(format, headers[0], headers[1], headers[2]))
	color.Cyan("+----------------+-------------------------------------------------------------------------------------------------------+----------------")

	// Print the table data
	for _, row := range data {
		color.White(fmt.Sprintf(format, row[0], row[1], row[2]))
	}

	// Print the table footer
	color.Cyan("+----------------+-------------------------------------------------------------------------------------------------------+---------------")
}
