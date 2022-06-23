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

func ExtractFirstNumbers(s string) string {
	var nums []string
	for _, r := range s {
		if r >= '0' && r <= '9' {
			nums = append(nums, string(r))
			continue
		}
		break
	}
	return strings.Join(nums, "")
}

func StringParamsToMap(params string) map[string]string {
	result := make(map[string]string)
	for _, param := range strings.Split(params, ",") {
		kv := strings.Split(param, "=")
		result[kv[0]] = kv[1]
	}
	return result
}

func MapToStringParams(m map[string]interface{}) string {
	var params []string
	for k, v := range m {
		params = append(params, k+"="+v.(string))
	}
	return strings.Join(params, ",")
}

func Titlelize(s string) string {
	return strings.Title(strings.ToLower(s))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
