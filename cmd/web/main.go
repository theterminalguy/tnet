package main

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/router"
	"github.com/10hourlabs/tentn/util/osutil"
)

func main() {
	e := router.DefineRoutes()
	port := fmt.Sprintf(":%v", osutil.Getenv("HTTP_PORT", "8080"))
	e.Logger.Fatal(e.Start(port))
}
