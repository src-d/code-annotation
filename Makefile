# Package configuration
PROJECT = code-annotation
COMMANDS = cli/server
DEPENDENCIES = github.com/golang/dep/cmd/dep

HOST ?= 127.0.0.1
PORT ?= 8080
SERVER_URL ?= //$(HOST):$(PORT)

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --depth 1 $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# set enviroment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

# Tools
YARN = yarn
REMOVE = rm -rf

godep:
	dep ensure

dependencies-frontend: godep
	$(YARN)	install

test-frontend: dependencies-frontend
	$(YARN) test

lint: dependencies-frontend
	$(YARN) lint

build: dependencies-frontend
	REACT_APP_SERVER_URL=$(SERVER_URL) $(YARN) build

## Compiles the dashboard assets, and serve the dashboard through its API
serve: build
	go run cli/server/server.go

gorun:
	go run cli/server/server.go
