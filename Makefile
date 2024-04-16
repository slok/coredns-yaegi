
SHELL := $(shell which bash)
OSTYPE := $(shell uname)
DOCKER := $(shell command -v docker)
GID := $(shell id -g)
UID := $(shell id -u)
VERSION ?= $(shell git describe --tags --always)

UNIT_TEST_CMD := ./scripts/check/unit-test.sh
CHECK_CMD := ./scripts/check/check.sh

DEV_IMAGE_NAME := local/coredns-yaegi

DOCKER_RUN_CMD := docker run --env ostype=$(OSTYPE) -v ${PWD}:/src --rm ${DEV_IMAGE_NAME}
BUILD_DEV_IMAGE_CMD := IMAGE=${DEV_IMAGE_NAME} DOCKER_FILE_PATH=./docker/dev/Dockerfile VERSION=latest ./scripts/build/docker/build-image-dev.sh

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

.PHONY: default
default: help

.PHONY: build-dev-image
build-dev-image:  ## Builds the development docker image.
	@$(BUILD_DEV_IMAGE_CMD)

build: build-dev-image ## Builds the production binary.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_CMD)'

build-all: build-dev-image ## Builds all archs production binaries.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_ALL_CMD)'

.PHONY: test
test: build-dev-image  ## Runs unit test.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(UNIT_TEST_CMD)'

.PHONY: check
check: build-dev-image  ## Runs checks.
	@$(DOCKER_RUN_CMD) /bin/sh -c '$(CHECK_CMD)'

.PHONY: go-gen
go-gen: build-dev-image  ## Generates go based code.
	@$(DOCKER_RUN_CMD) /bin/sh -c './scripts/gogen.sh'

.PHONY: gen
gen: go-gen ## Generates all.

.PHONY: deps
deps:  ## Fixes the dependencies
	@$(DOCKER_RUN_CMD) /bin/sh -c './scripts/deps.sh'

.PHONY: ci-build
ci-build: ## Builds the production binary in CI environment (without docker).
	@$(BUILD_BINARY_CMD)

.PHONY: ci-unit-test
ci-test:  ## Runs unit test in CI environment (without docker).
	@$(UNIT_TEST_CMD)

.PHONY: ci-check
ci-check:  ## Runs checks in CI environment (without docker).
	@$(CHECK_CMD)

