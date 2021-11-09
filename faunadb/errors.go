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
	HttpStatusCode() int  // HTTP status code
	Errors() []QueryError // Errors returned by the server
}

// The following errors wrap an HTTP 400, 403 and 404 error responses.
type InvalidArgumentError struct{ FaunaError }
type FunctionCallError struct{ FaunaError }
type PermissionDeniedError struct{ FaunaError }
type InvalidExpressionError struct{ FaunaError }
type InvalidUrlParameterError struct{ FaunaError }
type TransactionAbortedError struct{ FaunaError }
type InvalidWriteTimeError struct{ FaunaError }
type InvalidReferenceError struct{ FaunaError }
type MissingIdentityError struct{ FaunaError }
type InvalidTokenError struct{ FaunaError }
type StackOverflowError struct{ FaunaError }
type AuthenticationFailedError struct{ FaunaError }
type ValueNotFoundError struct{ FaunaError }
type InstanceNotFoundError struct{ FaunaError }
type InstanceAlreadyExistsError struct{ FaunaError }
type ValidationFailedError struct{ FaunaError }
type InstanceNotUniqueError struct{ FaunaError }
type FeatureNotAvailableError struct{ FaunaError }

// A Unauthorized wraps an HTTP 401 error response.
type Unauthorized struct{ FaunaError }

// A TransactionContention wraps an HTTP 409 error response.
type TransactionContention struct{ FaunaError }

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

func (err errorResponse) HttpStatusCode() int  { return err.status }
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
	case 400, 403, 404:
		return queryError(err)
	case 401:
		return Unauthorized{err}
	case 409:
		return TransactionContention{err}
	case 500:
		return InternalError{err}
	case 503:
		return Unavailable{err}
	default:
		return UnknownError{err}
	}
}

func queryError(err FaunaError) error {
	if len(err.Errors()) == 0 {
		return UnknownError{err}
	}

	code := err.Errors()[0].Code
	switch code {
	case "invalid argument":
		return InvalidArgumentError{err}
	case "call error":
		return FunctionCallError{err}
	case "permission denied":
		return PermissionDeniedError{err}
	case "invalid expression":
		return InvalidExpressionError{err}
	case "invalid url parameter":
		return InvalidUrlParameterError{err}
	case "transaction aborted":
		return TransactionAbortedError{err}
	case "invalid write time":
		return InvalidWriteTimeError{err}
	case "invalid ref":
		return InvalidReferenceError{err}
	case "missing identity":
		return MissingIdentityError{err}
	case "invalid token":
		return InvalidTokenError{err}
	case "stack overflow":
		return StackOverflowError{err}
	case "authentication failed":
		return AuthenticationFailedError{err}
	case "value not found":
		return ValueNotFoundError{err}
	case "instance not found":
		return InstanceNotFoundError{err}
	case "instance already exists":
		return InstanceAlreadyExistsError{err}
	case "validation failed":
		return ValidationFailedError{err}
	case "instance not unique":
		return InstanceNotUniqueError{err}
	case "feature not available":
		return FeatureNotAvailableError{err}
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
