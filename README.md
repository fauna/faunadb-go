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
- 1.8

## Using the Driver

### Installing

```bash
go get github.com/fauna/faunadb-go/faunadb
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

Tests can also be run via a Docker container. The `Dockerfile.test` file
creates an image that automatically runs the tests against a target FaunaDB
cluster. You can run tests this way via the following steps:

1. Build the image via
   `docker build -f Dockerfile.test --build-arg RUNTIME_IMAGE=<image> .`, where
   `<image>` is an Alpine-based image containing Go (for example,
   `golang:1.8-alpine`).
2. Run the image with a given endpoint (`FAUNA_ENDPOINT`) and secret
   (`FAUNA_ROOT_KEY`). By default, the endpoint is `https://db.fauna.com` and
   the secret is `secret`.
   To run tests against cloud,
   use `docker run -it -e FAUNA_ROOT_KEY=<your-cloud-secret> <built-image>`,
   where `<your-cloud-secret>` is a valid admin key secret for cloud, and
   `<built-image>` is the image id produced from building the Docker image.

An example of this build and run process:

```
$ docker build -f Dockerfile.test --build-arg RUNTIME_IMAGE=golang:1.8-alpine .
Sending build context to Docker daemon  146.4kB
... docker image builds ...
Successfully built 1438a4dc32b6
$ docker run -it -e FAUNA_ROOT_KEY="a-cloud-secret" 1438a4dc32b6
2017/11/27 18:44:10 Waiting for: https://db.fauna.com/ping
2017/11/27 18:44:11 Received 200 from https://db.fauna.com/ping
... go tests run ...
PASS
ok  	github.com/fauna/faunadb-go/faunadb	20.411s
2017/11/27 18:44:32 Command finished successfully.
```

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
