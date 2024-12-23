# Include toolbox tasks
include ./.toolbox.mk

# generate mocks
mocks: tb.mockgen
	$(TB_MOCKGEN) -destination pkg/mocks/logr/mock.go github.com/go-logr/logr LogSink
	$(TB_MOCKGEN) -destination pkg/mocks/client/mock.go sigs.k8s.io/controller-runtime/pkg/client Client

# Format code
fmt: tb.golines tb.gofumpt
	$(TB_GOLINES) --base-formatter="$(TB_GOFUMPT)" --max-len=120 --write-output .

# Run tests
test: mocks fmt lint tb.ginkgo
	$(TB_GINKGO) ./...

# Run go golanci-lint
lint: tb.golangci-lint
	$(TB_GOLANGCI_LINT) run --fix