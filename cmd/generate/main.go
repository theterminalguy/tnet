// +build ignore

package main

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/10hourlabs/tentn/util"
)

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
	if t.Exists() {
		fmt.Printf("[skipping] %v\n", t.OutDir())
		return nil
	}
	tmpl, err := template.ParseFiles(t.sourceDir())
	if err != nil {
		return err
	}
	f, err := os.Create(t.OutDir())
	if err != nil {
		fmt.Printf("Failed to generate template: %v\n", err)
	}
	err = tmpl.Execute(f, t)
	if err != nil {
		return err
	}
	fmt.Printf("[created] %v\n", t.OutDir())
	return nil
}

func (t *Template) OutDir() string {
	switch t.Kind {
	case "service":
		return fmt.Sprintf("internal/service/%v_service.go", t.FileName)
	case "repository":
		return fmt.Sprintf("internal/repository/%v_repository.go", t.FileName)
	case "handler":
		return fmt.Sprintf("internal/handler/%v_handler.go", t.FileName)
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
	case "handler":
		return "cmd/generate/templates/handler.tmpl"
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
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "handler",
	})
	for _, t := range templates {
		t.Generate()
	}
}

func InitEntSchema(resourceName string) {
	fmt.Printf("Generating ent schema for %v...\n", resourceName)
	cmd := exec.Command("go", "run", "entgo.io/ent/cmd/ent", "init", resourceName)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("An error occured initializing ent schema! %v\n", err)
	}
	fmt.Printf("[created] %v\n", "ent/schema/person.go")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a resource name. For example: job_application")
		os.Exit(2)
	}
	fileName := os.Args[1]
	resourceName := util.TitlelizeUnderscore(fileName)
	InitEntSchema(resourceName)
	GenerateTemplates(resourceName, fileName)
}
