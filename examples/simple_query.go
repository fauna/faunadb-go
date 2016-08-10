package main

import (
	"faunadb"
	"fmt"
)

func main() {
	client := &faunadb.FaunaClient{Secret: "secret", Endpoint: "http://localhost:8443"}
	res, _ := client.Query("HI")
	fmt.Println(res)
}
