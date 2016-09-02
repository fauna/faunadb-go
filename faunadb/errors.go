package faunadb

import "net/http"

type BadRequest struct{}

func (err BadRequest) Error() string {
	return "Bad request"
}

func checkForResponseErrors(response *http.Response) (err error) {
	if response.StatusCode < 300 {
		return
	}

	err = BadRequest{}

	return
}
