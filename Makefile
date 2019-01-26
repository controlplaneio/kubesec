NAME := kubesec
GITHUB_ORG = controlplaneio
DOCKER_HUB_ORG = controlplane

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
CONTAINER_NAME := $(REGISTRY)/$(NAME):$(CONTAINER_TAG)

# if no untracked changes and tag is not dev, release `latest` tag
ifeq ($(GIT_UNTRACKED_CHANGES),)
  ifneq ($(GIT_TAG),dev)
    CONTAINER_TAG_LATEST = latest
  endif
endif

CONTAINER_NAME_LATEST := $(REGISTRY)/$(NAME):$(CONTAINER_TAG_LATEST)

# golang buildtime, more at https://github.com/jessfraz/pepper/blob/master/Makefile
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

export NAME DOCKER_HUB_URL BUILD_DATE GIT_MESSAGE GIT_SHA GIT_TAG \
  CONTAINER_TAG CONTAINER_NAME CONTAINER_NAME_LATEST CONTAINER_NAME_TESTING
### github.com/controlplaneio/ensure-content.git makefile-header END ###

PACKAGE = none
HAS_DEP := $(shell command -v dep 2>/dev/null)
REMOTE_URL ?="https://kubesec.io/"

.PHONY: all
all: help

.PHONY: go
go: ## golang toolchain
	make test-go
	make build

.PHONY: test-go
test-go: ## golang unit tests
	go test $$(go list ./... | grep -v '/vendor/')

.PHONY: test-go-acceptance
test-go-acceptance: ## acceptance tests targeting golang build
	export BIN_UNDER_TEST=kube-sec-check; make test

.PHONY: dep
dep: ## golang and deployment dependencies
	command -v up &>/dev/null || curl -sfL https://raw.githubusercontent.com/apex/up/master/install.sh | sudo sh
	dep ensure -v

.PHONY: prune
prune: ## golang dependency prune
	dep prune -v

.PHONY: build
build: ## golang build
	bash -xc ' \
		PACKAGE="$(PACKAGE)"; \
		STATUS=$$(git diff-index --quiet HEAD 2>/dev/null || echo "-dirty"); \
		HASH="$$(git rev-parse --short HEAD 2>/dev/null)"; \
		VERSION="$$(git describe --tags 2>/dev/null|| echo $${HASH})$${STATUS}"; \
		go build -ldflags "\
			-X $${PACKAGE}.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p') \
			-X $${PACKAGE}.gitHash=$${HASH} \
			-X $${PACKAGE}.buildVersion=$${VERSION} \
		"; \
	'

.PHONY: dev
dev: ## non-golang dev
	make test && make build

# --- deployment recipes

.PHONY: deploy
deploy: ## deploy, test, promote to prod
	bash -xec ' \
		unalias make || true; \
		time make test \
      && time make dep \
      && time make up-deploy-staging \
			&& time make test-remote-staging \
			&& time make up-deploy \
			&& time make test-remote ; \
	'
.PHONY: hugo
hugo:
	bash -ec ' \
		( \
		(IS_KUBESEC=$$(wmctrl -l -x | grep -q kubesec.io && echo 1 || echo 0); cd kubesec.io && while read LINE; do echo "$${LINE}"; if [[ "$${IS_KUBESEC}" != 1 ]] && [[ "$${LINE}" =~ ^Web\ Server\ is\ available\ at ]]; then echo "$${LINE}" | sed -E "s,.*(localhost[^ ]*).*,\1,g" | xargs -I{} xdg-open "http://{}" || true; fi; done < <(hugo server --disableFastRender)); \
		)'

.PHONY: gen-html
gen-html:
	bash -ec ' \
		( \
		OUTPUT_PATH="kubesec.io/content/basics"; \
		INPUT_PATH="kubesec.io/_content/basics"; \
		mkdir -p "$${INPUT_PATH}" "$${OUTPUT_PATH}"; \
		find "$${OUTPUT_PATH}" -type f | grep -v index.md | xargs --no-run-if-empty rm; \
		IFS="$$(printf "\n+")"; \
		IFS="$${IFS%+}"; \
		for BLOB in $$(cat k8s-rules.json  | jq -c ".rules[]"); do \
			SELECTOR=$$(echo "$$BLOB" | jq -r "(select(.title!=null) | .title), (select(.title==null) | .selector)") \
			FILE_NAME=$$(echo "$$SELECTOR" | sed "s,[^a-zA-Z],-,g" \
							 | sed "s,--*,-,g" \
							 | sed "s,^-,," | sed "s,-$$,,").md; \
			FILE="$${OUTPUT_PATH}/$${FILE_NAME}"; \
			rm "$${FILE}" 2>/dev/null || true; \
			touch "$${FILE}" "$${INPUT_PATH}"/"$${FILE_NAME}"; \
			WEIGHT=$$(echo "$${BLOB}" | jq -r ".weight | select(values)"); \
			TITLE=$$(echo "$${BLOB}" | jq -r ".reason | select(values)"); \
			SELECTOR_ESCAPED=$${SELECTOR//\"/\\\"}; \
			SELECTOR_ESCAPED=$${SELECTOR//\"/\\\"}; \
			echo "+++" >> "$${FILE}"; \
			echo "title = \"$${SELECTOR_ESCAPED}\"" >> "$${FILE}"; \
			echo "weight = $${WEIGHT:-2}" >> "$${FILE}"; \
			echo "+++" >> "$${FILE}"; \
			printf "\n## $${TITLE}\n\n" >> "$${FILE}"; \
			cat "$${INPUT_PATH}"/"$${FILE_NAME}" >> "$${FILE}"; \
			printf "\n\n{{%% katacoda %%}}\n" >> "$${FILE}"; \
		done \
		); \
  touch . \
	'

.PHONY: logs
logs:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up logs -f) \
	'

.PHONY: test
test:
	bash -xc ' \
	  (COMMAND=./bin/bats/bin/bats; \
	  cd test && if command -v bats; then COMMAND=bats; fi && $${COMMAND} $(FLAGS) .) \
	'

.PHONY: test-remote
test-remote:
	bash -xc ' \
	  (cd test && TEST_REMOTE=1 REMOTE_URL=$(REMOTE_URL) ./bin/bats/bin/bats . )\
	'

.PHONY: test-remote-staging
test-remote-staging:
	bash -xc ' \
		(REMOTE_URL=$$(\make up-url-staging 2>/dev/null); \
	  cd test && TEST_REMOTE=1 REMOTE_URL=$${REMOTE_URL} ./bin/bats/bin/bats .) \
	'

.PHONY: test-new
test-new:
	bash -xc ' \
	  (cd test && /usr/src/bats-core/bin/bats .) \
	'

.PHONY: up-start
up-start:
	bash -xc ' \
		(cd up \
		&& AWS_PROFILE=binslug-s3 up run build \
		&& AWS_PROFILE=binslug-s3 up start) \
	'

.PHONY: up-deploy-staging
up-deploy-staging:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up deploy staging) \
	'

.PHONY: up-url-staging
up-url-staging:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up url --stage staging) \
	'

.PHONY: up-deploy
up-deploy:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up deploy production) \
	'

.PHONY: up-url
up-url:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up url --stage production) \
	'

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
    | grep -Ev '^help\b[[:space:]]*:' \
    | sort \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
