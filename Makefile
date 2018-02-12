# Package configuration
PROJECT = code-annotation
COMMANDS = cli/server
DEPENDENCIES = github.com/golang/dep/cmd/dep github.com/jteeuwen/go-bindata

HOST ?= 127.0.0.1
PORT ?= 8080
REACT_APP_SERVER_URL ?= //$(HOST):$(PORT) # frontend uses $(REACT_APP_SERVER_URL) as backend
UI_DOMAIN ?= $(REACT_APP_SERVER_URL) # /oauth-callback redirects to $(UI_DOMAIN)/?token=__TOKEN__

# Tools
YARN = yarn
GODEP = dep
GOLINT = golint
GOVET = go vet
BINDATA = go-bindata

# ci variables
TRAVIS_BUILD_DIR ?= $(shell pwd)
PKG_OS = linux
DOCKER_OS = linux
DOCKER_ARCH = amd64

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	@git clone --quiet --depth 1 -b v1 $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# Set enviroment variables from .env file
ENV ?= .env
-include $(ENV)
export $(shell [ -f "$(ENV)" ] && sed 's/=.*//' $(ENV))


# Frontend

dependencies-frontend:
	$(YARN)	install

test-frontend: dependencies-frontend
	$(YARN) test

lint-frontend: dependencies-frontend
	$(YARN) lint

build-frontend: dependencies-frontend
	$(YARN) build

dev-frontend: dependencies-frontend
	$(YARN) start

# Backend

dependencies-backend: $(DEPENDENCIES)
	$(GODEP) ensure

build-backend: dependencies-backend

lint-backend: dependencies-backend
	$(GOLINT) ./server/...
	$(GOVET) ./server/...
	
bindata:
	$(BINDATA) -o ./server/assets/asset.go -pkg assets build/static/... build/*.json build/*.png build/index.html

prepare-build: | build-frontend build-backend bindata

# Run only server
gorun:
	go run cli/server/server.go

## Compiles the assets, and serve the tool through its API
serve: build-frontend build-backend gorun

.PHONY: dependencies-frontend build-frontend dev-frontend \
		test-frontend lint-frontend \
		dependencies-backend build-backend release-build \
		lint-backend bindata \
		gorun serve
