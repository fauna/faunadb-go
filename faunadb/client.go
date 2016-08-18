package faunadb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"faunadb/query"
	"faunadb/values"
	"fmt"
	"io/ioutil"
	"net/http"
)

type resource struct {
	Resource values.Value `fauna:"resource"`
}

type FaunaClient struct {
	Secret   string
	Endpoint string
	Http     http.Client
}

func (client *FaunaClient) Query(expr query.Expr) (values.Value, error) {
	request, err := newRequest(client.Secret, client.Endpoint, expr)
	if err != nil {
		return values.Value{}, err
	}

	response, err := client.Http.Do(request)
	if err != nil {
		return values.Value{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		str, _ := ioutil.ReadAll(response.Body)
		return values.Value{}, fmt.Errorf("Query error %v: %s", response.StatusCode, str)
	}

	fullValueResponse, err := values.ReadValue(response.Body)
	if err != nil {
		return values.Value{}, err
	}

	var res resource
	err = fullValueResponse.Get(&res)
	if err != nil {
		return values.Value{}, err
	}

	return res.Resource, nil
}

func newRequest(secret, endpoint string, expr query.Expr) (*http.Request, error) {
	body, err := json.Marshal(expr)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", basicAuth(secret))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	return request, nil
}

func basicAuth(secret string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(secret))
	return fmt.Sprintf("Basic %s:", encoded)
}
