## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	mkdir -p $(TB_LOCALBIN)

## Tool Binaries
TB_GINKGO ?= $(TB_LOCALBIN)/ginkgo
TB_GOFUMPT ?= $(TB_LOCALBIN)/gofumpt
TB_GOLANGCI_LINT ?= $(TB_LOCALBIN)/golangci-lint
TB_GOLINES ?= $(TB_LOCALBIN)/golines
TB_MOCKGEN ?= $(TB_LOCALBIN)/mockgen

## Tool Versions
TB_GINKGO_VERSION ?= v2.22.1
TB_GOFUMPT_VERSION ?= v0.7.0
TB_GOLANGCI_LINT_VERSION ?= v1.62.2
TB_GOLINES_VERSION ?= v0.12.2
TB_MOCKGEN_VERSION ?= v0.5.0

## Tool Installer
.PHONY: tb.ginkgo
tb.ginkgo: $(TB_GINKGO) ## Download ginkgo locally if necessary.
$(TB_GINKGO): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/ginkgo || GOBIN=$(TB_LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@$(TB_GINKGO_VERSION)
.PHONY: tb.gofumpt
tb.gofumpt: $(TB_GOFUMPT) ## Download gofumpt locally if necessary.
$(TB_GOFUMPT): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/gofumpt || GOBIN=$(TB_LOCALBIN) go install mvdan.cc/gofumpt@$(TB_GOFUMPT_VERSION)
.PHONY: tb.golangci-lint
tb.golangci-lint: $(TB_GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(TB_GOLANGCI_LINT): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/golangci-lint || GOBIN=$(TB_LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(TB_GOLANGCI_LINT_VERSION)
.PHONY: tb.golines
tb.golines: $(TB_GOLINES) ## Download golines locally if necessary.
$(TB_GOLINES): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/golines || GOBIN=$(TB_LOCALBIN) go install github.com/segmentio/golines@$(TB_GOLINES_VERSION)
.PHONY: tb.mockgen
tb.mockgen: $(TB_MOCKGEN) ## Download mockgen locally if necessary.
$(TB_MOCKGEN): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/mockgen || GOBIN=$(TB_LOCALBIN) go install go.uber.org/mock/mockgen@$(TB_MOCKGEN_VERSION)

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_LOCALBIN)/ginkgo \
		$(TB_LOCALBIN)/gofumpt \
		$(TB_LOCALBIN)/golangci-lint \
		$(TB_LOCALBIN)/golines \
		$(TB_LOCALBIN)/mockgen

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile -f $(TB_LOCALDIR)/Makefile \
		github.com/onsi/ginkgo/v2/ginkgo \
		mvdan.cc/gofumpt@github.com/mvdan/gofumpt \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/segmentio/golines \
		go.uber.org/mock/mockgen@github.com/uber-go/mock
## toolbox - end