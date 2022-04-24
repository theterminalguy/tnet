package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var verbActionSingular = map[string]string{
	"GET":    "read",
	"POST":   "create",
	"PUT":    "update",
	"DELETE": "delete",
}

var verbActionPlural = map[string]string{
	"GET":    "reads",
	"POST":   "creates",
	"PUT":    "updates",
	"DELETE": "deletes",
}

func PathToScope(reqPath, verb string) string {
	if _, ok := verbActionSingular[verb]; !ok {
		return ""
	}
	if _, ok := verbActionPlural[verb]; !ok {
		return ""
	}
	paths := strings.Split(reqPath, "v1/")
	paths = strings.Split(paths[1], "/")
	if len(paths) < 3 {
		// No ID in the path so it must be a plural action
		// e.g. /v1/recruiter/talents
		// e.g. /v1/recruiter/talent-collections
		return fmt.Sprintf("%s.%s", strings.Join(paths, "."), verbActionPlural[verb])
	}
	// get the first three elements of the path
	chunks := []string{}
	for i := 0; i < 3; i++ {
		chunks = append(chunks, paths[i])
	}
	lastIdx := len(chunks) - 1
	id := chunks[lastIdx]
	if _, err := uuid.Parse(id); err == nil {
		chunks[lastIdx] = verbActionSingular[verb]
	}
	return strings.Join(chunks[0:], ".")
}

func main() {
	fmt.Println(PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "GET"))
	fmt.Println(PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "POST"))
	fmt.Println(PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "PUT"))
	fmt.Println(PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "DELETE"))
	fmt.Println(PathToScope("/v1/recruiter/talents", "GET"))
	fmt.Println(PathToScope("/v1/recruiter/talents", "POST"))
	fmt.Println(PathToScope("/v1/recruiter/talents", "PUT"))
	fmt.Println(PathToScope("/v1/recruiter/talents", "DELETE"))
	fmt.Println(PathToScope("/v1/recruiter/talents/search", "GET"))
}
