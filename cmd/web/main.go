package main

import (
	"fmt"
	"io"
	"os"
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

func main() {
	e := router.DefineRoutes()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer

	httpPort := fmt.Sprintf(":%v", os.Getenv("PORT"))
	if os.Getenv("ENV") == "dev" {
		e.Logger.Fatal(e.StartTLS(httpPort, "cert.pem", "cert-key.pem"))
	} else {
		// Cloud Run automatically enforces TLS
		e.Logger.Fatal(e.Start(httpPort))
	}
}
