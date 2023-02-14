#!/bin/sh

set -eou

apk add --update make gcc musl-dev

make docker-wait
make test
make coverage
