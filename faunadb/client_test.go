package faunadb

import (
	"faunadb/values"
	"testing"
)

func TestEcho(t *testing.T) {
	client := &FaunaClient{
		Secret:   "secret",
		Endpoint: "http://localhost:8443",
	}

	res, err := client.Query("HI")
	if err != nil {
		t.Error(err)
		return
	}

	expected := values.NewValue("HI")
	if res != expected {
		t.Errorf("Expected: %s got: %s", expected, res)
	}
}
