package faunadb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type FaunaClient struct {
	Secret   string
	Endpoint string "https://cloud.faunadb.com"
	Http     http.Client
}

func (client *FaunaClient) Query(expr string) (value string, err error) {
	request, err := newRequest(client.Secret, client.Endpoint, expr)
	if err != nil {
		return
	}

	response, err := client.Http.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = parseResponse(response.Body, &value)
	return
}

func newRequest(secret string, endpoint string, expr string) (*http.Request, error) {
	body := bytes.NewBufferString(fmt.Sprintf("\"%s\"", expr))
	request, err := http.NewRequest("POST", endpoint, body)

	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(secret, "")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	return request, nil
}

func parseResponse(raw io.Reader, result *string) error {
	parsed := new(bytes.Buffer)
	_, err := parsed.ReadFrom(raw)

	if err != nil {
		return err
	}

	*result = parsed.String()
	return nil
}
