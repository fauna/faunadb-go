package faunadb

import "net/http"

type FaunaError struct{}

func (err FaunaError) Error() string { return "" }

type BadRequest struct{ FaunaError }
type Unauthorized struct{ FaunaError }
type NotFound struct{ FaunaError }
type InternalError struct{ FaunaError }
type Unavailable struct{ FaunaError }
type UnknownError struct{ FaunaError }

func checkForResponseErrors(response *http.Response) (err error) {
	if response.StatusCode < 300 {
		return
	}

	switch response.StatusCode {
	case 400:
		err = BadRequest{}
	case 401:
		err = Unauthorized{}
	case 404:
		err = NotFound{}
	case 500:
		err = InternalError{}
	case 503:
		err = Unavailable{}
	default:
		err = UnknownError{}
	}

	return
}
