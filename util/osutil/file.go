package osutil

import (
	"net/http"
)

func ReadCSVFromURL(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
