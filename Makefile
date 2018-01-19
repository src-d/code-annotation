# Package configuration
PROJECT = code-annotation
COMMANDS = server/cmd/code-annotation

# Including ci Makefile
MAKEFILE = Makefile.main
CI_REPOSITORY = https://github.com/src-d/ci.git
CI_FOLDER = .ci

# Tools
YARN = yarn
REMOVE = rm -rf

SERVER_URL ?= /api
API_PORT ?= 8080

$(MAKEFILE):
	@git clone --quiet $(CI_REPOSITORY) $(CI_FOLDER); \
	cp $(CI_FOLDER)/$(MAKEFILE) .;

-include $(MAKEFILE)

# set enviroment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

dependencies-frontend: dependencies
	$(YARN)	install

test-frontend: dependencies-frontend
	$(YARN) test

lint: dependencies-frontend
	$(YARN) lint

build: dependencies-frontend
	REACT_APP_SERVER_URL=$(SERVER_URL) $(YARN) build

## Compiles the dashboard assets, and serve the dashboard through its API
serve: build
	go run server/cmd/code-annotation/*

gorun:
	go run server/cmd/code-annotation/*
