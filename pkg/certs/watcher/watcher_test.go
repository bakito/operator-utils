package watcher

import (
	"context"
	"encoding/base64"

	"github.com/bakito/operator-utils/pkg/certs"
	mock_client "github.com/bakito/operator-utils/pkg/mocks/client"
	mock_logr "github.com/bakito/operator-utils/pkg/mocks/logr"
	gm "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	arv1 "k8s.io/api/admissionregistration/v1"
	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Watcher", func() {
	var (
		ctx          context.Context
		w            *watcher
		mockCtrl     *gm.Controller //gomock struct
		mockLog      *mock_logr.MockLogger
		mockClient   *mock_client.MockClient
		cert         []byte
		oldCert      []byte
		expectedCert string
	)

	BeforeEach(func() {
		ctx = context.TODO()
		mockCtrl = gm.NewController(GinkgoT())
		mockLog = mock_logr.NewMockLogger(mockCtrl)
		mockClient = mock_client.NewMockClient(mockCtrl)
		w = New(certs.Options{}).(*watcher)
		w.InjectClient(mockClient)
		w.InjectLogger(mockLog)
		w.InjectConfig(&rest.Config{})

		cert = []byte("cert")
		oldCert = []byte("old-cert")
		expectedCert = base64.StdEncoding.EncodeToString(cert)

	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("NeedLeaderElection", func() {
		It("should be false", func() {
			Ω(w.NeedLeaderElection()).To(BeFalse())
		})
	})

	Context("supportsV1", func() {
		var (
			apiGroupList *metav1.APIGroupList
		)
		BeforeEach(func() {
			apiGroupList = &metav1.APIGroupList{
				Groups: []metav1.APIGroup{},
			}
		})
		Context("supportsV1 v1", func() {
			BeforeEach(func() {
				apiGroupList.Groups = append(apiGroupList.Groups, metav1.APIGroup{
					Name:     arv1.GroupName,
					Versions: []metav1.GroupVersionForDiscovery{{Version: "v1"}},
				})
			})
			It("should be true", func() {
				supports, err := w.supportsV1Internal(apiGroupList)
				Ω(supports).To(BeTrue())
				Ω(err).To(BeNil())
			})
		})
		Context("supportsV1 v1beta1", func() {
			BeforeEach(func() {
				apiGroupList.Groups = append(apiGroupList.Groups, metav1.APIGroup{
					Name:     arv1.GroupName,
					Versions: []metav1.GroupVersionForDiscovery{{Version: "v1beta1"}},
				})
			})
			It("should be false", func() {
				supports, err := w.supportsV1Internal(apiGroupList)
				Ω(supports).To(BeFalse())
				Ω(err).To(BeNil())
			})
		})
		Context("supportsV1 error", func() {
			BeforeEach(func() {
				apiGroupList.Groups = append(apiGroupList.Groups, metav1.APIGroup{
					Name:     "foo",
					Versions: []metav1.GroupVersionForDiscovery{{Version: "foo"}},
				})
			})
			It("should be error", func() {
				supports, err := w.supportsV1Internal(apiGroupList)
				Ω(supports).To(BeFalse())
				Ω(err).To(HaveOccurred())
			})
		})
	})

	Context("patchHooksV1", func() {
		BeforeEach(func() {
			mockClient.EXPECT().Get(ctx, gm.Any(), gm.AssignableToTypeOf(&arv1.MutatingWebhookConfiguration{})).
				Do(func(ctx context.Context, key client.ObjectKey, obj *arv1.MutatingWebhookConfiguration) error {
					obj.Webhooks = []arv1.MutatingWebhook{{Name: "mwc", ClientConfig: arv1.WebhookClientConfig{CABundle: oldCert}}}
					return nil
				})
			mockClient.EXPECT().Get(ctx, gm.Any(), gm.AssignableToTypeOf(&arv1.ValidatingWebhookConfiguration{})).
				Do(func(ctx context.Context, key client.ObjectKey, obj *arv1.ValidatingWebhookConfiguration) error {
					obj.Webhooks = []arv1.ValidatingWebhook{{Name: "vwc", ClientConfig: arv1.WebhookClientConfig{CABundle: oldCert}}}
					return nil
				})
		})
		Context("NeedLeaderElection with changes", func() {
			BeforeEach(func() {
				mockClient.EXPECT().Patch(ctx, gm.AssignableToTypeOf(&arv1.MutatingWebhookConfiguration{}), gm.Any(), gm.Any()).
					Do(func(ctx context.Context, obj client.Object, patch client.Patch) error {
						Ω(patch.Type()).To(BeEquivalentTo(types.StrategicMergePatchType))
						b, err := patch.Data(obj)
						Ω(err).To(BeNil())
						Ω(string(b)).To(ContainSubstring(expectedCert))
						return nil
					})
				mockClient.EXPECT().Patch(ctx, gm.AssignableToTypeOf(&arv1.ValidatingWebhookConfiguration{}), gm.Any(), gm.Any()).
					Do(func(ctx context.Context, obj client.Object, patch client.Patch) error {
						Ω(patch.Type()).To(BeEquivalentTo(types.StrategicMergePatchType))
						b, err := patch.Data(obj)
						Ω(err).To(BeNil())
						Ω(string(b)).To(ContainSubstring(expectedCert))
						return nil
					})
				mockLog.EXPECT().WithValues(gm.Any(), gm.Any(), gm.Any(), gm.Any()).Return(mockLog).Times(2)
				mockLog.EXPECT().WithValues(gm.Any(), gm.Any()).Return(mockLog).Times(2)
				mockLog.EXPECT().Info(gm.Any()).Times(2)
			})

			It("should not fail", func() {
				Ω(w.patchHooksV1(ctx, cert)).To(BeNil())
			})
		})
		Context("NeedLeaderElection no changes", func() {
			It("should not fail", func() {
				Ω(w.patchHooksV1(ctx, oldCert)).To(BeNil())
			})
		})
	})

	Context("patchHooksBeta1V1", func() {
		BeforeEach(func() {
			mockClient.EXPECT().Get(ctx, gm.Any(), gm.AssignableToTypeOf(&arv1beta1.MutatingWebhookConfiguration{})).
				Do(func(ctx context.Context, key client.ObjectKey, obj *arv1beta1.MutatingWebhookConfiguration) error {
					obj.Webhooks = []arv1beta1.MutatingWebhook{{Name: "mwc", ClientConfig: arv1beta1.WebhookClientConfig{CABundle: oldCert}}}
					return nil
				})
			mockClient.EXPECT().Get(ctx, gm.Any(), gm.AssignableToTypeOf(&arv1beta1.ValidatingWebhookConfiguration{})).
				Do(func(ctx context.Context, key client.ObjectKey, obj *arv1beta1.ValidatingWebhookConfiguration) error {
					obj.Webhooks = []arv1beta1.ValidatingWebhook{{Name: "vwc", ClientConfig: arv1beta1.WebhookClientConfig{CABundle: oldCert}}}
					return nil
				})
		})
		Context("NeedLeaderElection with changes", func() {
			BeforeEach(func() {
				mockClient.EXPECT().Patch(ctx, gm.AssignableToTypeOf(&arv1beta1.MutatingWebhookConfiguration{}), gm.Any(), gm.Any()).
					Do(func(ctx context.Context, obj client.Object, patch client.Patch) error {
						Ω(patch.Type()).To(BeEquivalentTo(types.StrategicMergePatchType))
						b, err := patch.Data(obj)
						Ω(err).To(BeNil())
						Ω(string(b)).To(ContainSubstring(expectedCert))
						return nil
					})
				mockClient.EXPECT().Patch(ctx, gm.AssignableToTypeOf(&arv1beta1.ValidatingWebhookConfiguration{}), gm.Any(), gm.Any()).
					Do(func(ctx context.Context, obj client.Object, patch client.Patch) error {
						Ω(patch.Type()).To(BeEquivalentTo(types.StrategicMergePatchType))
						b, err := patch.Data(obj)
						Ω(err).To(BeNil())
						Ω(string(b)).To(ContainSubstring(expectedCert))
						return nil
					})
				mockLog.EXPECT().WithValues(gm.Any(), gm.Any(), gm.Any(), gm.Any()).Return(mockLog).Times(2)
				mockLog.EXPECT().WithValues(gm.Any(), gm.Any()).Return(mockLog).Times(2)
				mockLog.EXPECT().Info(gm.Any()).Times(2)
			})

			It("should not fail", func() {
				Ω(w.patchHooksBeta1V1(ctx, cert)).To(BeNil())
			})
		})
		Context("NeedLeaderElection no changes", func() {
			It("should not fail", func() {
				Ω(w.patchHooksBeta1V1(ctx, oldCert)).To(BeNil())
			})
		})
	})
})
