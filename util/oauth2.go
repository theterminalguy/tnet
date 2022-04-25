package util

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

func PathToOauth2Scope(reqPath, verb string) string {
	if _, ok := verbActionSingular[verb]; !ok {
		return ""
	}
	if _, ok := verbActionPlural[verb]; !ok {
		return ""
	}
	paths := strings.Split(reqPath, "v1/")
	if len(paths) < 2 {
		return ""
	}
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
