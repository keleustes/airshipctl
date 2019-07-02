SHELL := /bin/bash

GO_FLAGS            := -ldflags '-extldflags "-static"' -tags=netgo
GOPATH              := ${GOPATH}

# Directories.
BINDIR              := bin
EXECUTABLE_CLI      := airshipctl

TOOLS_DIR           := tools
TOOLSBINDIR         := $(TOOLS_DIR)/bin
MANIFEST_ROOT       ?= config
CRD_ROOT            ?= $(MANIFEST_ROOT)/crd/bases

# Binaries.
CONTROLLER_GEN      := $(TOOLSBINDIR)/controller-gen
LINTER              := $(TOOLSBINDIR)/golangci-lint
MOCKGEN             := $(TOOLSBINDIR)/mockgen
CONVERSION_GEN      := $(TOOLSBINDIR)/conversion-gen
KUBEBUILDER         := $(TOOLSBINDIR)/kubebuilder

# linting
LINTER_CONFIG       := .golangci.yaml

# docker
DOCKER_MAKE_TARGET  := build

# docker image options
DOCKER_REGISTRY     ?= quay.io
DOCKER_IMAGE_NAME   ?= airshipctl
DOCKER_IMAGE_PREFIX ?= airshipit
DOCKER_IMAGE_TAG    ?= dev
DOCKER_IMAGE        ?= $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_PREFIX)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_TARGET_STAGE ?= release

# go options
PKG                 := ./...
TESTS               := .
TEST_FLAGS          :=
COVER_FLAGS         :=
COVER_PROFILE       := cover.out
COVER_PKG           := $(shell go list ./... | tail -n+2 | grep -v "opendev.org/airship/airshipctl/testutil" | paste -sd"," -)

.PHONY: get-modules
get-modules:
	@go mod download
	cd $(TOOLS_DIR); go mod download

.PHONY: modules
modules: ## Runs go mod to ensure modules are up to date.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

## --------------------------------------
## Generate
## --------------------------------------

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN) ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) crd \
		paths=$(GOPATH)/src/sigs.k8s.io/cluster-api/api/v1alpha2/... \
		crd:trivialVersions=true \
		output:crd:dir=$(CRD_ROOT) \
		output:stdout
	$(CONTROLLER_GEN) crd \
		paths=$(GOPATH)/src/sigs.k8s.io/cluster-api-provider-baremetal/api/v1alpha2/... \
		crd:trivialVersions=true \
		output:crd:dir=$(CRD_ROOT) \
		output:stdout
	$(CONTROLLER_GEN) crd \
		paths=$(GOPATH)/src/sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2/... \
		crd:trivialVersions=true \
		output:crd:dir=$(CRD_ROOT) \
		output:stdout

## --------------------------------------
## Binaries
## --------------------------------------

.PHONY: build
build: get-modules
	@CGO_ENABLED=0 go build -o $(BINDIR)/$(EXECUTABLE_CLI) $(GO_FLAGS)

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(CONTROLLER_GEN): $(TOOLS_DIR)/go.mod # Build controller-gen from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BINDIR)/controller-gen sigs.k8s.io/controller-tools/cmd/controller-gen

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod # Build golangci-lint from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BINDIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(MOCKGEN): $(TOOLS_DIR)/go.mod # Build mockgen from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BINDIR)/mockgen github.com/golang/mock/mockgen

$(CONVERSION_GEN): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -tags=tools -o $(BINDIR)/conversion-gen k8s.io/code-generator/cmd/conversion-gen

$(KUBEBUILDER): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); ./install_kubebuilder.sh

.PHONY: installtools
installtools: $(KUBEBUILDER) $(GOLANGCI_LINT)

## --------------------------------------
## Testing
## --------------------------------------

.PHONY: test
test: build
test: cover

.PHONY: unit-tests
unit-tests: TESTFLAGS += -race -v
unit-tests:
	@echo "Performing unit test step..."
	@go test -run $(TESTS) $(PKG) $(TESTFLAGS) $(COVER_FLAGS)
	@echo "All unit tests passed"

.PHONY: cover
cover: COVER_FLAGS = -covermode=atomic -coverprofile=$(COVER_PROFILE) -coverpkg=$(COVER_PKG)
cover: unit-tests
	@./tools/coverage_check $(COVER_PROFILE)

.PHONY: lint
lint: $(LINTER)
	@echo "Performing linting step..."
	@./$(LINTER) run --config $(LINTER_CONFIG)
	@echo "Linting completed successfully"

.PHONY: docker-image
docker-image:
	@docker build . --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) --tag $(DOCKER_IMAGE) --target $(DOCKER_TARGET_STAGE)

.PHONY: print-docker-image-tag
print-docker-image-tag:
	@echo "$(DOCKER_IMAGE)"

.PHONY: docker-image-unit-tests
docker-image-unit-tests: DOCKER_MAKE_TARGET = cover
docker-image-unit-tests: DOCKER_TARGET_STAGE = builder
docker-image-unit-tests: docker-image

.PHONY: docker-image-lint
docker-image-lint: DOCKER_MAKE_TARGET = lint
docker-image-lint: DOCKER_TARGET_STAGE = builder
docker-image-lint: docker-image

.PHONY: clean
clean:
	@rm -fr $(BINDIR)
	@rm -fr $(COVER_PROFILE)

.PHONY: docs
docs:
	@echo "TODO"

$(TOOLBINDIR):
	mkdir -p $(TOOLBINDIR)

$(LINTER): $(TOOLBINDIR)
	./tools/install_linter

.PHONY: update-golden
update-golden: delete-golden
update-golden: TESTFLAGS += -update
update-golden: PKG = opendev.org/airship/airshipctl/cmd/...
update-golden: unit-tests

# The delete-golden target is a utility for update-golden
.PHONY: delete-golden
delete-golden:
	@find cmd -type f -name "*.golden" -delete
