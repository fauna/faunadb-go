package faunadb

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkForResponseErrors(response *http.Response) error {
	if response.StatusCode >= 300 {
		str, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Query error %v: %s", response.StatusCode, str)
	}

	return nil
}
