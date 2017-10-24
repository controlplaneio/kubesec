SHELL := /bin/bash
PACKAGE = github.com/controlplane/theseus/cmd
HAS_DEP := $(shell command -v dep 2>/dev/null)

.PHONY: build container dep dev local test test test-acceptance test-unit
.SILENT:

all:
	make dep
	make test
	make build

test-go :
	go test $$(go list ./... | grep -v '/vendor/')

test:
	bash -xc ' \
	  cd test && ./bin/bats/bin/bats . \
	'

test-new:
	bash -xc ' \
	  cd test && /usr/src/bats-core/bin/bats . \
	'

up-start:
	bash -xc ' \
		cd up && AWS_PROFILE=binslug-s3 up start \
	'

up-deploy:
	bash -xc ' \
		cd up && AWS_PROFILE=binslug-s3 up deploy production \
	'

up-url:
	bash -xc ' \
		cd up && AWS_PROFILE=binslug-s3 up deploy production \
	'

test-acceptance-old:
	cd test && ./test-acceptance.sh

test-unit-old:
	cd test && ./test-theseus.sh

dev:
	make test && make build

dep-safe:
	bash -xc ' \
		make dep && make test || { \
			echo "Attempting to remedy gopkg.in/yaml.v2"; \
			rm -rf $$(pwd)/vendor/gopkg.in/yaml.v2; \
			go get -v gopkg.in/yaml.v2 && \
				mkdir -p $$(pwd)/vendor/gopkg.in && \
				ln -s $${GOPATH}/src/gopkg.in/yaml.v2 $$(pwd)/vendor/gopkg.in/ && \
				make test; \
		}; \
	'

dep: get-dep
	dep ensure -v

prune:
	dep prune -v

build:
	bash -c ' \
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

cloud:
	cat cloudbuild.yaml
	gcloud container builds submit --config cloudbuild.yaml .

local:
	bash -c "container-builder-local --config cloudbuild.yaml --dryrun=false . 2>&1"

alpine:
	bash -xc ' \
		pwd; ls -lasp; \
		mkdir -p /gocode/src/github.com/controlplane/; \
		ln -s /workspace /gocode/src/github.com/controlplane/theseus; \
		cd /gocode/src/github.com/controlplane/theseus; \
		pwd; \
		ls -lasp; \
		\
		make dep-safe && \
		make build; \
	'

release:
	hub release create \
		-d \
		-a $$(basename $$(pwd)) \
		-m "Version $$(git describe --tags 2>/dev/null)" \
		-m "$$(git log --format=oneline \
			| cut -d' ' -f 2- \
			| awk '!x[$$0]++' \
			| grep -iE '^[^ :]*:' \
			| grep -iEv '^(build|refactor):')" \
		$$(git describe --tags)

get-dep:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
