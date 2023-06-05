package watch

import "github.com/omidxplimbo/mustache/db"

func GetAllSubdomain(projectName string) {

	subdomain := db.Subdomain{}
	subdomain.GetAllSubdomain(projectName)
}
