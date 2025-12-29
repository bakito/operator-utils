## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi

# Helper functions
STRIP_V = $(patsubst v%,%,$(1))

## Tool Binaries
TB_GINKGO ?= $(TB_LOCALBIN)/ginkgo
TB_GOFUMPT ?= $(TB_LOCALBIN)/gofumpt
TB_GOLANGCI_LINT ?= $(TB_LOCALBIN)/golangci-lint
TB_GOLINES ?= $(TB_LOCALBIN)/golines
TB_MOCKGEN ?= $(TB_LOCALBIN)/mockgen

## Tool Versions
TB_GINKGO_VERSION ?= v2.27.3
TB_GOFUMPT_VERSION ?= v0.9.2
TB_GOLANGCI_LINT_VERSION ?= v2.7.2
TB_GOLANGCI_LINT_VERSION_NUM ?= $(call STRIP_V,$(TB_GOLANGCI_LINT_VERSION))
TB_GOLINES_VERSION ?= v0.13.0
TB_MOCKGEN_VERSION ?= v0.6.0

## Tool Installer
.PHONY: tb.ginkgo
tb.ginkgo: ## Download ginkgo locally if necessary.
	@test -s $(TB_GINKGO) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@$(TB_GINKGO_VERSION)
.PHONY: tb.gofumpt
tb.gofumpt: ## Download gofumpt locally if necessary.
	@test -s $(TB_GOFUMPT) || \
		GOBIN=$(TB_LOCALBIN) go install mvdan.cc/gofumpt@$(TB_GOFUMPT_VERSION)
.PHONY: tb.golangci-lint
tb.golangci-lint: ## Download golangci-lint locally if necessary.
	@test -s $(TB_GOLANGCI_LINT) && $(TB_GOLANGCI_LINT) --version | grep -q $(TB_GOLANGCI_LINT_VERSION_NUM) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(TB_GOLANGCI_LINT_VERSION)
.PHONY: tb.golines
tb.golines: ## Download golines locally if necessary.
	@test -s $(TB_GOLINES) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/segmentio/golines@$(TB_GOLINES_VERSION)
.PHONY: tb.mockgen
tb.mockgen: ## Download mockgen locally if necessary.
	@test -s $(TB_MOCKGEN) || \
		GOBIN=$(TB_LOCALBIN) go install go.uber.org/mock/mockgen@$(TB_MOCKGEN_VERSION)

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_GINKGO) \
		$(TB_GOFUMPT) \
		$(TB_GOLANGCI_LINT) \
		$(TB_GOLINES) \
		$(TB_MOCKGEN)

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile -f $(TB_LOCALDIR)/Makefile \
		github.com/onsi/ginkgo/v2/ginkgo \
		mvdan.cc/gofumpt@github.com/mvdan/gofumpt \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint?--version \
		github.com/segmentio/golines \
		go.uber.org/mock/mockgen@github.com/uber-go/mock
## toolbox - end
