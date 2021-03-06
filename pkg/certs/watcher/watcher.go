package watcher

import (
	"context"
	"path/filepath"

	"github.com/bakito/operator-utils/pkg/certs"
	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// New create a new watcher
func New(opts certs.Options) manager.Runnable {
	w := &watcher{
		opts: opts.ApplyDefaults(opts.Name),
	}
	w.certFile = filepath.Join(opts.CertDir, opts.CACert)
	return w
}

type watcher struct {
	opts      certs.Options
	certFile  string
	client    client.Client
	config    *rest.Config
	caWatcher *fsnotify.Watcher
	patch     func(ctx context.Context, caCert []byte) error
	logger    logr.Logger
}

func (w *watcher) Start(ctx context.Context) error {
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

	if err = w.syncHooks(); err != nil {
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
	go w.watchCA()

	// Block until the stop channel is closed.
	<-ctx.Done()
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
