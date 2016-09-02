package faunadb

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReturnBadRequestOn400(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 400})
	require.Equal(t, BadRequest{}, err)
}

func TestReturnUnauthorizedOn401(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 401})
	require.Equal(t, Unauthorized{}, err)
}

func TestReturnNotFoundOn404(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 404})
	require.Equal(t, NotFound{}, err)
}

func TestReturnInternalErrorOn500(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 500})
	require.Equal(t, InternalError{}, err)
}

func TestReturnUnavailableOn503(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 503})
	require.Equal(t, Unavailable{}, err)
}
