package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	html := `<!DOCTYPE html>
<html>
<head>
<title>Tentn</title>
</head>
<body>
<h1>Tentn</h1> 
<p>Talent Network API version 0.0.1</p>

<p>
	<a href="/recruiter/auth">Add to Slack </a>
</p>

<p>
	<a href="/talent/auth">Talent Login</a>
</p>
</body>
</html>`
	// render html
	return c.HTMLBlob(http.StatusOK, []byte(html))
}
