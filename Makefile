SHELL := /bin/bash
PACKAGE = github.com/controlplane/theseus/cmd
HAS_DEP := $(shell command -v dep 2>/dev/null)
REMOTE_URL="https://kubesec.io/"

.PHONY: build container dep dev local test test test-acceptance test-unit
.SILENT:

all:
	make dep
	make test
	make build

deploy:
	bash -xec ' \
		unalias make || true; \
		make test \
		&& make up-deploy-staging \
			&& make test-remote-staging \
			&& make up-deploy \
			&& make test-remote ; \
	'
gen-html:
	bash -ec ' \
		(mkdir -p html/basics kubesec.io/contents/basics; \
		IFS="$$(printf "\n+")"; \
		IFS="$${IFS%+}"; \
		for BLOB in $$(cat k8s-rules.json  | jq -c ".rules[]"); do \
			SELECTOR=$$(echo "$$BLOB" | jq .selector -r) \
			FILE_NAME=$$(echo "$$SELECTOR" | sed "s,[^a-zA-Z],-,g" \
							 | sed "s,--*,-,g" \
							 | sed "s,^-,," | sed "s,-$$,,").md; \
			FILE="kubesec.io/content/basics/$${FILE_NAME}"; \
			rm "$${FILE}" || true; \
			touch "$${FILE}" html/basics/"$${FILE_NAME}"; \
			TITLE=$$(echo "$${BLOB}" | jq -r ".reason | select(values)"); \
			SELECTOR_ESCAPED=$${SELECTOR//\"/\\\"}; \
			SELECTOR_ESCAPED=$${SELECTOR//\"/\\\"}; \
			echo "+++" >> "$${FILE}"; \
			echo "title = \"$${SELECTOR_ESCAPED}\"" >> "$${FILE}"; \
			echo "weight = 15" >> "$${FILE}"; \
			echo "+++" >> "$${FILE}"; \
			printf "\n## $${TITLE}\n" >> "$${FILE}"; \
			cat html/basics/$${FILE_NAME} >> $${FILE}; \
		done \
		) \
	'

test:
	bash -xc ' \
	  (cd test && ./bin/bats/bin/bats .) \
	'
test-remote:
	bash -xc ' \
	  (cd test && TEST_REMOTE=1 REMOTE_URL=$(REMOTE_URL) ./bin/bats/bin/bats . )\
	'

test-remote-staging:
	bash -xc ' \
		(REMOTE_URL=$$(\make up-url-staging 2>/dev/null); \
	  cd test && TEST_REMOTE=1 REMOTE_URL=$${REMOTE_URL} ./bin/bats/bin/bats .) \
	'

test-new:
	bash -xc ' \
	  (cd test && /usr/src/bats-core/bin/bats .) \
	'

up-start:
	bash -xc ' \
		(cd up \
		&& AWS_PROFILE=binslug-s3 up run build \
		&& AWS_PROFILE=binslug-s3 up start) \
	'

up-deploy-staging:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up deploy staging) \
	'

up-url-staging:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up url staging) \
	'

up-deploy:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up deploy production) \
	'

up-url:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up url production) \
	'

test-acceptance-old:
	cd test && ./test-acceptance.sh

test-unit-old:
	cd test && ./test-theseus.sh

test-go :
	go test $$(go list ./... | grep -v '/vendor/')

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
