package util

import (
	"fmt"
	"net/http"

	"github.com/10hourlabs/tenlog"
	"github.com/labstack/echo/v4"
)

// LogAndReturnError logs the error and returns it
// TODO: https://go.dev/blog/go1.13-errors
// Update tenlog to do this
// should be able to post errors to slack
func LogAndReturnErr(msg string, err error) error {
	tenlog.Error(msg, err)
	return err
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if er, ok := err.(*echo.HTTPError); ok {
		tenlog.Error(string(er.Error()))
		var pcode string
		switch er.Code {
		case 404:
			pcode = "404"
		}
		c.Render(http.StatusOK, fmt.Sprintf("%s.html", pcode), nil)
	}
}
