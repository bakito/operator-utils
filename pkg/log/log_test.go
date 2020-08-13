package log_test

import (
	"github.com/bakito/operator-utils/pkg/log"
	mock_logr "github.com/bakito/operator-utils/pkg/mocks/logr"
	gm "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Log", func() {
	var (
		object   *corev1.Pod
		mockCtrl *gm.Controller //gomock struct
		mockLog  *mock_logr.MockLogger
	)

	BeforeEach(func() {
		object = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-name",
			},
		}
		mockCtrl = gm.NewController(GinkgoT())
		mockLog = mock_logr.NewMockLogger(mockCtrl)
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("All fields correctly set if no kind available", func() {
		BeforeEach(func() {
			mockLog.EXPECT().WithValues("namespace", "test-namespace", "name", "test-name").Return(mockLog)
			mockLog.EXPECT().WithValues("kind", "Pod").Return(mockLog)
		})
		It("should get a logger", func() {
			Expect(log.With(mockLog, object)).To(BeEquivalentTo(mockLog))
		})
	})

	Context("All fields correctly set if no kind", func() {
		BeforeEach(func() {
			object.TypeMeta = metav1.TypeMeta{
				Kind: "PodKind",
			}
			mockLog.EXPECT().WithValues("namespace", "test-namespace", "name", "test-name").Return(mockLog)
			mockLog.EXPECT().WithValues("kind", "PodKind").Return(mockLog)
		})
		It("should get a logger", func() {
			Expect(log.With(mockLog, object)).To(BeEquivalentTo(mockLog))
		})
	})
})
