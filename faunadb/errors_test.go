package faunadb

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	emptyErrorBody = `{ "errors": [] }`
	noErrors       = []QueryError{}
)

func TestReturnBadRequestOn400(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(400, emptyErrorBody))
	require.Equal(t, BadRequest{errorResponseWith(400, noErrors)}, err)
}

func TestReturnUnauthorizedOn401(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(401, emptyErrorBody))
	require.Equal(t, Unauthorized{errorResponseWith(401, noErrors)}, err)
}

func TestReturnNotFoundOn404(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(404, emptyErrorBody))
	require.Equal(t, NotFound{errorResponseWith(404, noErrors)}, err)
}

func TestReturnInternalErrorOn500(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(500, emptyErrorBody))
	require.Equal(t, InternalError{errorResponseWith(500, noErrors)}, err)
}

func TestReturnUnavailableOn503(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(503, emptyErrorBody))
	require.Equal(t, Unavailable{errorResponseWith(503, noErrors)}, err)
}

func TestReturnUnknownErrorByDefault(t *testing.T) {
	err := checkForResponseErrors(httpErrorResponseWith(1001, emptyErrorBody))
	require.Equal(t, UnknownError{errorResponseWith(1001, noErrors)}, err)
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

	err := checkForResponseErrors(httpErrorResponseWith(401, json))

	expectedError := Unauthorized{
		errorResponseWith(401,
			[]QueryError{
				{
					Position:    []string{"data", "token"},
					Code:        "invalid token",
					Description: "Invalid token.",
					Failures: []ValidationFailure{
						{
							Field:       []string{"data", "token"},
							Code:        "invalid token",
							Description: "invalid token",
						},
					},
				},
			},
		),
	}

	require.Equal(t, expectedError, err)
	require.EqualError(t, err, "Response error 401. Errors: [data/token](invalid token): Invalid token.")
}

func TestUnparseableResponse(t *testing.T) {
	json := "can't parse this as an error"
	err := checkForResponseErrors(httpErrorResponseWith(503, json))

	require.Equal(t, Unavailable{errorResponse{status: 503}}, err)
	require.EqualError(t, err, "Response error 503. Unparseable server response.")
}

func httpErrorResponseWith(status int, errorBody string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewBufferString(errorBody)),
	}
}

func errorResponseWith(status int, errors []QueryError) errorResponse {
	return errorResponse{
		parseable: true,
		status:    status,
		errors:    errors,
	}
}
