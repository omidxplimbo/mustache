package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/omidxplimbo/mustache/config"
	"os"
)

func Info(m string) {
	color.Green("[Info] %s \n", m)
}

func Process(m string) {
	color.HiGreen("[+] %s \n", m)
}

func Warning(m string) {
	color.Yellow("[War] %s \n", m)
	fmt.Println()
	os.Exit(0)
}

func Fetal(m string) {
	color.Red("[Fatal] %s \n", m)
	fmt.Println()
	os.Exit(0)
}

func WarningC(m string) {
	color.HiYellow("[Warning] %s \n", m)
	fmt.Println()
}

func AddLog(m string) {
	logsPath := config.ProjectConfig().Addresses.Logs

	// Check if the file exists
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(logsPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
	}

	// Open the file in append mode
	file, err := os.OpenFile(logsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the message to the file
	if _, err := file.WriteString(m + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
