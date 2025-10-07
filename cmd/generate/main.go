// +build ignore

package main

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/theterminalguy/tentn/util"
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
	PackageName string
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
	defer f.Close()
	if err != nil {
		fmt.Printf("Unable to create file: %v\n", err)
	}
	err = tmpl.Execute(f, t)
	if err != nil {
		fmt.Printf("Failed to generate template: %v\n", err)
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

func GenerateTemplates(entityName, fileName, packageName string) {
	var templates []Templater
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "service",
		PackageName: packageName,
	})
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "repository",
		PackageName: packageName,
	})
	templates = append(templates, &Template{
		Entity:   entityName,
		FileName: fileName,
		Kind:     "handler",
		PackageName: packageName,
	})
	for _, t := range templates {
		t.Generate()
	}
}

func InitEntSchema(resourceName, fileName string) error {
	outDir := fmt.Sprintf("ent/schema/%v.go", fileName)
	if _, err := os.Stat(outDir); err == nil {
		fmt.Printf("[skipping] %v\n", outDir)
		return nil
	}
	cmd := exec.Command("go", "run", "entgo.io/ent/cmd/ent", "init", resourceName)
	cmd.Stderr = os.Stderr
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("An error occured initializing ent schema! %v\n", err)
	}
	fmt.Printf("[created] %v\n", outDir)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a resource name. For example: job_application")
		os.Exit(64)
	}
	fileName := os.Args[1]
	resourceName := util.TitlelizeUnderscore(fileName)
	packageName := util.RemoveUnderscore(fileName)
	fmt.Printf("Generating scaffold for %v...\n", resourceName)
	if err := InitEntSchema(resourceName, fileName); err != nil {
		fmt.Println(err)
		os.Exit(64)
	}
	GenerateTemplates(resourceName, fileName, packageName)
}
