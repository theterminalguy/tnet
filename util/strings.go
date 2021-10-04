package util

import "strings"

// TitlelizeUnderscore
// Takes a string in the form `hello_world`
// Then converts it to `HelloWorld` with each first letter capitalized
func TitlelizeUnderscore(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) > 1 {
		var s []string
		for _, part := range parts {
			s = append(s, strings.Title(part))
		}
		return strings.Join(s, "")
	}
	return s
}
