package filter

import "sigs.k8s.io/controller-runtime/pkg/event"

// NamePredicate only watches objects with given name
type NamePredicate struct {
	Namespace string
	Names     []string
}

// Create implements Predicate
func (p NamePredicate) Create(e event.CreateEvent) bool {
	if p.Namespace == "" || e.Meta.GetNamespace() == p.Namespace {
		for _, n := range p.Names {
			if e.Meta.GetName() == n {
				return true
			}
		}
	}
	return false
}

// Delete implements Predicate
func (p NamePredicate) Delete(e event.DeleteEvent) bool {
	if p.Namespace == "" || e.Meta.GetNamespace() == p.Namespace {
		for _, n := range p.Names {
			if e.Meta.GetName() == n {
				return true
			}
		}
	}
	return false
}

// Update implements Predicate
func (p NamePredicate) Update(e event.UpdateEvent) bool {
	if p.Namespace == "" || e.MetaNew.GetNamespace() == p.Namespace {
		for _, n := range p.Names {
			if e.MetaNew.GetName() == n {
				return true
			}
		}
	}
	return false
}

// Generic implements Predicate
func (p NamePredicate) Generic(e event.GenericEvent) bool {
	if p.Namespace == "" || e.Meta.GetNamespace() == p.Namespace {
		for _, n := range p.Names {
			if e.Meta.GetName() == n {
				return true
			}
		}
	}
	return false
}
