package util

import "github.com/10hourlabs/tenlog"

// LogAndReturnError logs the error and returns it
// TODO: https://go.dev/blog/go1.13-errors
func LogAndReturnErr(err error, severity tenlog.LogLevel, msg string) error {
	tenlog.Error(msg, err)
	return err
}
