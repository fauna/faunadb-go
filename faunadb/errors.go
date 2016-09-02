package faunadb

import "net/http"

type FaunaError struct {
}

func (err FaunaError) Error() string {
	return ""
}

type BadRequest struct {
	FaunaError
}

type Unauthorized struct {
	FaunaError
}

func checkForResponseErrors(response *http.Response) (err error) {
	if response.StatusCode < 300 {
		return
	}

	switch response.StatusCode {
	case 400:
		err = BadRequest{}
	case 401:
		err = Unauthorized{}
	}

	return
}
