# Name of this service/application
SERVICE_NAME := brigade-exporter

# Path of the go service inside docker
DOCKER_GO_SERVICE_PATH := /go/src/github.com/slok/brigade-exporter

# Shell to use for running scripts
SHELL := $(shell which bash)

# Get docker path or an empty string
DOCKER := $(shell command -v docker)

# Get docker-compose path or an empty string
DOCKER_COMPOSE := $(shell command -v docker-compose)

# Get the main unix group for the user running make (to be used by docker-compose later)
GID := $(shell id -g)

# Get the unix user id for the user running make (to be used by docker-compose later)
UID := $(shell id -u)

# Version from Git.
VERSION=$(shell git describe --tags --always)

# Dev direcotry has all the required dev files.
DEV_DIR := ./docker/dev

# cmds
UNIT_TEST_CMD := ./hack/scripts/unit-test.sh
INTEGRATION_TEST_CMD := ./hack/scripts/integration-test.sh
MOCKS_CMD := ./hack/scripts/mockgen.sh
DOCKER_RUN_CMD := docker run -v ${PWD}:$(DOCKER_GO_SERVICE_PATH) --rm -it $(SERVICE_NAME)
BUILD_BINARY_CMD := VERSION=${VERSION} ./hack/scripts/build-binary.sh
BUILD_IMAGE_CMD := VERSION=${VERSION} ./hack/scripts/build-image.sh
DEPS_CMD := GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor
DEBUG_CMD := go run ./cmd/brigade-exporter/* --debug
DEV_CMD := $(DEBUG_CMD) --development
FAKE_CMD := $(DEV_CMD) --fake
K8S_VERSION := "1.10.7"
SET_K8S_DEPS_CMD := go mod edit \
	-require=k8s.io/client-go@kubernetes-${K8S_VERSION} \
	-require=k8s.io/apimachinery@kubernetes-${K8S_VERSION} \
	-require=k8s.io/api@kubernetes-${K8S_VERSION} && \
	go mod tidy

# environment dirs
DEV_DIR := docker/dev

# The default action of this Makefile is to build the development docker image
.PHONY: default
default: build

# Test if the dependencies we need to run this Makefile are installed
.PHONY: deps-development
deps-development:
ifndef DOCKER
	@echo "Docker is not available. Please install docker"
	@exit 1
endif
ifndef DOCKER_COMPOSE
	@echo "docker-compose is not available. Please install docker-compose"
	@exit 1
endif

# Build the development docker images
.PHONY: build
build:
	docker build -t $(SERVICE_NAME) --build-arg uid=$(UID) --build-arg  gid=$(GID) -f $(DEV_DIR)/Dockerfile .

# run the development stack.
.PHONY: shell
stack: deps-development
	cd $(DEV_DIR) && \
    ( docker-compose -p $(SERVICE_NAME) up --build; \
      docker-compose -p $(SERVICE_NAME) stop; \
      docker-compose -p $(SERVICE_NAME) rm -f; )

# Build production stuff.
build-binary:
	$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_CMD)'

.PHONY: build-image
build-image:
	$(BUILD_IMAGE_CMD)

# Dependencies stuff.
.PHONY: set-k8s-deps
set-k8s-deps:
	$(SET_K8S_DEPS_CMD)

.PHONY: deps
deps:
	$(DEPS_CMD)

# Test stuff in dev
.PHONY: unit-test
unit-test: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(UNIT_TEST_CMD)'
.PHONY: integration-test
integration-test: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(INTEGRATION_TEST_CMD)'
.PHONY: test
test: integration-test

# Mocks stuff in dev
.PHONY: mocks
mocks: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(MOCKS_CMD)'

.PHONY: dev
dev:
	$(DEV_CMD)

.PHONY: fake
fake:
	$(FAKE_CMD)

.PHONY: push
push: export PUSH_IMAGE=true
push: build-image

