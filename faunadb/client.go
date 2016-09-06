package faunadb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type FaunaClient struct {
	Secret   string
	Endpoint string
	HTTP     http.Client
}

func (client *FaunaClient) Query(expr Expr) (value Value, err error) {
	response, err := client.performRequest(expr)

	if response != nil {
		defer func() {
			_, _ = io.Copy(ioutil.Discard, response.Body) // Discard remaining bytes so the connection can be reused
			_ = response.Body.Close()
		}()
	}

	if err == nil {
		if err = checkForResponseErrors(response); err == nil {
			value, err = client.parseResponse(response)
		}
	}

	return
}

func (client *FaunaClient) BatchQuery(exprs []Expr) (values []Value, err error) {
	arr := make(Arr, len(exprs))

	for i, expr := range exprs {
		arr[i] = expr
	}

	var res Value

	if res, err = client.Query(arr); err == nil {
		err = res.Get(&values)
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

	return value.At(ObjKey("resource")).GetValue()
}
