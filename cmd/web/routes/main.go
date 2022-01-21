package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/10hourlabs/tentn/internal/router"
)

func main() {
	e := router.DefineRoutes()
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("routes.json", data, 0644)
}
