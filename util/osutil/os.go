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
	"JWT_SIGNED_SECRET",
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

	"PAYMENT_DRIVER",
	"STRIPE_API_KEY",
	"STRIPE_PRODUCT_KEY",
	"STRIPE_ENDPOINT_SECRET",
}

var requiredRuntimeEnvForProd = []string{
	"CLOUDSQL_PG_USER",
	"CLOUDSQL_PG_PASSWORD",
	"CLOUDSQL_PG_DBNAME",
	"CLOUDSQL_PG_SOCKET_DIR",
	"CLOUDSQL_PG_INSTANCE",

	"FILESYSTEM_DRIVER",
	"GOOGLE_APPLICATION_CREDENTIALS",
}

func CheckEnv() {
	if os.Getenv("DEBUG") == "true" {
		return
	}
	// First check if we're running in dev or test mode
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "test" {
		if os.Getenv("PLATFORM") != "docker" {
			// check if .env file exists, this means you are running locally
			if _, err := os.Stat(".env"); os.IsNotExist(err) {
				panic("Missing .env file")
			}
			// set env vars from .env file
			SetEnvFromFile(".env")
		}
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
	} else if os.Getenv("ENV") == "production" {
		for _, key := range requiredRuntimeEnvForProd {
			if os.Getenv(key) == "" {
				errs = append(errs, "Missing required env var for production mode: "+key)
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
		// Do not set ENV var
		if parts[0] == "ENV" {
			continue
		}
		os.Setenv(parts[0], parts[1])
	}
}

func InDevMode() bool {
	return strings.HasPrefix(os.Getenv("ENV"), "dev")
}
