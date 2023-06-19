package watch

import "github.com/omidxplimbo/mustache/db"

func GetAllSubdomain(projectName string, justCount *bool) {

	subdomain := db.Subdomain{}
	subdomain.GetAllSubdomain(projectName, justCount)
}

func LatestSubdomain(projectName string, count int) {

	subdomain := db.Subdomain{}
	subdomain.LatestSubdomain(projectName, count)
}

func LivesSubdomain(projectName string, justCount *bool) {

	subdomain := db.Subdomain{}
	subdomain.LivesSubdomain(projectName, justCount)
}

func LatestLivesSubdomain(projectName string, count int) {

	subdomain := db.Subdomain{}
	subdomain.LatestLivesSubdomain(projectName, count)
}

func GetSub(projectName string, target string) {

	subdomain := db.Subdomain{}
	subdomain.GetSub(projectName, target)
}

func ResolvedSubdomain(projectName string, justCount *bool) {

	subdomain := db.Subdomain{}
	subdomain.ResolvedSubdomain(projectName, justCount)
}

func Domain(projectName string, justCount *bool) {

	domain := db.Domain{}
	domain.GetAllDomains(projectName, justCount)
}

func GetDomain(projectName string, target string) {

	domain := db.Domain{}
	domain.GetDomain(projectName, target)
}
