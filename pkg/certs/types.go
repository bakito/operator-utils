package certs

import "time"

const (
	// Dir directory of the certs
	Dir = "certs"
	// ServerKey is the name of the key associated with the secret's private key.
	ServerKey = "tls.key"
	// ServerCert is the name of the key associated with the secret's public key.
	ServerCert = "tls.crt"
	// CACert is the name of the key associated with the certificate of the CA for
	// the keypair.
	CACert = "ca.crt"
	// OneWeek Time used for updating a certificate before it expires.
	OneWeek = 7 * 24 * time.Hour
	// Organization Default cert organisation
	Organization = "cluster.local"
)

// Options cert options
type Options struct {
	CertDir                     string
	ServerKey                   string
	ServerCert                  string
	CACert                      string
	UpdateBefore                time.Duration
	Name                        string
	MutatingWebhookConfigName   string
	ValidatingWebhookConfigName string
	Organization                string
}

// ApplyDefaults apply default options
func (o *Options) ApplyDefaults(name string) Options {
	o.Name = name
	if o.CertDir == "" {
		o.CertDir = Dir
	}
	if o.ServerKey == "" {
		o.ServerKey = ServerKey
	}
	if o.ServerCert == "" {
		o.ServerCert = ServerCert
	}
	if o.CACert == "" {
		o.CACert = CACert
	}
	if o.UpdateBefore == 0 {
		o.UpdateBefore = OneWeek
	}
	if o.Organization == "" {
		o.Organization = Organization
	}
	if o.MutatingWebhookConfigName == "" {
		o.MutatingWebhookConfigName = o.Name
	}
	if o.ValidatingWebhookConfigName == "" {
		o.ValidatingWebhookConfigName = o.Name
	}
	return *o
}
