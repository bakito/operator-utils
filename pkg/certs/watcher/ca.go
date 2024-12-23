package watcher

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (w *watcher) watchCA() {
	for {
		select {
		case event, ok := <-w.caWatcher.Events:
			// Channel is closed.
			if !ok {
				return
			}

			_ = w.handleEvent(&event)

		case err, ok := <-w.caWatcher.Errors:
			// Channel is closed.
			if !ok {
				return
			}

			w.logger.Error(err, "webhook ca certificate watch error")
		}
	}
}

func (w *watcher) handleEvent(event *fsnotify.Event) error {
	// Only care about events which may modify the contents of the file.
	if !(isWrite(event) || isRemove(event) || isCreate(event)) {
		return nil
	}

	w.logger.V(1).Info("webhook ca certificate event", "event", event)

	// If the file was removed, re-add the watch.
	if isRemove(event) {
		if err := w.caWatcher.Add(w.certFile); err != nil {
			w.logger.Error(err, "error re-watching file")
			return err
		}
	}

	return w.syncHooks()
}

func (w *watcher) syncHooks() error {
	dat, err := os.ReadFile(w.certFile)
	if err != nil {
		w.logger.Error(err, "Error reading webhook ca cert")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err = w.patch(ctx, dat); err != nil {
		w.logger.Error(err, "Error patching webhook ca cert")
	}

	return err
}

func isWrite(event *fsnotify.Event) bool {
	return event != nil && event.Op&fsnotify.Write == fsnotify.Write
}

func isCreate(event *fsnotify.Event) bool {
	return event != nil && event.Op&fsnotify.Create == fsnotify.Create
}

func isRemove(event *fsnotify.Event) bool {
	return event != nil && event.Op&fsnotify.Remove == fsnotify.Remove
}
