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

// New create a new reconciler
func New(log logr.Logger, namespace string, secretName string, opts certs.Options) Reconciler {
	return &reconciler{
		log:  log,
		opts: opts.ApplyDefaults(secretName),
		nn: types.NamespacedName{
			Namespace: namespace,
			Name:      secretName,
		},
	}
}

func (r *reconciler) SetupWithManager(globalMgr, namespacedMgr ctrl.Manager) error {
	if r.opts.Name == "" {
		return errors.New("no name defined")
	}

	// setup ca cert watcher
	w := watcher.New(r.opts)
	if err := globalMgr.Add(w); err != nil {
		return err
	}

	r.Client = namespacedMgr.GetClient()

	return ctrl.NewControllerManagedBy(namespacedMgr).
		For(&corev1.Secret{}).
		WithEventFilter(filter.NamePredicate{
			Namespace: r.nn.Namespace,
			Names:     []string{r.nn.Name},
		}).
		Complete(r)
}
