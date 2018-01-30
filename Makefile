# Package configuration
PROJECT = code-annotation
COMMANDS = cli/server
DEPENDENCIES = github.com/golang/dep/cmd/dep

# Including ci Makefile
MAKEFILE = Makefile.main
CI_REPOSITORY = https://github.com/src-d/ci.git
CI_FOLDER = .ci

# Tools
YARN = yarn
REMOVE = rm -rf

HOST ?= 127.0.0.1
PORT ?= 8080
SERVER_URL ?= //$(HOST):$(PORT)

$(MAKEFILE):
	@git clone --quiet $(CI_REPOSITORY) $(CI_FOLDER); \
	cp $(CI_FOLDER)/$(MAKEFILE) .;

-include $(MAKEFILE)

# set enviroment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

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
