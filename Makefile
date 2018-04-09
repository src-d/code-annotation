# Package configuration
PROJECT = code-annotation
COMMANDS = cli/server
DEPENDENCIES = \
	github.com/golang/dep/cmd/dep \
	github.com/jteeuwen/go-bindata \
	github.com/golang/lint/golint
GO_LINTABLE_PACKAGES = $(shell go list ./... | grep -v '/vendor/')

YARN_PRODUCTION ?= true

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
CI_BRANCH ?= v1
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	@git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# Set enviroment variables from .env file
DOT_ENV ?= .env
-include $(DOT_ENV)
export $(shell [ -f "$(DOT_ENV)" ] && sed 's/=.*//' $(DOT_ENV))


# Frontend

dependencies-frontend-development:
	$(MAKE) dependencies-frontend YARN_PRODUCTION=false

dependencies-frontend:
	$(YARN) install --production=$(YARN_PRODUCTION)

test-frontend: dependencies-frontend-development
	$(YARN) test

lint-frontend: dependencies-frontend-development
	$(YARN) lint

build-frontend: dependencies-frontend
	$(YARN) build

dev-frontend: dependencies-frontend
	$(YARN) start

# Backend

dependencies-backend: $(DEPENDENCIES)
	$(GODEP) ensure

build-backend: dependencies-backend

$(GO_LINTABLE_PACKAGES):
	$(GOLINT) $@
	$(GOVET) $@

lint-backend: dependencies-backend $(GO_LINTABLE_PACKAGES)

bindata:
	$(BINDATA) \
		-pkg assets \
		-o ./server/assets/asset.go \
		build/static/... \
		build/*.json \
		build/*.png \
		build/*.svg \
		build/*.ico \
		build/index.html

prepare-build: | build-frontend build-backend bindata

validate-commit: | dependencies-backend no-changes-in-commit

build-app: | prepare-build packages

# Run only server
gorun:
	go run cli/server/server.go

## Compiles the assets, and serve the tool through its API
serve: | build-frontend build-backend gorun

.PHONY: dependencies-frontend build-frontend dev-frontend \
		dependencies-frontend-development prepare-build build-app \
		test-frontend lint-frontend \
		dependencies-backend build-backend release-build \
		lint-backend bindata \
		gorun serve validate-commit
