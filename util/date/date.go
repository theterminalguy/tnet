package date

import (
	"time"
)

// JSStringToRFC3339 converts a datatime string in the form `2006-01-02T15:04:05Z07:00` to RFC3339
func JSStringToRFC3339(dt string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ToRFC3339(dt time.Time) string {
	return dt.Format(time.RFC3339)
}
