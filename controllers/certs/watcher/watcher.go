package watcher

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	apiGroup = "admissionregistration.k8s.io"
)

func New(name string, certFile string) manager.Runnable {
	return &watcher{
		name:     name,
		certFile: certFile,
	}
}

type watcher struct {
	name      string
	certFile  string
	isV1      bool
	client    client.Client
	config    *rest.Config
	caWatcher *fsnotify.Watcher
	patch     func(ctx context.Context, caCert []byte) error
	logger    logr.Logger
}

func (w *watcher) Start(stopCh <-chan struct{}) error {
	var err error

	isV1, err := w.supportsV1()
	if err != nil {
		return err
	}
	if isV1 {
		w.patch = w.patchHooksV1
	} else {
		w.patch = w.patchHooksBeta1V1
	}

	if err = w.caChanged(nil); err != nil {
		return err
	}

	w.caWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := w.caWatcher.Add(w.certFile); err != nil {
		return err
	}

	w.logger.Info("Starting webhook ca certificate watcher")
	go w.WatchCA()

	// Block until the stop channel is closed.
	<-stopCh
	return w.caWatcher.Close()
}

func (w *watcher) NeedLeaderElection() bool {
	return false
}

func (w *watcher) InjectClient(c client.Client) error {
	w.client = c
	return nil
}

func (w *watcher) InjectConfig(c *rest.Config) error {
	w.config = c
	return nil
}

func (w *watcher) InjectLogger(l logr.Logger) error {
	w.logger = l
	return nil
}
