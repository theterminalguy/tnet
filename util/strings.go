package util

import (
	"math/rand"
	"strings"
)

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
	return strings.Title(s)
}

// RemoveUnderscore
// Takes a string in the form `hello_world`
// Then converts it to `helloworld`
func RemoveUnderscore(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) > 1 {
		var s []string
		s = append(s, parts...)
		return strings.Join(s, "")
	}
	return s
}

func RandomStringElement(s []string) string {
	randPos := rand.Intn(len(s))
	return s[randPos]
}
