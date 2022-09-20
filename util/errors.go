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

func LogAndReturnErrs(e []error, level tenlog.LogLevel) error {
	errs := CombineErrors(e)
	switch level {
	case tenlog.CRITICAL:
		tenlog.Critical(errs)
	case tenlog.DEBUG:
		tenlog.Debug(errs)
	case tenlog.INFO:
		tenlog.Info(errs)
	case tenlog.WARN:
		tenlog.Warn(errs)
	case tenlog.NOTICE:
		tenlog.Notice(errs)
	default:
		tenlog.Error(errs)
	}
	return errs
}

func CombineErrors(errs []error) error {
	var errStr string
	for i, err := range errs {
		errStr += err.Error() + fmt.Sprintf("(%d)", i)
	}
	return fmt.Errorf(errStr)
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if er, ok := err.(*echo.HTTPError); ok {
		tenlog.Error(string(er.Error()))
		var pcode string
		switch er.Code {
		case http.StatusNotFound:
			pcode = "404"
		case http.StatusUnauthorized:
			c.String(http.StatusUnauthorized, er.Message.(string))
		}
		c.Render(http.StatusOK, fmt.Sprintf("%s.html", pcode), nil)
	}
}
