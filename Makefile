# find or install mockgen
ifeq (, $(shell which mockgen))
 $(shell go install github.com/golang/mock/mockgen@v1.6.0)
endif
ifeq (, $(shell which ginkgo))
 $(shell go install github.com/onsi/ginkgo/ginkgo@latest)
endif

# generate mocks
mocks:
	mockgen -destination pkg/mocks/logr/mock.go github.com/go-logr/logr LogSink
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