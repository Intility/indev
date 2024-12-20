# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

BINARY_NAME=indev

DOT := $(shell command -v dot 2> /dev/null)

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: mocks
mocks: mockery ## Generate mocks
	$(MOCKERY)

.PHONY: docs
docs: vhs ## Generate docs
	@echo "Generating gifs..."
	@$(VHS) docs/create_cluster_interactive.tape -o docs/create_cluster_interactive.gif

.PHONY: rename
rename: gum ## Rename the project
	@set -e; \
	OLD_NAME=$(shell awk 'NR==1{print $$2}' go.mod); \
	NEW_NAME=$(shell $(GUM) input --prompt "github.com/intility/" --placeholder "new_name"); \
	./scripts/rename-project.sh $$OLD_NAME github.com/intility/$$NEW_NAME

.PHONY: build-graph
build-graph: actiongraph ## Generate build graph
	@echo "Generating build graph..."
	@go clean -cache
	@CGO_ENABLED=0 go build -o $(LOCALBIN)/$(BINARY_NAME) -debug-actiongraph=$(LOCALBIN)/build-graph.json ./cmd/indev/main.go
	@$(ACTIONGRAPH) -f $(LOCALBIN)/build-graph.json top -n 100

Q ?= "**"
.PHONY: dependency-graph
ifndef DOT
dependency-graph: dependency-graph-text ## Generate dependency analysis. Optional Q="<filter>". See gomod graph: https://github.com/Helcaraxan/gomod
else
dependency-graph: dependency-graph-svg
endif

.PHONY: dependency-graph-text
dependency-graph-text: gomod
	@echo "Generating dependency graph..." 3>&2 2>&1 1>&3
	@echo -e '\033[0;33m[NOTE]\033[0m: Graphviz not installed (https://www.graphviz.org/download/). Printing dot file to stdout' 3>&2 2>&1 1>&3
	@$(GOMOD) graph '$(Q)'

.PHONY: dependency-graph-svg
dependency-graph-svg: gomod
	@echo "Generating dependency graph..." 3>&2 2>&1 1>&3
	$(eval SVG := $(shell mktemp -u).svg)
	@$(GOMOD) graph '$(Q)' | dot -Tsvg -o $(SVG)
	open $(SVG)

##@ Build

.PHONY: build
build: fmt vet ## Build the code generator.
	CGO_ENABLED=0 go build -o $(LOCALBIN)/$(BINARY_NAME) ./cmd/indev/main.go

.PHONY: run
run: fmt vet check-env-vars ## Run the example app.
	go run ./cmd/indev/main.go

.PHONY: generate
generate: ## Generate code.
	go generate ./...

.PHONY: cross-compile
cross-compile: ## Cross compile the code.
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(LOCALBIN)/$(BINARY_NAME)-linux-arm64 ./main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(LOCALBIN)/$(BINARY_NAME)-linux-amd64 ./main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(LOCALBIN)/$(BINARY_NAME)-darwin-arm64 ./main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(LOCALBIN)/$(BINARY_NAME)-darwin-amd64 ./main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(LOCALBIN)/$(BINARY_NAME)-windows-amd64.exe ./main.go

##@ Pre Deployment

.PHONY: security-scan
security-scan: gosec ## Run security scan on the codebase.
	$(GOSEC) ./...

.PHONY: dependency-scan
dependency-scan: govulncheck ## Run dependency scan on the codebase.
	$(GOVULNCHECK) ./...

##@ Deployment


##@ Dependencies

.PHONY: check-env-vars
check-env-vars:
	@missing_vars=""; \
	for var in AOAI_API_KEY AOAI_ENDPOINT AOAI_API_VERSION AOAI_MODEL_DEPLOYMENT; do \
			if [ -z "$${!var}" ]; then \
					missing_vars="$$missing_vars $$var"; \
			fi \
	done; \
	if [ -n "$$missing_vars" ]; then \
			echo "Error: the following env vars are not set:$$missing_vars"; \
			exit 1; \
	fi

