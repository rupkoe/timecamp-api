// Package tc-api contains api types and helper functions on these types (like filtering).
package api

import (
	"io/ioutil"
	"net/http"
)

const DateFormat = "2006-01-02"

type Connection struct {
	ApiUrl string
	Token  string
}

func httpGet(url string) ([]byte, error) {
	var data []byte

	response, err := http.Get(url)
	if err != nil {
		return data, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}
