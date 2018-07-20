SHELL := /bin/bash
PACKAGE = none
HAS_DEP := $(shell command -v dep 2>/dev/null)
REMOTE_URL ?="https://kubesec.io/"

.PHONY: build container dep dev local test test test-acceptance test-unit
.SILENT:

all:
	help

all-go: ## golang toolchain
	make dep
	make test
	make build

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
hugo:
	bash -ec ' \
		( \
		(IS_KUBESEC=$$(wmctrl -l -x | grep -q kubesec.io && echo 1 || echo 0); cd kubesec.io && while read LINE; do echo "$${LINE}"; if [[ "$${IS_KUBESEC}" != 1 ]] && [[ "$${LINE}" =~ ^Web\ Server\ is\ available\ at ]]; then echo "$${LINE}" | sed -E "s,.*(localhost[^ ]*).*,\1,g" | xargs -I{} xdg-open "http://{}" || true; fi; done < <(hugo server --disableFastRender)); \
		)'

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

logs:
	bash -xc ' \
		(cd up && AWS_PROFILE=binslug-s3 up logs -f) \
	'

test:
	bash -xc ' \
	  (COMMAND=./bin/bats/bin/bats; \
	  cd test && if command -v bats; then COMMAND=bats; fi && $${COMMAND} $(FLAGS) .) \
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

dep:
	command -v up &>/dev/null || curl -sfL https://raw.githubusercontent.com/apex/up/master/install.sh | sh

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

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
    | grep -Ev '^help\b[[:space:]]*:' \
    | sort \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
