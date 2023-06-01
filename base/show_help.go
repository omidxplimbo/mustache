package base

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func ShowHelp() {

	color.Green("Command Line Usage: ")
	usage := `	mustache -module continues -flow [flowName] -target [targetName] -path [pathFile] -db [databaseName]
	mustache -module recon -flow [flowName] -target [targetName] -path [pathFile] -db [databaseName]
	mustache -module project -flow init -project [projectName] 
	mustache -module project -flow backup -project [projectName]
	mustache -module project -flow remove -project [projectName]
	mustache -module project -flow backup -project [projectName] -collection [collectionName]
	`

	color.Yellow(usage)
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
	}

	watch := [][]string{
		{"subdomain", "Get all subdomain of project", "Watch"},
		{"domain", "Get all domain of project", "Watch"},
		{"tech", "Get all technologies of project", "Watch"},
		{"cidr", "Get all cidr of project", "Watch"},
		{"asn", "Get all asn of project", "Watch"},
		{"urls", "Get all urls of project", "Watch"},
		{"lives", "Get all lives subdomain of project", "Watch"},
		{"latest-sub", "Get latest added subdomain of project (default 10 record)", "Watch"},
		{"latest-lives", "Get latest added live subdomain of project (default 10 record)", "Watch"},
		{"latest-urls", "Get latest added urls of project (default 10 record)", "Watch"},
		{"latest-domain", "Get latest added domain of project (default 10 record)", "Watch"},
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
