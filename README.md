# FaunaDB Go Driver

[![Build Status](https://travis-ci.org/fauna/faunadb-go.svg?branch=master)](https://travis-ci.org/fauna/faunadb-go)
[![Coverage Status](https://codecov.io/gh/fauna/faunadb-go/branch/master/graph/badge.svg)](https://codecov.io/gh/fauna/faunadb-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/fauna/faunadb-go)](https://goreportcard.com/report/github.com/fauna/faunadb-go)
[![GoDoc](https://godoc.org/github.com/fauna/faunadb-go/faunadb?status.svg)](https://godoc.org/github.com/fauna/faunadb-go/faunadb)
[![License](https://img.shields.io/badge/license-MPL_2.0-blue.svg?maxAge=2592000)](https://raw.githubusercontent.com/fauna/faunadb-go/master/LICENSE)

A Go lang driver for [FaunaDB](https://fauna.com/).

## Supported Go Versions

Currently, the driver is tested on:
- 1.5
- 1.6
- 1.7

## Using the Driver

### Installing

```bash
go get github.com/fauna/faunadb-go
```

### Importing

For better usage, we recommend that you import this driver with an alias import
such as:

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

	res, err := client.Query(f.Get(f.Ref("classes/users/42")))
	if err != nil {
		panic(err)
	}

	var user User

	if err := res.Get(&user); err != nil {
		panic(err)
	}

	fmt.Println(user)
}
```

The [tutorials](https://fauna.com/tutorials) in the FaunaDB documentation
contain driver-specific examples.

For more information about FaunaDB query language, consult our query language
[reference documentation](https://fauna.com/documentation/queries).

Specific reference documentation for the driver is hosted at
[GoDoc](https://godoc.org/github.com/fauna/faunadb-go/faunadb).

## Contributing

GitHub pull requests are very welcome.

### Driver Development

Run `go get -t ./...` in order to install project's dependencies.

Run tests with `FAUNA_ROOT_KEY="your-cloud-secret" go test ./...`.

## LICENSE

Copyright 2017 [Fauna, Inc.](https://fauna.com/)

Licensed under the Mozilla Public License, Version 2.0 (the
"License"); you may not use this software except in compliance with
the License. You may obtain a copy of the License at

[http://mozilla.org/MPL/2.0/](http://mozilla.org/MPL/2.0/)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.
