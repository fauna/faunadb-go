package faunadb

import (
	"fmt"
	"net/http"
	"strings"
)

var errorsField = ObjKey("errors")

type FaunaError interface {
	error
	Status() int
	Errors() []QueryError
}

type BadRequest struct{ FaunaError }
type Unauthorized struct{ FaunaError }
type NotFound struct{ FaunaError }
type InternalError struct{ FaunaError }
type Unavailable struct{ FaunaError }
type UnknownError struct{ FaunaError }

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

type errorResponse struct {
	parseable bool
	status    int
	errors    []QueryError
}

func (err errorResponse) Status() int          { return err.status }
func (err errorResponse) Errors() []QueryError { return err.errors }

func (err errorResponse) Error() string {
	return fmt.Sprintf("Response error %d. %s", err.status, err.queryErrors())
}

func (err *errorResponse) queryErrors() string {
	if !err.parseable {
		return "Unparseable server response."
	}

	var errors []string

	for _, queryError := range err.errors {
		errors = append(errors,
			fmt.Sprintf("[%s](%s): %s", strings.Join(queryError.Position, "/"), queryError.Code, queryError.Description))
	}

	return fmt.Sprintf("Errors: %s", strings.Join(errors, ", "))
}

func checkForResponseErrors(response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	err := parseErrorResponse(response)

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

func parseErrorResponse(response *http.Response) FaunaError {
	var errors []QueryError

	if response.Body != nil {
		if value, err := parseJSON(response.Body); err == nil {
			if err := value.At(errorsField).Get(&errors); err == nil {
				return errorResponse{true, response.StatusCode, errors}
			}
		}
	}

	return errorResponse{false, response.StatusCode, errors}
}
