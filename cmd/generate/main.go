// +build ignore

package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/10hourlabs/tentn/util"
)

// template files
const ()

type Templater interface {
	OutDir() string
	Exists() bool
	Name() string
	Generate() error
}

type Template struct {
	FileName string
	Entity   string
	Kind     string
}

func (t *Template) Exists() bool {
	_, err := os.Stat(t.OutDir())
	return err == nil
}

func (t *Template) Name() string {
	return t.FileName
}

func (t *Template) Generate() error {
	tmpl, err := template.ParseFiles(t.sourceDir())
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Template) OutDir() string {
	switch t.Kind {
	case "service":
		return fmt.Sprintf("internal/service/%v_service.go", t.FileName)
	case "repository":
		return fmt.Sprintf("internal/repository/%v_repository.go", t.FileName)
	default:
		panic("unknown template kind")
	}
}

func (t *Template) sourceDir() string {
	switch t.Kind {
	case "service":
		return "cmd/generate/templates/service.tmpl"
	case "repository":
		return "cmd/generate/templates/repository.tmpl"
	default:
		panic("unknown template kind")
	}
}

func GenerateTemplates(entityName, fileName string) {
	var templates []Templater
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "service",
	})
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "repository",
	})
	for _, t := range templates {
		if t.Exists() {
			fmt.Println("[skipping] " + t.OutDir())
			continue
		}
		t.Generate()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a resource name. For example: job_application")
		os.Exit(2)
	}
	fileName := os.Args[1]
	resourceName := util.TitlelizeUnderscore(fileName)
	GenerateTemplates(resourceName, fileName)
}
