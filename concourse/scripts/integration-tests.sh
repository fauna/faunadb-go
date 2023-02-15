#!/bin/sh

set -eou

apk add --update make gcc musl-dev

while ! $(curl --output /dev/null --silent --fail --max-time 1 http://faunadb:8443/ping); do sleep 1; done

make test
make coverage
