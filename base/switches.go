package base

import (
	"flag"
	"fmt"
	projectModule "github.com/omidxplimbo/mustache/modules/project"
	"os"
)

func HandleSwitch() {

	// Define the flow flag with the options
	help := flag.Bool("h", false, "Display Help")
	flow := flag.String("flow", "", "The Flow To Use (general, init-http, init-dns, init-ip)")
	target := flag.String("target", "", "Target Name")
	pathFile := flag.String("path", "", "Path Of File For Initiate To Database")
	dbName := flag.String("db", "", "Database Name For Target")
	module := flag.String("module", "", "Set module for run")
	projectName := flag.String("project", "", "Create project in database")
	collection := flag.String("collection", "all", "Collection project in database")
	flag.Parse()

	// Show help function
	if *help {
		ShowHelp()
	}

	// Handle Switches
	switch *module {
	case "project":
		projectHandle(*flow, *projectName, *collection, *pathFile)
	case "continues":
		continuesHandle(*flow, *projectName, *dbName, *pathFile, *target)
	case "recon":
		reconHandle(*flow, *projectName, *dbName, *pathFile, *target)
	default:
		fmt.Println("Wrong Selected Module. Please use -h for get help")
	}

}

func reconHandle(flow string, projectName string, dbName string, path string, target string) {
	switch flow {
	default:
		fmt.Println("Recon Flow Very Soon")
	}
}

func continuesHandle(flow string, projectName string, dbName string, path string, target string) {

	switch flow {
	case "general":
		fmt.Println("Using general flow")

	case "init-http":

		if path == "" {
			fmt.Println("Path flag cannot be empty When use Initiate Flags. Please Use -h For Get Help")
			os.Exit(0)
		}
		if dbName == "" {
			fmt.Println("DB flag cannot be empty When use Initiate Flags. Please Use -h For Get Help")
			os.Exit(0)
		}

	case "init-dns":
		fmt.Println("Using init-dns flow")
	case "init-ip":
		fmt.Println("Using init-ip flow")
	case "http-probe":
		fmt.Println("Using http-probe flow")
	default:
		fmt.Println("Flow Flag Cannot Be Empty. Please Use -h For Get Help")
		os.Exit(0)

	}
}

func projectHandle(flow string, projectName string, collection string, path string) {

	project := projectModule.Project{}
	switch flow {
	case "init":
		project.Init(projectName)
	case "backup":
		project.Backup(projectName, collection, path)
	case "report":
		project.Report(projectName)
	case "remove":
		project.Remove(projectName)
	default:
		fmt.Println("Flow Flag Cannot Be Empty. Please Use -h For Get Help")
		os.Exit(0)
	}

}
