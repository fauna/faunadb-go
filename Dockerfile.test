# Take in a runtime image to use for the base system
# Expects a debian image
ARG RUNTIME_IMAGE

# Use the docker image provided via build arg
FROM $RUNTIME_IMAGE

# Copy in the dockerize utility
ARG DOCKERIZE_VERSION=0.6.0
RUN curl -sL https://github.com/jwilder/dockerize/releases/download/v$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-v$DOCKERIZE_VERSION.tar.gz | tar -xzC /usr/local/bin

# Copy project into the image
COPY . /go/src/github.com/fauna/faunadb-go

# Shift over to the project and fetch dependencies
WORKDIR /go/src/github.com/fauna/faunadb-go
RUN go get -t -v ./... github.com/jstemmer/go-junit-report

# Define the default variables for the tests
ENV FAUNA_ROOT_KEY=secret FAUNA_ENDPOINT=https://db.fauna.com FAUNA_TIMEOUT=30s

# Run the tests (after target database is up)
CMD ["make", "docker-wait", "test"]
