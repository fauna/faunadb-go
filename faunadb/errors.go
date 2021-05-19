package faunadb

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var errorsField = ObjKey("errors")

// A FaunaError wraps HTTP errors when sending queries to a FaunaDB cluster.
type FaunaError interface {
	error
	Status() int          // HTTP status code
	Errors() []QueryError // Errors returned by the server
}

// A BadRequest wraps an HTTP 400 error response.
type BadRequest struct{ FaunaError }

// A Unauthorized wraps an HTTP 401 error response.
type Unauthorized struct{ FaunaError }

// A PermissionDenied wraps an HTTP 403 error response.
type PermissionDenied struct{ FaunaError }

// A NotFound wraps an HTTP 404 error response.
type NotFound struct{ FaunaError }

// A InternalError wraps an HTTP 500 error response.
type InternalError struct{ FaunaError }

// A Unavailable wraps an HTTP 503 error response.
type Unavailable struct{ FaunaError }

// A UnknownError wraps any unknown http error response.
type UnknownError struct{ FaunaError }

// QueryError describes query errors returned by the server.
type QueryError struct {
	Position    []string            `fauna:"position"`
	Code        string              `fauna:"code"`
	Description string              `fauna:"description"`
	Cause       []ValidationFailure `fauna:"cause"`
}

// ValidationFailure describes validation errors on a submitted query.
type ValidationFailure struct {
	Position    []string `fauna:"position"`
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

	errs := make([]string, len(err.errors))

	for i, queryError := range err.errors {

		errs[i] =
			fmt.Sprintf("[%s](%s): %s, details: %s", strings.Join(queryError.Position, "/"), queryError.Code, queryError.Description, queryError.Cause)
	}

	return fmt.Sprintf("Errors: %s", strings.Join(errs, ", "))
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
	case 403:
		return PermissionDenied{err}
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

func errorFromStreamError(obj ObjectV) (err error) {
	var sb strings.Builder
	sb.WriteString("stream_error:")
	keys := make([]string, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, ok := obj[k]; ok {
			sb.WriteString(" ")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(fmt.Sprintf("'%s'", obj[k]))
		}

	}
	err = errors.New(sb.String())
	return
}