.PHONY: lang-gen
lang-gen:
	cd generator && go build -o $(LOCALBIN)/lang-gen lang-gen.go

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
TOOLKIT_TOOLS_GEN = $(LOCALBIN)/toolkit-tools-gen-$(TOOLKIT_TOOLS_GEN_VERSION)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
GOSEC = $(LOCALBIN)/gosec-$(GOSEC_VERSION)
GOVULNCHECK = $(LOCALBIN)/govulncheck-$(GOVULNCHECK_VERSION)
MOCKERY = $(LOCALBIN)/mockery-$(MOCKERY_VERSION)
VHS = $(LOCALBIN)/vhs-$(VHS_VERSION)
GUM = $(LOCALBIN)/gum-$(GUM_VERSION)
ACTIONGRAPH = $(LOCALBIN)/actiongraph-$(ACTIONGRAPH_VERSION)
GOMOD = $(LOCALBIN)/gomod-$(GOMOD_VERSION)

## Tool Versions
TOOLKIT_TOOLS_GEN_VERSION ?= latest
GOLANGCI_LINT_VERSION ?= v1.59.1
GOSEC_VERSION ?= latest
GOVULNCHECK_VERSION ?= latest
MOCKERY_VERSION ?= v2.42.1
VHS_VERSION ?= latest
GUM_VERSION ?= latest
ACTIONGRAPH_VERSION ?= latest
GOMOD_VERSION ?= latest

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: gosec
gosec: $(GOSEC) ## Download gosec locally if necessary.
$(GOSEC): $(LOCALBIN)
	$(call go-install-tool,$(GOSEC),github.com/securego/gosec/v2/cmd/gosec,$(GOSEC_VERSION))

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK) ## Download govulncheck locally if necessary.
$(GOVULNCHECK): $(LOCALBIN)
	$(call go-install-tool,$(GOVULNCHECK),golang.org/x/vuln/cmd/govulncheck,$(GOVULNCHECK_VERSION))

.PHONY: toolkit-tools-gen
toolkit-tools-gen: $(TOOLKIT_TOOLS_GEN) ## Download toolkit-tools-gen locally if necessary.
$(TOOLKIT_TOOLS_GEN): $(LOCALBIN)
	$(call go-install-tool,$(TOOLKIT_TOOLS_GEN),github.com/intility/go-openai-toolkit/cmd/toolkit-tools-gen,$(TOOLKIT_TOOLS_GEN_VERSION))

.PHONY: mockery
mockery: $(MOCKERY) ## Download mockery locally if necessary.
$(MOCKERY): $(LOCALBIN)
	$(call go-install-tool,$(MOCKERY),github.com/vektra/mockery/v2,$(MOCKERY_VERSION))

.PHONY: vhs
vhs: $(VHS) ## Download vhs locally if necessary.
$(VHS): $(LOCALBIN)
	# check for the precense of ttyd and ffmpeg
	@which ttyd || { echo "ttyd is not installed. Please install it from https://github.com/tsl0922/ttyd"; exit 1; }
	@which ffmpeg || { echo "ffmpeg is not installed. Please install it from https://ffmpeg.org"; exit 1; }
	$(call go-install-tool,$(VHS),github.com/charmbracelet/vhs,$(VHS_VERSION))

.PHONY: gum
gum: $(GUM) ## Download gum locally if necessary.
$(GUM): $(LOCALBIN)
	$(call go-install-tool,$(GUM),github.com/charmbracelet/gum,$(GUM_VERSION))

.PHONY: actiongraph
actiongraph: $(ACTIONGRAPH) ## Download actiongraph locally if necessary.
$(ACTIONGRAPH): $(LOCALBIN)
	$(call go-install-tool,$(ACTIONGRAPH),github.com/icio/actiongraph,$(ACTIONGRAPH_VERSION))

.PHONY: gomod
gomod: $(GOMOD) ## Download godepgraph locally if necessary.
$(GOMOD): $(LOCALBIN)
	$(call go-install-tool,$(GOMOD),github.com/Helcaraxan/gomod,$(GOMOD_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
