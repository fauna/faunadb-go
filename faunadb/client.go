package faunadb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
)

type resource struct { //FIXME: Replace this with value transversal once it's implemented
	Resource Value `fauna:"resource"`
}

type FaunaClient struct {
	Secret   string
	Endpoint string
	HTTP     http.Client
}

func (client *FaunaClient) Query(expr Expr) (value Value, err error) {
	response, err := client.performRequest(expr)

	if response != nil {
		defer func() { _ = response.Body.Close() }()
	}

	if err == nil {
		if err = checkForResponseErrors(response); err == nil {
			value, err = client.parseResponse(response)
		}
	}

	return
}

func (client *FaunaClient) performRequest(expr Expr) (response *http.Response, err error) {
	var request *http.Request

	if request, err = client.prepareRequest(expr); err == nil {
		response, err = client.HTTP.Do(request)
	}

	return
}

func (client *FaunaClient) prepareRequest(expr Expr) (request *http.Request, err error) {
	var body []byte

	if body, err = writeJSON(expr); err == nil {
		if request, err = http.NewRequest("POST", client.Endpoint, bytes.NewReader(body)); err == nil {
			request.Header.Add("Authorization", client.basicAuth())
			request.Header.Add("Content-Type", "application/json; charset=utf-8")
		}
	}

	return
}

func (client *FaunaClient) basicAuth() string {
	encoded := base64.StdEncoding.EncodeToString([]byte(client.Secret))
	return fmt.Sprintf("Basic %s:", encoded)
}

func (client *FaunaClient) parseResponse(response *http.Response) (Value, error) {
	value, err := parseJSON(response.Body)

	if err != nil {
		return nil, err
	}

	var res resource

	if err := value.To(&res); err != nil {
		return nil, err
	}

	return res.Resource, nil
}
