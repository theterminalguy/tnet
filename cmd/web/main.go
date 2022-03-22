package main

import (
	"fmt"
	"os"

	"github.com/10hourlabs/tentn/internal/router"
)

func main() {
	e := router.DefineRoutes()
	httpPort := fmt.Sprintf(":%v", os.Getenv("PORT"))
	if os.Getenv("ENV") == "dev" {
		e.Logger.Fatal(e.StartTLS(httpPort, "cert.pem", "cert-key.pem"))
	} else {
		// Cloud Run automatically enforces TLS
		e.Logger.Fatal(e.Start(httpPort))
	}
}
