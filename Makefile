RUNTIME_IMAGE ?= golang:1.8
DOCKER_RUN_FLAGS = -it --rm

ifdef FAUNA_ROOT_KEY
DOCKER_RUN_FLAGS += -e FAUNA_ROOT_KEY=$(FAUNA_ROOT_KEY)
endif

ifdef FAUNA_ENDPOINT
DOCKER_RUN_FLAGS += -e FAUNA_ENDPOINT=$(FAUNA_ENDPOINT)
endif

test:
	go test -v ./...

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

jenkins-test:
	go test -v -race -coverprofile=results/coverage.txt -covermode=atomic ./... 2>&1 | tee log.txt
	go-junit-report -package-name faunadb -set-exit-code < log.txt > results/report.xml

docker-wait:
	dockerize -wait $(FAUNA_ENDPOINT)/ping -timeout $(FAUNA_TIMEOUT)

docker-test:
	docker build -f Dockerfile.test -t faunadb-go-test:latest --build-arg RUNTIME_IMAGE=$(RUNTIME_IMAGE) .
	docker run $(DOCKER_RUN_FLAGS) faunadb-go-test:latest
