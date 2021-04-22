#
# Copyright 2016 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# set default shell
SHELL = bash -e -o pipefail

# Variables
VERSION                  ?= $(shell cat ./VERSION)

DOCKER_LABEL_VCS_DIRTY     = false
ifneq ($(shell git ls-files --others --modified --exclude-standard 2>/dev/null | wc -l | sed -e 's/ //g'),0)
    DOCKER_LABEL_VCS_DIRTY = true
endif
## Docker related
DOCKER_EXTRA_ARGS        ?=
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${VERSION}
OPENDEVMGR_IMAGENAME     := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}opendevice-manager:${DOCKER_TAG}
TYPE                     ?= minimal

## Docker labels. Only set ref and commit date if committed
DOCKER_LABEL_VCS_URL       ?= $(shell git remote get-url $(shell git remote))
DOCKER_LABEL_VCS_REF       = $(shell git rev-parse HEAD)
DOCKER_LABEL_BUILD_DATE    ?= $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
DOCKER_LABEL_COMMIT_DATE   = $(shell git show -s --format=%cd --date=iso-strict HEAD)

DOCKER_BUILD_ARGS ?= \
	${DOCKER_EXTRA_ARGS} \
	--build-arg org_label_schema_version="${VERSION}" \
	--build-arg org_label_schema_vcs_url="${DOCKER_LABEL_VCS_URL}" \
	--build-arg org_label_schema_vcs_ref="${DOCKER_LABEL_VCS_REF}" \
	--build-arg org_label_schema_build_date="${DOCKER_LABEL_BUILD_DATE}" \
	--build-arg org_opencord_vcs_commit_date="${DOCKER_LABEL_COMMIT_DATE}" \
	--build-arg org_opencord_vcs_dirty="${DOCKER_LABEL_VCS_DIRTY}"

# tool containers
VOLTHA_TOOLS_VERSION ?= 2.3.1

GO                = docker run --rm --user $$(id -u):$$(id -g) -v ${CURDIR}:/app $(shell test -t 0 && echo "-it") -v gocache:/.cache -v gocache-${VOLTHA_TOOLS_VERSION}:/go/pkg voltha/voltha-ci-tools:${VOLTHA_TOOLS_VERSION}-golang go
GO_JUNIT_REPORT   = docker run --rm --user $$(id -u):$$(id -g) -v ${CURDIR}:/app -i voltha/voltha-ci-tools:${VOLTHA_TOOLS_VERSION}-go-junit-report go-junit-report
GOCOVER_COBERTURA = docker run --rm --user $$(id -u):$$(id -g) -v ${CURDIR}:/app/src/github.com/opencord/opendevice-manager -i voltha/voltha-ci-tools:${VOLTHA_TOOLS_VERSION}-gocover-cobertura gocover-cobertura
GOLANGCI_LINT     = docker run --rm --user $$(id -u):$$(id -g) -v ${CURDIR}:/app $(shell test -t 0 && echo "-it") -v gocache:/.cache -v gocache-${VOLTHA_TOOLS_VERSION}:/go/pkg voltha/voltha-ci-tools:${VOLTHA_TOOLS_VERSION}-golangci-lint golangci-lint
HADOLINT          = docker run --rm --user $$(id -u):$$(id -g) -v ${CURDIR}:/app $(shell test -t 0 && echo "-it") voltha/voltha-ci-tools:${VOLTHA_TOOLS_VERSION}-hadolint hadolint

.PHONY: docker-build local-protos local-lib-go help
.DEFAULT_GOAL := help

