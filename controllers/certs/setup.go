package certs

import (
	"path/filepath"

	"github.com/bakito/operator-utils/controllers/certs/watcher"
	"github.com/bakito/operator-utils/filter"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	// setup ca cert watcher
	w := watcher.New(r.NamespacedName.Name, filepath.Join(Dir, CACert))
	if err := mgr.Add(w); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(filter.NamePredicate{Names: []string{r.NamespacedName.Name}}).
		Complete(r)
}
