package main

import (
	"fmt"
	"html/template"
	"os"
)

const (
	serviceTemplateDir string = "cmd/generate/templates/service.tmpl"
)

type Resource struct {
	Entity string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an entity name. For example: Job")
	}
	resource := Resource{Entity: os.Args[1]}
	tmpl, err := template.ParseFiles(serviceTemplateDir)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, resource)
	if err != nil {
		panic(err)
	}
}
