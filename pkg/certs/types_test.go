package certs_test

import (
	"time"

	. "github.com/bakito/operator-utils/pkg/certs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Context("ApplyDefaults", func() {
		var (
			o    *Options
			name string
		)
		BeforeEach(func() {
			o = &Options{}
			name = "test-name"
			Ω(false).To(BeFalse())
		})

		It("nothing should be empty by default", func() {
			oo := o.ApplyDefaults(name)

			Ω(oo.CertDir).To(Equal(Dir))
			Ω(oo.ServerKey).To(Equal(ServerKey))
			Ω(oo.ServerCert).To(Equal(ServerCert))
			Ω(oo.CACert).To(Equal(CACert))
			Ω(oo.UpdateBefore).To(Equal(OneWeek))
			Ω(oo.Organization).To(Equal(Organization))

			Ω(oo.Name).To(Equal(name))
			Ω(oo.MutatingWebhookConfigName).To(Equal(name))
			Ω(oo.ValidatingWebhookConfigName).To(Equal(name))
		})

		It("not overwritten if set", func() {
			o.CertDir = "CertDir"
			o.ServerKey = "ServerKey"
			o.ServerCert = "ServerCert"
			o.CACert = "CACert"
			o.UpdateBefore = 1 * time.Hour
			o.Name = "old-name"
			o.MutatingWebhookConfigName = "MutatingWebhookConfigName"
			o.ValidatingWebhookConfigName = "ValidatingWebhookConfigName"
			o.Organization = "Organization"
			oo := o.ApplyDefaults(name)

			Ω(oo.CertDir).To(Equal("CertDir"))
			Ω(oo.ServerKey).To(Equal("ServerKey"))
			Ω(oo.ServerCert).To(Equal("ServerCert"))
			Ω(oo.CACert).To(Equal("CACert"))
			Ω(oo.UpdateBefore).To(Equal(1 * time.Hour))
			Ω(oo.Organization).To(Equal("Organization"))

			Ω(oo.Name).To(Equal(name))
			Ω(oo.MutatingWebhookConfigName).To(Equal("MutatingWebhookConfigName"))
			Ω(oo.ValidatingWebhookConfigName).To(Equal("ValidatingWebhookConfigName"))
		})
	})
})