## Local Development Helpers
local-protos: ## Copies a local version of the dmi-protos dependency into the vendor directory
ifdef LOCAL_PROTOS
	rm -rf vendor/github.com/opencord/device-management-interface/go
	mkdir -p vendor/github.com/opencord/device-management-interface/go
	cp -r ${LOCAL_PROTOS}/go/* vendor/github.com/opencord/device-management-interface/go
	rm -rf vendor/github.com/opencord/device-management-interface/go/vendor
endif

local-lib-go: ## Copies a local version of the voltha-lib-go dependency into the vendor directory
ifdef LOCAL_LIB_GO
	mkdir -p vendor/github.com/opencord/voltha-lib-go/v4/pkg
	cp -r ${LOCAL_LIB_GO}/pkg/* vendor/github.com/opencord/voltha-lib-go/v4/pkg/
endif

## Docker targets
build: docker-build ## Alias for 'docker build'

# docker-build: local-protos local-lib-go ## Build opendevice manager docker image (set BUILD_PROFILED=true to also build the profiled image)
docker-build: ## Build opendevice manager docker image (set BUILD_PROFILED=true to also build the profiled image)
	docker build $(DOCKER_BUILD_ARGS) -t ${OPENDEVMGR_IMAGENAME} -f docker/Dockerfile.opendevmgr .
ifdef BUILD_PROFILED
	docker build $(DOCKER_BUILD_ARGS) --build-arg EXTRA_GO_BUILD_TAGS="-tags profile" -t ${OPENDEVMGR_IMAGENAME}-profile -f docker/Dockerfile.opendevmgr .
endif

docker-push: ## Push the docker images to an external repository
	docker push ${OPENDEVMGR_IMAGENAME}
ifdef BUILD_PROFILED
	docker push ${OPENDEVMGR_IMAGENAME}-profile
endif

docker-kind-load: ## Load docker images into a KinD cluster
	@if [ "`kind get clusters | grep opendevmgr-$(TYPE)`" = '' ]; then echo "no opendevmgr-$(TYPE) cluster found" && exit 1; fi
	kind load docker-image ${OPENDEVMGR_IMAGENAME} --name=opendevmgr-$(TYPE) --nodes $(shell kubectl get nodes --template='{{range .items}}{{.metadata.name}},{{end}}' | sed 's/,$$//')

## lint and unit tests
lint-dockerfile: ## Perform static analysis on Dockerfile
	@echo "Running Dockerfile lint check ..."
	@${HADOLINT} $$(find . -name "Dockerfile.*")
	@echo "Dockerfile lint check OK"

lint-mod: ## Verify the Go dependencies
	@echo "Running dependency check..."
	@${GO} mod verify
	@echo "Dependency check OK. Running vendor check..."
	@git status > /dev/null
	@git diff-index --quiet HEAD -- go.mod go.sum vendor || (echo "ERROR: Staged or modified files must be committed before running this test" && git status -- go.mod go.sum vendor && exit 1)
	@[[ `git ls-files --exclude-standard --others go.mod go.sum vendor` == "" ]] || (echo "ERROR: Untracked files must be cleaned up before running this test" && git status -- go.mod go.sum vendor && exit 1)
	${GO} mod tidy
	${GO} mod vendor
	@git status > /dev/null
	@git diff-index --quiet HEAD -- go.mod go.sum vendor || (echo "ERROR: Modified files detected after running go mod tidy / go mod vendor" && git status -- go.mod go.sum vendor && git checkout -- go.mod go.sum vendor && exit 1)
	@[[ `git ls-files --exclude-standard --others go.mod go.sum vendor` == "" ]] || (echo "ERROR: Untracked files detected after running go mod tidy / go mod vendor" && git status -- go.mod go.sum vendor && git checkout -- go.mod go.sum vendor && exit 1)
	@echo "Vendor check OK."

lint: local-lib-go lint-mod lint-dockerfile ## Run all lint targets

test: ## Run unit tests
	@echo "Call unit test case suite"

sca: ## Runs static code analysis with the golangci-lint tool
	@rm -rf ./sca-report
	@mkdir -p ./sca-report
	@echo "Running static code analysis..."
	@${GOLANGCI_LINT} run --deadline=6m --out-format junit-xml ./... | tee ./sca-report/sca-report.xml
	@echo ""
	@echo "Static code analysis OK"

clean: distclean ## Removes any local filesystem artifacts generated by a build

distclean: ## Removes any local filesystem artifacts generated by a build or test run
	rm -rf ./sca-report

mod-update: ## Update go mod files
	${GO} mod tidy
	${GO} mod vendor

# For each makefile target, add ## <description> on the target line and it will be listed by 'make help'
help: ## Print help for each Makefile target
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@grep '^[[:alpha:]_-]*:.* ##' $(MAKEFILE_LIST) \
		| sort | awk 'BEGIN {FS=":.* ## "}; {printf "%-25s : %s\n", $$1, $$2};'