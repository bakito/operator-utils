package log_test

import (
	"github.com/bakito/operator-utils/pkg/log"
	mock_logr "github.com/bakito/operator-utils/pkg/mocks/logr"
	"github.com/go-logr/logr"
	gm "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Log", func() {
	var (
		object   *corev1.Pod
		mockCtrl *gm.Controller //gomock struct
		mockSink *mock_logr.MockLogSink
		l        logr.Logger
	)

	BeforeEach(func() {
		object = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-name",
			},
		}
		mockCtrl = gm.NewController(GinkgoT())
		mockSink = mock_logr.NewMockLogSink(mockCtrl)
		l = logr.New(mockSink)
		mockSink.EXPECT().Init(gm.Any())
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("All fields correctly set if no kind available", func() {
		BeforeEach(func() {
			mockSink.EXPECT().WithValues("namespace", "test-namespace", "name", "test-name").Return(mockSink)
			mockSink.EXPECT().WithValues("kind", "Pod").Return(mockSink)
		})
		It("should get a logger", func() {
			Ω(log.With(l, object)).To(BeEquivalentTo(mockSink))
		})
	})

	Context("All fields correctly set if no kind", func() {
		BeforeEach(func() {
			object.TypeMeta = metav1.TypeMeta{
				Kind: "PodKind",
			}
			mockSink.EXPECT().WithValues("namespace", "test-namespace", "name", "test-name").Return(mockSink)
			mockSink.EXPECT().WithValues("kind", "PodKind").Return(mockSink)
		})
		It("should get a logger", func() {
			Ω(log.With(l, object)).To(BeEquivalentTo(mockSink))
		})
	})
})
