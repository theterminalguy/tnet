// +build ignore

package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/10hourlabs/tentn/util"
)

// template files
const (
	serviceTemplate     string = "cmd/generate/templates/service.tmpl"
	respositoryTemplate string = "cmd/generate/templates/repository.tmpl"
)

type Resource struct {
	Entity   string
	FileName string
}

func (r *Resource) ServiceFileExists() bool {
	_, err := os.Stat(r.ServiceFilePath())
	return err == nil
}

func (r *Resource) ServiceFilePath() string {
	return fmt.Sprintf("internal/service/%v_service.go", r.FileName)
}

func (r *Resource) GenerateTemplate() {
	tmpl, err := template.ParseFiles(serviceTemplate, respositoryTemplate)
	if err != nil {
		panic(err)
	}
	for _, t := range tmpl.Templates() {
		switch t.Name() {
		case "service.tmpl":
			if r.ServiceFileExists() {
				fmt.Println("[skipping] " + r.ServiceFilePath())
				continue
			}
			err = t.Execute(os.Stdout, r)
			if err != nil {
				panic(err)
			}
		case "repository.tmpl":
			if r.ServiceFileExists() {
				fmt.Println("[skipping] " + r.ServiceFilePath())
				continue
			}
			err = t.Execute(os.Stdout, r)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a resource name. For example: job_application")
		os.Exit(2)
	}
	fileName := os.Args[1]
	resourceName := util.TitlelizeUnderscore(fileName)
	resource := Resource{Entity: resourceName, FileName: fileName}
	if resource.ServiceFileExists() {
		fmt.Println("[skipping] " + resource.ServiceFilePath())
	} else {
		resource.GenerateTemplate()
	}
}
