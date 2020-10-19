module github.com/bakito/operator-utils

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9
	// fix untli 0.2.1 is released https://github.com/go-logr/logr/issues/22
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/golang/mock v1.4.4
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.3
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.6.3
)
