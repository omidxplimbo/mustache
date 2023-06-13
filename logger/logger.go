package logger

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Info(m string) {
	color.Green("[Info] %s \n", m)
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
