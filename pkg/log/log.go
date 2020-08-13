package log

import (
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// With get a logger for a given runtime object. namespace, name and kind are added as value.
func With(base logr.Logger, object runtime.Object) logr.Logger {

	l := base

	if meta, ok := object.(metav1.Object); ok {
		l = l.WithValues("namespace", meta.GetNamespace(), "name", meta.GetName())
	}
	kind := object.GetObjectKind().GroupVersionKind().Kind
	if kind == "" {
		split := strings.Split(reflect.TypeOf(object).String(), ".")
		kind = split[len(split)-1]
	}
	return l.WithValues("kind", kind)
}
