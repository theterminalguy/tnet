package main

import (
	"fmt"
	"os"

	"github.com/10hourlabs/tentn/internal/router"
)

func main() {
	e := router.DefineRoutes()
	httpPort := fmt.Sprintf(":%v", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(httpPort))
}
