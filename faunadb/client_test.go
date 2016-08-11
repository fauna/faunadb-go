package faunadb

import "testing"

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

	expected := "{\"resource\":\"HI\"}"
	if res != expected {
		t.Errorf("Expected: %s got: %s", expected, res)
	}
}
