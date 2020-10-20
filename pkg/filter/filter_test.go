package filter_test

import (
	"github.com/bakito/operator-utils/pkg/filter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("Filter", func() {
	var (
		np  filter.NamePredicate
		pod *corev1.Pod
	)

	BeforeEach(func() {
		np = filter.NamePredicate{
			Names: []string{"foo"},
		}
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{},
		}
	})

	Describe("Create", func() {
		var (
			e event.CreateEvent
		)
		BeforeEach(func() {
			e = event.CreateEvent{
				Meta: pod,
			}
		})
		Context("Create match", func() {
			BeforeEach(func() {
				pod.Name = "foo"
			})
			It("should match", func() {
				Ω(np.Create(e)).To(BeTrue())
			})
		})

		Context("Create not match", func() {
			BeforeEach(func() {
				pod.Name = "bar"
			})
			It("should match", func() {
				Ω(np.Create(e)).To(BeFalse())
			})
		})
	})

	Describe("Delete", func() {
		var (
			e event.DeleteEvent
		)
		BeforeEach(func() {
			e = event.DeleteEvent{
				Meta: pod,
			}
		})
		Context("Delete match", func() {
			BeforeEach(func() {
				pod.Name = "foo"
			})
			It("should match", func() {
				Ω(np.Delete(e)).To(BeTrue())
			})
		})

		Context("Delete not match", func() {
			BeforeEach(func() {
				pod.Name = "bar"
			})
			It("should match", func() {
				Ω(np.Delete(e)).To(BeFalse())
			})
		})
	})

	Describe("Update", func() {
		var (
			e event.UpdateEvent
		)
		BeforeEach(func() {
			e = event.UpdateEvent{
				MetaNew: pod,
			}
		})
		Context("Update match", func() {
			BeforeEach(func() {
				pod.Name = "foo"
			})
			It("should match", func() {
				Ω(np.Update(e)).To(BeTrue())
			})
		})

		Context("Update not match", func() {
			BeforeEach(func() {
				pod.Name = "bar"
			})
			It("should match", func() {
				Ω(np.Update(e)).To(BeFalse())
			})
		})
	})

	Describe("Generic", func() {
		var (
			e event.GenericEvent
		)
		BeforeEach(func() {
			e = event.GenericEvent{
				Meta: pod,
			}
		})
		Context("Generic match", func() {
			BeforeEach(func() {
				pod.Name = "foo"
			})
			It("should match", func() {
				Ω(np.Generic(e)).To(BeTrue())
			})
		})

		Context("Generic not match", func() {
			BeforeEach(func() {
				pod.Name = "bar"
			})
			It("should match", func() {
				Ω(np.Generic(e)).To(BeFalse())
			})
		})
	})
})
