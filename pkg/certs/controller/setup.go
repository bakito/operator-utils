package controller

import (
	"errors"

	"github.com/bakito/operator-utils/pkg/certs"
	"github.com/bakito/operator-utils/pkg/certs/watcher"
	"github.com/bakito/operator-utils/pkg/filter"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func New(log logr.Logger, namespace string, secretName string, opts certs.Options) Reconciler {
	return &reconciler{
		log:  log,
		opts: opts.ApplyDefaults(secretName),
		namespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      secretName,
		},
	}
}

func (r *reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.opts.Name == "" {
		return errors.New("no name defined")
	}
	r.Client = mgr.GetClient()

	// setup ca cert watcher
	w := watcher.New(r.opts)
	if err := mgr.Add(w); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(filter.NamePredicate{Names: []string{r.namespacedName.Name}}).
		Complete(r)
}
