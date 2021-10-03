package main

import (
	"github.com/10hourlabs/tentn/internal/router"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	e := router.DefineRoutes()

	e.Logger.Fatal(e.Start(":1323"))
}
