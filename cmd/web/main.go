package main

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/router"
	"github.com/10hourlabs/tentn/util/osutil"
)

func main() {
	e := router.DefineRoutes()
	httpPort := fmt.Sprintf(":%v", osutil.Getenv("PORT", "8080"))
	e.Logger.Fatal(e.Start(httpPort))
}
