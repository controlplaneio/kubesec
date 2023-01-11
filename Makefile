NAME := kubesec
GITHUB_ORG = controlplaneio
DOCKER_HUB_ORG = controlplane
K8S_SCHEMA_VER = 1.25.4

### github.com/controlplaneio/ensure-content.git makefile-header START ###
ifeq ($(NAME),)
  $(error NAME required, please add "NAME := project-name" to top of Makefile)
else ifeq ($(GITHUB_ORG),)
    $(error GITHUB_ORG required, please add "GITHUB_ORG := controlplaneio" to top of Makefile)
else ifeq ($(DOCKER_HUB_ORG),)
    $(error DOCKER_HUB_ORG required, please add "DOCKER_HUB_ORG := controlplane" to top of Makefile)
endif

PKG := github.com/$(GITHUB_ORG)/$(NAME)
DOCKER_REGISTRY_FQDN ?= docker.io
DOCKER_HUB_URL := $(DOCKER_REGISTRY_FQDN)/$(DOCKER_HUB_ORG)/$(NAME)

SHELL := /bin/bash
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GIT_MESSAGE := $(shell git -c log.showSignature=false \
	log --max-count=1 --pretty=format:"%H")
GIT_SHA := $(shell git -c log.showSignature=false rev-parse HEAD)
GIT_TAG := $(shell bash -c 'TAG=$$(git -c log.showSignature=false \
	describe --tags --exact-match --abbrev=0 $(GIT_SHA) 2>/dev/null); echo "$${TAG:-dev}"')
GIT_UNTRACKED_CHANGES := $(shell git -c log.showSignature=false \
	status --porcelain)

ifneq ($(GIT_UNTRACKED_CHANGES),)
  GIT_COMMIT := $(GIT_COMMIT)-dirty
  ifneq ($(GIT_TAG),dev)
    GIT_TAG := $(GIT_TAG)-dirty
  endif
endif

CONTAINER_TAG ?= $(GIT_TAG)
CONTAINER_TAG_LATEST := $(CONTAINER_TAG)
CONTAINER_NAME := $(DOCKER_REGISTRY_FQDN)/$(NAME):$(CONTAINER_TAG)

# if no untracked changes and tag is not dev, release `latest` tag
ifeq ($(GIT_UNTRACKED_CHANGES),)
  ifneq ($(GIT_TAG),dev)
    CONTAINER_TAG_LATEST = latest
  endif
endif

CONTAINER_NAME_LATEST := $(DOCKER_REGISTRY_FQDN)/$(NAME):$(CONTAINER_TAG_LATEST)

# golang buildtime, more at https://github.com/jessfraz/pepper/blob/master/Makefile
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

export NAME DOCKER_HUB_URL BUILD_DATE GIT_MESSAGE GIT_SHA GIT_TAG \
  CONTAINER_TAG CONTAINER_NAME CONTAINER_NAME_LATEST CONTAINER_NAME_TESTING
### github.com/controlplaneio/ensure-content.git makefile-header END ###

PACKAGE = none
BATS_PARALLEL_JOBS := $(shell command -v parallel 2>/dev/null && echo '--jobs 20')
REMOTE_URL ?= "https://v2.kubesec.io/scan"

.PHONY: all
all: help

# ---

.PHONY: all
lint:
	@echo "+ $@"
	-make lint-markdown
	make lint-go-fmt

.PHONY: lint-go-fmt
lint-go-fmt: ## golang fmt check
	@echo "+ $@"
	gofmt -l -s ./pkg | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

MARKDOWN_IMAGE ?= registry.gitlab.com/06kellyjac/docker_markdownlint-cli
MARKDOWN_IMAGE_TAG ?= 0.19.0
.PHONY: lint-markdown
lint-markdown:
	@echo "+ $@"
	docker run -v ${PWD}:/markdown ${MARKDOWN_IMAGE}:${MARKDOWN_IMAGE_TAG} '**/*.md' --ignore 'test/bin/'

# ---

.PHONY: test
test: ## unit and local acceptance tests
	@echo "+ $@"
	make test-unit build test-acceptance

test/bin/%:
	git submodule update --init -- $@

.PHONY: bats
bats: test/bin/bats test/bin/bats-assert test/bin/bats-support ## fetch bats dependencies

.PHONY: test-acceptance
test-acceptance: bats build ## acceptance tests
	@echo "+ $@"
	bash -xc 'cd test && ./bin/bats/bin/bats $(BATS_PARALLEL_JOBS) .'

.PHONY: test-remote
test-remote: bats build ## acceptance tests against remote URL
	@echo "+ $@"
	bash -xc 'cd test && REMOTE_URL=$(REMOTE_URL) ./bin/bats/bin/bats $(BATS_PARALLEL_JOBS) .'

.PHONY: test-unit
test-unit: ## golang unit tests
	@echo "+ $@"
	go test -race $$(go list ./... | grep -v '/vendor/') -run "$${RUN:-.*}"

.PHONY: test-unit-verbose
test-unit-verbose: ## golang unit tests (verbose)
	@echo "+ $@"
	go test -race -v $$(go list ./... | grep -v '/vendor/') -run "$${RUN:-.*}"

# ---

.PHONY: build
build: ## golang build
	@echo "+ $@"
	go build -a -o ./dist/kubesec .

.PHONY: prune
prune: ## golang dependency prune
	@echo "+ $@"
	go mod tidy

# ---

.PHONY: docker-build
docker-build: ## builds a docker image
	@echo "+ $@"
	docker build --tag "${CONTAINER_NAME}" --build-arg "K8S_SCHEMA_VER=${K8S_SCHEMA_VER}" .
	docker tag "${CONTAINER_NAME}" "${CONTAINER_NAME_LATEST}"
	@echo "Successfully tagged ${CONTAINER_NAME} as ${CONTAINER_NAME_LATEST}"

.PHONY: docker-run
docker-run: ## runs the last build docker image
	@echo "+ $@"
	docker run -it "${CONTAINER_NAME}"

.PHONY: docker-push
docker-push: ## pushes the last build docker image
	@echo "+ $@"
	docker push "${CONTAINER_NAME}"
	docker push "${CONTAINER_NAME_LATEST}"

# ---

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
    | grep -Ev '^help\b[[:space:]]*:' \
    | sort \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
