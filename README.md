# FaunaDB Go Driver

[![Go Report Card](https://goreportcard.com/badge/github.com/fauna/faunadb-go)](https://goreportcard.com/report/github.com/fauna/faunadb-go)
[![GoDoc](https://godoc.org/github.com/fauna/faunadb-go/faunadb?status.svg)](https://godoc.org/github.com/fauna/faunadb-go/faunadb)
[![License](https://img.shields.io/badge/license-MPL_2.0-blue.svg?maxAge=2592000)](https://raw.githubusercontent.com/fauna/faunadb-go/master/LICENSE)

A Go lang driver for [FaunaDB](https://fauna.com/).

## Supported Go Versions

Currently, the driver is tested on:
- 1.11
- 1.12
- 1.13

## Using the Driver

### Installing

To get the latest version run:

```bash
go get github.com/fauna/faunadb-go/faunadb
```

Please note that our driver undergoes breaking changes from time to time, so depending on our `master` branch is not recommended.
It is recommended to use one of the following methods instead:

#### Using `gopkg.in`

To get a specific version when using `gopkg.in`, use:

```bash
go get gopkg.in/fauna/faunadb-go.v2/faunadb
```

#### Using `dep`

To get a specific version when using `dep`, use:

```bash
dep ensure -add github.com/fauna/faunadb-go/faunadb@v2.11.0
```

### Importing

For better usage, we recommend that you import this driver with an alias import.

#### Using `gopkg.in`

To import a specific version when using `gopkg.in`, use:

```go
import f "gopkg.in/fauna/faunadb-go.v2/faunadb"
```

#### Using `dep` or `go get`

To import a specific version when using `dep` or `go get`, use:

```go
import f "github.com/fauna/faunadb-go/faunadb"
```

### Basic Usage

```go
package main

import (
	"fmt"

	f "github.com/fauna/faunadb-go/faunadb"
)

type User struct {
	Name string `fauna:"name"`
}

func main() {
	client := f.NewFaunaClient("your-secret-here")

	res, err := client.Query(f.Get(f.RefCollection(f.Collection("user"), "42")))
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

For more information about FaunaDB query language, consult our query language
[reference documentation](https://docs.fauna.com/fauna/current/reference/queryapi/).

Specific reference documentation for the driver is hosted at
[GoDoc](https://godoc.org/github.com/fauna/faunadb-go/faunadb).


Most users found tests for the driver as a very useful form of documentation
[Check it out here](https://github.com/fauna/faunadb-go/tree/master/faunadb).


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

## LICENSE

Copyright 2018 [Fauna, Inc.](https://fauna.com/)

Licensed under the Mozilla Public License, Version 2.0 (the
"License"); you may not use this software except in compliance with
the License. You may obtain a copy of the License at

[http://mozilla.org/MPL/2.0/](http://mozilla.org/MPL/2.0/)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.
