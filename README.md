# FaunaDB Go Driver

[![Go Report Card](https://goreportcard.com/badge/github.com/fauna/faunadb-go)](https://goreportcard.com/report/github.com/fauna/faunadb-go)
[![GoDoc](https://godoc.org/github.com/fauna/faunadb-go/faunadb?status.svg)](https://pkg.go.dev/github.com/fauna/faunadb-go/v4)
[![License](https://img.shields.io/badge/license-MPL_2.0-blue.svg?maxAge=2592000)](https://raw.githubusercontent.com/fauna/faunadb-go/main/LICENSE)

A Go lang driver for [FaunaDB](https://fauna.com/).

## Supported Go Versions

Currently, the driver is tested on:
- 1.13
- 1.14
- 1.15
- 1.16

## Using the Driver

### Installing

To get the latest version run:

```bash
go get github.com/fauna/faunadb-go/v4/faunadb
```

Please note that our driver undergoes breaking changes from time to time, so depending on our `main` branch is not recommended.
It is recommended to use one of the following methods instead:

### Importing

For better usage, we recommend that you import this driver with an alias import.

#### Using `dep` or `go get`

To import a specific version when using `go get`, use:

```go
import f "github.com/fauna/faunadb-go/v4/faunadb"
```

### Basic Usage

```go
package main

import (
	"fmt"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

type User struct {
	Name string `fauna:"name"`
}

func main() {
	client := f.NewFaunaClient("your-secret-here")

	res, err := client.Query(f.Get(f.Ref(f.Collection("user"), "42")))
	if err != nil {
		panic(err)
	}

	var user User

	if err := res.At(f.ObjKey("data")).Get(&user); err != nil {
		panic(err)
	}

	fmt.Println(user)
}
```

### Streaming feature usage
```go
package main

import f "github.com/fauna/faunadb-go/v4/faunadb"

func main() {
	secret := ""
	dbClient = f.NewFaunaClient(secret)
	var ref f.RefV
	value, err := dbClient.Query(f.Get(&ref))
	if err != nil {
		panic(err)
	}
	err = value.At(f.ObjKey("ref")).Get(&docRef)
	var subscription f.StreamSubscription
	subscription = dbClient.Stream(docRef)
	err = subscription.Start()
	if err != nil {
		panic("Panic")
	}
	for a := range subscription.StreamEvents() {
		switch a.Type() {
	
		case f.StartEventT:
			// do smth on start event
	
		case f.HistoryRewriteEventT:
			// do smth on historyRewrite event	
			
		case f.VersionEventT:
			// do smth on version event
			
		case f.ErrorEventT:
			// do smth on error event
			subscription.Close() // if you want to close streaming on errors
		}
	}
}
```

### Omitempty usage
```go
package main

import f "github.com/fauna/faunadb-go/v4/faunadb"

func main() {
	secret := ""
	dbClient = f.NewFaunaClient(secret)
	var ref f.RefV
	value, err := dbClient.Query(f.Get(&ref))
	if err != nil {
		panic(err)
	}
	type OmitStruct struct {
		Name           string      `fauna:"name,omitempty"`
		Age            int         `fauna:"age,omitempty"`
		Payment        float64     `fauna:"payment,omitempty"`
		AgePointer     *int        `fauna:"agePointer,omitempty"`
		PaymentPointer *float64    `fauna:"paymentPointer,omitempty"`
	}
	_, err := dbClient.Query(
		f.Create(f.Collection("categories"), f.Obj{"data": OmitStruct{Name: "John", Age: 0}}))
	if err != nil {
		panic(err)
	}
}
```
**Result (empty/zero fields will be ignored):**
```text
{
  "ref": Ref(Collection("categories"), "295143889346494983"),
  "ts": 1617729997710000,
  "data": {
    "name": "John"
  }
}
```
### Http2 support
Driver uses http2 by default. To use http 1.x provide custom http client to `FaunaClient`
```go
package main

import f "github.com/fauna/faunadb-go/v4/faunadb"

func main() {
	secret := ""
	customHttpClient := http.Client{}
	dbClient = f.NewFaunaClient(secret, f.HTTP(&customHttpClient))
}
```
<br>
For more information about Fauna Query Language (FQL), consult our query language
[reference documentation](https://docs.fauna.com/fauna/current/api/fql/).

Specific reference documentation for the driver is hosted at
[GoDoc](https://pkg.go.dev/github.com/fauna/faunadb-go/v4).


Most users found tests for the driver as a very useful form of documentation
[Check it out here](https://github.com/fauna/faunadb-go/tree/main/faunadb).


## Contributing

GitHub pull requests are very welcome.

### Driver Development

Run `go get -t ./...` in order to install project's dependencies.

Run tests against FaunaDB Cloud by passing your root database key to the
test suite, as follows: `FAUNA_ROOT_KEY="your-cloud-secret" go test ./...`.

If you have access to another running FaunaDB database, use the
`FAUNA_ENDPOINT` environment variable to specify its URI.

Alternatively, tests can be run via a Docker container with
`FAUNA_ROOT_KEY="your-cloud-secret" make docker-test` (an alternate
Debian-based Go image can be provided via `RUNTIME_IMAGE`).

Tip: Setting the `FAUNA_QUERY_TIMEOUT_MS` environment variable will
set a timeout in milliseconds for all queries.

## LICENSE

Copyright 2020 [Fauna, Inc.](https://fauna.com/)

Licensed under the Mozilla Public License, Version 2.0 (the
"License"); you may not use this software except in compliance with
the License. You may obtain a copy of the License at

[http://mozilla.org/MPL/2.0/](http://mozilla.org/MPL/2.0/)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.
