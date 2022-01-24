package controller

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-logr/logr"
	"time"

	"github.com/bakito/operator-utils/pkg/certs"
	mock_client "github.com/bakito/operator-utils/pkg/mocks/client"
	mock_logr "github.com/bakito/operator-utils/pkg/mocks/logr"
	gm "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Controller", func() {
	var (
		r        *reconciler
		mockCtrl *gm.Controller //gomock struct
		mockSink *mock_logr.MockLogSink
		ctx      context.Context
	)
	BeforeEach(func() {
		mockCtrl = gm.NewController(GinkgoT())
		mockSink = mock_logr.NewMockLogSink(mockCtrl)
		log := logr.New(mockSink)
		r = New(log, "", "", certs.Options{}).(*reconciler)
		ctx = context.Background()

		mockSink.EXPECT().Init(gm.Any())
		mockSink.EXPECT().Enabled(gm.Any()).AnyTimes().Return(true)
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})
	Context("Reconcile", func() {
		var (
			mockClient *mock_client.MockClient
			req        ctrl.Request
		)
		BeforeEach(func() {
			mockClient = mock_client.NewMockClient(mockCtrl)
			req = ctrl.Request{}

			r.Client = mockClient
			r.nn = types.NamespacedName{}
		})
		Context("Recreate certs", func() {
			BeforeEach(func() {
				mockClient.EXPECT().Patch(gm.Any(), gm.AssignableToTypeOf(&corev1.Secret{}), gm.Any())
				mockSink.EXPECT().WithValues(gm.Any(), "certs", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "cert", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "namespace", gm.Any(), "name", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "kind", "Secret").Return(mockSink)
				mockSink.EXPECT().Info(gm.Any(), gm.Any()).Times(2)
			})
			It("All certs missing", func() {
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{}))
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
			It("key missing", func() {
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerCert: {2},
							certs.CACert:     {3},
						}
						return nil
					})
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
			It("cert missing", func() {
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerKey: {1},
							certs.CACert:    {3},
						}
						return nil
					})
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
			It("ca missing", func() {
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerKey:  {1},
							certs.ServerCert: {2},
						}
						return nil
					})
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})

		})

		Context("invalid cert", func() {
			BeforeEach(func() {
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerKey:  {1},
							certs.ServerCert: {2},
							certs.CACert:     {3},
						}
						return nil
					})
				mockClient.EXPECT().Patch(gm.Any(), gm.AssignableToTypeOf(&corev1.Secret{}), gm.Any())

				mockSink.EXPECT().WithValues(gm.Any(), "certs", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "namespace", gm.Any(), "name", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "kind", "Secret").Return(mockSink)
				mockSink.EXPECT().Error(gm.Any(), gm.Any(), gm.Any())
				mockSink.EXPECT().Info(gm.Any(), gm.Any())
			})
			It("invalid cert", func() {
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
		})
		Context("Certs expired", func() {
			BeforeEach(func() {
				k, c, ca, _ := r.createCerts(time.Now())
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerKey:  k,
							certs.ServerCert: c,
							certs.CACert:     ca,
						}
						return nil
					})
				mockClient.EXPECT().Patch(gm.Any(), gm.AssignableToTypeOf(&corev1.Secret{}), gm.Any())
				mockSink.EXPECT().WithValues(gm.Any(), "certs", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "namespace", gm.Any(), "name", gm.Any()).Return(mockSink)
				mockSink.EXPECT().WithValues(gm.Any(), "kind", "Secret").Return(mockSink)
				mockSink.EXPECT().Info(gm.Any(), gm.Any())
			})
			It("Certs expired", func() {
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
		})

		Context("Certs available", func() {
			BeforeEach(func() {
				k, c, ca, _ := r.createCerts(time.Now().AddDate(1, 0, 0))
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).
					Do(func(tx context.Context, key interface{}, sec *corev1.Secret) error {
						sec.Data = map[string][]byte{
							certs.ServerKey:  k,
							certs.ServerCert: c,
							certs.CACert:     ca,
						}
						return nil
					})
				mockSink.EXPECT().WithValues("certs", gm.Any()).Return(mockSink)
			})
			It("Recreate available", func() {
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
		})

		Context("Deleted", func() {
			BeforeEach(func() {
				err := errors.NewNotFound(schema.GroupResource{}, "")
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).Return(err)
				mockSink.EXPECT().WithValues(gm.Any(), "certs", gm.Any()).Return(mockSink)
				mockSink.EXPECT().Error(gm.Any(), err, gm.Any())
			})
			It("Deleted", func() {
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(BeNil())
			})
		})

		Context("Get with error", func() {
			BeforeEach(func() {
				err := fmt.Errorf("Get with error")
				mockClient.EXPECT().Get(gm.Any(), r.nn, gm.AssignableToTypeOf(&corev1.Secret{})).Return(err)
				mockSink.EXPECT().WithValues("certs", gm.Any()).Return(mockSink)
			})
			It("Get with error", func() {
				res, err := r.Reconcile(ctx, req)
				Ω(res.Requeue).To(BeFalse())
				Ω(err).To(HaveOccurred())
			})
		})
	})
	Context("createCerts", func() {
		var (
			serverKey  []byte
			serverCert []byte
			caCert     []byte
		)

		BeforeEach(func() {
			r.opts.Organization = "test.org"
			var err error
			serverKey, serverCert, caCert, err = r.createCerts(time.Now().AddDate(1, 0, 1))
			Ω(err).To(BeNil())
			Ω(serverKey).To(Not(BeNil()))
			Ω(serverCert).To(Not(BeNil()))
			Ω(caCert).To(Not(BeNil()))
		})
		It("key / cert", func() {

			cert := readCert(serverCert)

			Ω(cert.IsCA).To(BeFalse())
			Ω(cert.NotAfter.After(time.Now().AddDate(1, 0, 0))).To(BeTrue())
			Ω(cert.Subject.Organization).To(HaveLen(1))
			Ω(cert.Subject.Organization[0]).To(Equal("test.org"))
			Ω(cert.PublicKey).To(BeAssignableToTypeOf(&rsa.PublicKey{}))

			key := readKey(serverKey)

			Ω(key.PublicKey).To(Equal(*cert.PublicKey.(*rsa.PublicKey)))

		})
		It("ca", func() {
			ca := readCert(caCert)

			Ω(ca.IsCA).To(BeTrue())
			Ω(ca.NotAfter.After(time.Now().AddDate(1, 0, 0))).To(BeTrue())
			Ω(ca.Subject.Organization).To(HaveLen(1))
			Ω(ca.Subject.Organization[0]).To(Equal("test.org"))
		})
	})
})

func readKey(key []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(key)
	Ω(block).To(Not(BeNil()))

	k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	Ω(err).To(BeNil())
	Ω(k).To(Not(BeNil()))
	return k
}

func readCert(cert []byte) *x509.Certificate {
	block, _ := pem.Decode(cert)
	Ω(block).To(Not(BeNil()))

	c, err := x509.ParseCertificate(block.Bytes)
	Ω(err).To(BeNil())
	Ω(c).To(Not(BeNil()))
	return c
}
