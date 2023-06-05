package watch

import "github.com/omidxplimbo/mustache/db"

func GetAllSubdomain(projectName string) {

	subdomain := db.Subdomain{}
	subdomain.GetAllSubdomain(projectName)
}

func LatestSubdomain(projectName string, count int) {

	subdomain := db.Subdomain{}
	subdomain.LatestSubdomain(projectName, count)
}

func LivesSubdomain(projectName string) {

	subdomain := db.Subdomain{}
	subdomain.LivesSubdomain(projectName)
}

func LatestLivesSubdomain(projectName string, count int) {

	subdomain := db.Subdomain{}
	subdomain.LatestLivesSubdomain(projectName, count)
}

func GetSub(projectName string, target string) {

	subdomain := db.Subdomain{}
	subdomain.GetSub(projectName, target)
}
