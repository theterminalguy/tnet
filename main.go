package main

import (
	"net/http"

	"github.com/10hourlabs/jobsapi/internal/job"
	"github.com/labstack/echo"
)

func main() {
	j := &job.Job{}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, j.Hello())
	})
	e.Logger.Fatal(e.Start(":1323"))
}
