# find or install mockgen
ifeq (, $(shell which mockgen))
 $(shell go get github.com/golang/mock/mockgen)
endif
ifeq (, $(shell which ginkgo))
 $(shell go get github.com/onsi/ginkgo/ginkgo)
endif

# generate mocks
mocks:
	mockgen -destination pkg/mocks/logr/mock.go github.com/go-logr/logr Logger
	mockgen -destination pkg/mocks/client/mock.go sigs.k8s.io/controller-runtime/pkg/client Client

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run tests
test: mocks fmt vet
	go test ./...