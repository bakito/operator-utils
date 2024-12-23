package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// New create a new pprof runnable
func New(addr string) manager.Runnable {
	return &pprofRunner{
		addr: addr,
	}
}

type pprofRunner struct {
	addr string
}

func (ppr *pprofRunner) Start(ctx context.Context) error {
	log := ctrl.Log.WithName("pprof").WithValues("addr", ppr.addr)
	log.Info("metrics server is starting to pprof\"")

	r := http.NewServeMux()
	// Register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:              ppr.addr,
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		log.Error(srv.ListenAndServe(), "error running pprof service")
	}()
	<-ctx.Done()
	log.Info("stopping pprof service")
	_ = srv.Shutdown(context.TODO())
	return nil
}
