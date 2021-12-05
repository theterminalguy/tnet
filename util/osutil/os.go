package osutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var requiredRuntimeEnv = []string{
	"ENV",
	"PORT",
	"APP_HOST",
}

var requiredRuntimeEnvForDev = []string{
	"POSTGRES_HOST",
	"POSTGRES_PORT",
	"POSTGRES_USER",
	"POSTGRES_PASSWORD",
	"POSTGRES_DB",
	"POSTGRES_SSL_MODE",

	"GOOGLE_OAUTH_CLIENT_ID",
	"GOOGLE_OAUTH_CLIENT_SECRET",

	"JWT_SIGNED_SECRET",
}

var requiredRuntimeEnvForProd = []string{
	"CLOUDSQL_PG_USER",
	"CLOUDSQL_PG_PASSWORD",
	"CLOUDSQL_PG_SOCKET",
	"CLOUDSQL_PG_INSTANCE",
	"CLOUDSQL_PG_DBNAME",
}

func CheckEnv() {
	// First check if we're running in dev or test mode
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "test" {
		// check if .env file exists
		if _, err := os.Stat(".env"); os.IsNotExist(err) {
			panic("Missing .env file")
		}
		// set env vars from .env file
		SetEnvFromFile(".env")
	}

	var errs []string

	for _, key := range requiredRuntimeEnv {
		if os.Getenv(key) == "" {
			errs = append(errs, "Missing required env var: "+key)
		}
	}
	if os.Getenv("ENV") == "dev" {
		for _, key := range requiredRuntimeEnvForDev {
			if os.Getenv(key) == "" {
				errs = append(errs, "Missing required env var for dev mode: "+key)
			}
		}
	} else if os.Getenv("ENV") == "prod" {
		for _, key := range requiredRuntimeEnvForProd {
			if os.Getenv(key) == "" {
				errs = append(errs, "Missing required env var for prod mode: "+key)
			}
		}
	}

	if len(errs) > 0 {
		panic(strings.Join(errs, "\n"))
	}
}

// SetEnvFromFile sets env vars from a file
// Only use this for dev mode
func SetEnvFromFile(file string) {
	fmt.Println("Setting env vars from file:", file)
	fileContent, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fileContent.Close()

	scanner := bufio.NewScanner(fileContent)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			panic("invalid line in .env file: " + line)
		}
		os.Setenv(parts[0], parts[1])
	}
}
