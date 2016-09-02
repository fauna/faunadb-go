package faunadb

import (
	"fmt"
	"net/http"
)

var errorsField = ObjKey("errors")

type FaunaError interface {
	Status() int
	Errors() []QueryError
}

type QueryError struct {
	Position    []string            `fauna:"position"`
	Code        string              `fauna:"code"`
	Description string              `fauna:"description"`
	Failures    []ValidationFailure `fauna:"failures"`
}

type ValidationFailure struct {
	Field       []string `fauna:"field"`
	Code        string   `fauna:"code"`
	Description string   `fauna:"description"`
}

type responseError struct {
	status int
	errors []QueryError
}

func (err responseError) Error() string        { return "" } //FIXME better message
func (err responseError) Status() int          { return err.status }
func (err responseError) Errors() []QueryError { return err.errors }

type BadRequest struct{ responseError }
type Unauthorized struct{ responseError }
type NotFound struct{ responseError }
type InternalError struct{ responseError }
type Unavailable struct{ responseError }
type UnknownError struct{ responseError }

func checkForResponseErrors(response *http.Response) error {
	if response.StatusCode < 300 {
		return nil
	}

	errors, ok := parseResponseErrors(response)
	if !ok {
		return fmt.Errorf("Fixme") //FIXME
	}

	err := responseError{response.StatusCode, errors}

	switch response.StatusCode {
	case 400:
		return BadRequest{err}
	case 401:
		return Unauthorized{err}
	case 404:
		return NotFound{err}
	case 500:
		return InternalError{err}
	case 503:
		return Unavailable{err}
	default:
		return UnknownError{err}
	}
}

func parseResponseErrors(response *http.Response) (errors []QueryError, ok bool) {
	if response.Body != nil {
		if res, err := parseJSON(response.Body); err == nil {
			if err := res.At(errorsField).Get(&errors); err == nil {
				ok = true
			}
		}
	}

	return
}
