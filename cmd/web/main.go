package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/theterminalguy/tentn/internal/router"
	"github.com/theterminalguy/tentn/util"
	"github.com/theterminalguy/tentn/util/osutil"
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

type HelloWorldHandler struct {
	Subscription string
}

func main() {
	e := router.DefineRoutes()
	e.Static("/public/views/css", "public/views/css")
	e.Static("/public/views/img", "public/views/img")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer
	e.HTTPErrorHandler = util.CustomHTTPErrorHandler

	// bus := whisper.NewEventBus(context.Background(), os.Getenv("PUBSUB_EVENT_CONNECTION_NAME"))
	// bus.RegisterEvents(&event.HelloWorldEvent{})
	// go func() {
	// 	if err := whisper.Listen(bus, whisper.NewGooglePubSub()); err != nil {
	// 		log.Fatalf("failed to subscribe: %v\n", err)
	// 	}
	// }()

	httpPort := fmt.Sprintf(":%v", os.Getenv("PORT"))
	if osutil.InDevMode() && os.Getenv("SSL") == "" {
		fmt.Println("Running in development mode with SSL")
		e.Logger.Fatal(e.StartTLS(httpPort, "cert.pem", "cert-key.pem"))
	} else {
		fmt.Println("Running in development mode without SSL")
		// Cloud Run automatically enforces TLS
		e.Logger.Fatal(e.Start(httpPort))
	}
}
