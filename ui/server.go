package ui

import (
	"fmt"
	"github.com/omidxplimbo/mustache/config"
	"github.com/omidxplimbo/mustache/logger"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func Server(portServer int) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("Run url: %s", r.URL.String()))
		http.ServeFile(w, r, "ui/index.html")
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			command := r.FormValue("command")
			logger.Info(fmt.Sprintf("Run url: %s with -- mustache %s --  command at the %s", r.URL.String(), command, time.Now().Format(config.ProjectConfig().TimeShow)))

			// Split the command string into separate arguments
			args := strings.Fields(command)

			// Run the appropriate function based on the inputs
			cmd := exec.Command("go", append([]string{"run", "main.go"}, args...)...)
			output, err := cmd.Output()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(output)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", portServer), nil)

}
