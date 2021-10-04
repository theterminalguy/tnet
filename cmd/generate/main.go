// +build ignore

package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/10hourlabs/tentn/util"
)

// template directories
const (
	serviceTemplate string = "cmd/generate/templates/service.tmpl"
)

type Resource struct {
	Entity   string
	FileName string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an entity name. For example: job_application")
		os.Exit(2)
	}
	fileName := os.Args[1]
	resourceName := util.TitlelizeUnderscore(fileName)
	resource := Resource{Entity: resourceName, FileName: fileName}
	tmpl, err := template.ParseFiles(serviceTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, resource)
	if err != nil {
		panic(err)
	}
}
