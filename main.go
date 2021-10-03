package main

import (
	"github.com/10hourlabs/tentn/internal/router"
)

func main() {
	e := router.DefineRoutes()

	e.Logger.Fatal(e.Start(":1323"))
}
