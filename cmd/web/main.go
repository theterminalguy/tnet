package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/10hourlabs/tentn/internal/router"
	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func inDevMode() bool {
	return strings.HasPrefix(os.Getenv("ENV"), "dev")
}

func main() {
	e := router.DefineRoutes()
	e.Static("/public/views/css", "public/views/css")
	e.Static("/public/views/img", "public/views/img")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/index.html")),
	}
	e.Renderer = renderer

	httpPort := fmt.Sprintf(":%v", os.Getenv("PORT"))
	if inDevMode() && os.Getenv("SSL") == "" {
		fmt.Println("Running in development mode with SSL")
		e.Logger.Fatal(e.StartTLS(httpPort, "cert.pem", "cert-key.pem"))
	} else {
		fmt.Println("Running in development mode without SSL")
		// Cloud Run automatically enforces TLS
		e.Logger.Fatal(e.Start(httpPort))
	}
}
