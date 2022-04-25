package util

import "github.com/10hourlabs/tenlog"

// LogAndReturnError logs the error and returns it
// TODO: https://go.dev/blog/go1.13-errors
// Update tenlog to do this
// should be able to post errors to slack
func LogAndReturnErr(msg string, err error) error {
	tenlog.Error(msg, err)
	return err
}
