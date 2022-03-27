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
<h1>Welcome to the future of Recruiting</h1> 
<i>Talent Network API version 0.0.1</i>

<h2>Features</h2>
<ul>
<li>AI Powered Job Recruiting and Matching Engine</li>
<li>Automate your JOB posting and searching process</li>
<li>Use Natural Language to interface with our AI powered job search engine</li>
</ul>


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
