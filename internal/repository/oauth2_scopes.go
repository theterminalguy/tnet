package repository

import "strings"

// Scopes

// recruiter.talents.search // GET
// recruiter.talents.reads // GET
// recruiter.talents.creates // POST
// recruiter.talents.updates // PUT
// recruiter.talents.deletes // DELETE

var OauthScopes = map[string]string{
	"recruiter.talents.search": "Search and filter for talents",
	"recruiter.talents.read":   "Find a talent by ID",
}

// /v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060 == "recruiter.talents.readByID"

// implement a function that checks if a path matches a scope

func PathToScope(path string) string {
	// this algorithm is not perfect, but it's good enough for now
	paths := strings.Split(path, "/")
	if len(paths) < 2 {
		return ""
	}
	return strings.Join(paths[1:], ".")
}

func PathRequiresScope(path string, scope string) bool {
	return strings.HasPrefix(path, scope)
}
