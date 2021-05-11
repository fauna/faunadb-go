#!/bin/sh

set -eou

apk add --update make gcc musl-dev

make test
make coverage
