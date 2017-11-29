RUNTIME_IMAGE ?= golang:1.8-alpine
DOCKER_RUN_FLAGS = -it --rm

ifdef FAUNA_ROOT_KEY
DOCKER_RUN_FLAGS += -e FAUNA_ROOT_KEY=$(FAUNA_ROOT_KEY)
endif

ifdef FAUNA_ENDPOINT
DOCKER_RUN_FLAGS += -e FAUNA_ENDPOINT=$(FAUNA_ENDPOINT)
endif

docker-test:
	docker build -f Dockerfile.test -t faunadb-go-test:latest --build-arg RUNTIME_IMAGE=$(RUNTIME_IMAGE) .
	docker run $(DOCKER_RUN_FLAGS) faunadb-go-test:latest
