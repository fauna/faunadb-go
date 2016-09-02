package faunadb

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	emptyErrorBody = `{ "errors": [] }`
	noErrors       = []QueryError{}
)

type fakeBody struct{ io.Reader }

func (f fakeBody) Close() error { return nil }

func TestReturnBadRequestOn400(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(400))
	require.Equal(t, BadRequest{errorResponseWith(400)}, err)
}

func TestReturnUnauthorizedOn401(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(401))
	require.Equal(t, Unauthorized{errorResponseWith(401)}, err)
}

func TestReturnNotFoundOn404(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(404))
	require.Equal(t, NotFound{errorResponseWith(404)}, err)
}

func TestReturnInternalErrorOn500(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(500))
	require.Equal(t, InternalError{errorResponseWith(500)}, err)
}

func TestReturnUnavailableOn503(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(503))
	require.Equal(t, Unavailable{errorResponseWith(503)}, err)
}

func TestReturnUnknownErrorByDefault(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(1001))
	require.Equal(t, UnknownError{errorResponseWith(1001)}, err)
}

func TestParseErrorResponse(t *testing.T) {
	json := `
	{
		"errors": [
			{
				"position": [ "data", "token" ],
				"code": "invalid token",
				"description": "Invalid token.",
				"failures": [
					{
						"field": [ "data", "token" ],
						"code": "invalid token",
						"description": "invalid token"
					}
				]
			}
		]
	}
	`

	expectedError := Unauthorized{
		responseError{
			status: 401,
			errors: []QueryError{
				QueryError{
					Position:    []string{"data", "token"},
					Code:        "invalid token",
					Description: "Invalid token.",
					Failures: []ValidationFailure{
						ValidationFailure{
							Field:       []string{"data", "token"},
							Code:        "invalid token",
							Description: "invalid token",
						},
					},
				},
			},
		},
	}

	response := &http.Response{
		StatusCode: 401,
		Body:       fakeBody{bytes.NewBufferString(json)},
	}

	require.Equal(t, expectedError, checkForResponseErrors(response))
}

func httpErrorResponseWith(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       fakeBody{bytes.NewBufferString(emptyErrorBody)},
	}
}

func errorResponseWith(status int) responseError {
	return responseError{
		status: status,
		errors: noErrors,
	}
}
