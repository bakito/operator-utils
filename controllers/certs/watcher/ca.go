package watcher

import (
	"context"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

// WatchCA reads events from the watcher's channel and reacts to changes.
func (w *watcher) WatchCA() {
	for {
		select {
		case event, ok := <-w.caWatcher.Events:
			// Channel is closed.
			if !ok {
				return
			}
			if isWrite(event) || isCreate(event) {
				_ = w.caChanged()
			}

		case err, ok := <-w.caWatcher.Errors:
			// Channel is closed.
			if !ok {
				return
			}

			w.logger.Error(err, "ca certificate watch error")
		}
	}
}

func (w *watcher) caChanged() error {
	ctx := context.TODO()
	dat, err := ioutil.ReadFile(w.certFile)
	if err != nil {
		w.logger.Error(err, "Error reading ca cert")
	}

	if err = w.patch(ctx, dat); err != nil {
		w.logger.Error(err, "Error patching webhook ca cert")
	}
	return err
}

func isWrite(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write
}

func isCreate(event fsnotify.Event) bool {
	return event.Op&fsnotify.Create == fsnotify.Create
}
