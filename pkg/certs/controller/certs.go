package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// Create the common parts of the cert. These don't change between
// the root/CA cert and the server cert.
func (r *reconciler) createCertTemplate(notAfter time.Time) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	serviceName := r.nn.Name + "." + r.nn.Namespace
	commonName := serviceName + ".svc"
	serviceNames := []string{
		r.nn.Name,
		serviceName,
		commonName,
		serviceName + ".svc.cluster.local",
	}

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{r.opts.Organization},
			CommonName:   commonName,
		},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		DNSNames:              serviceNames,
	}
	return &tmpl, nil
}

// Create cert template suitable for CA and hence signing
func (r *reconciler) createCACertTemplate(notAfter time.Time) (*x509.Certificate, error) {
	rootCert, err := r.createCertTemplate(notAfter)
	if err != nil {
		return nil, err
	}
	// Make it into a CA cert and change it so we can use it to sign certs
	rootCert.IsCA = true
	rootCert.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	return rootCert, nil
}

// Create cert template that we can use on the server for TLS
func (r *reconciler) createServerCertTemplate(notAfter time.Time) (*x509.Certificate, error) {
	serverCert, err := r.createCertTemplate(notAfter)
	if err != nil {
		return nil, err
	}
	serverCert.KeyUsage = x509.KeyUsageDigitalSignature
	serverCert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	return serverCert, err
}

// Actually sign the cert and return things in a form that we can use later on
func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error,
) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func (r *reconciler) createCA(notAfter time.Time) (*rsa.PrivateKey, *x509.Certificate, []byte, error) {
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating random key: %w", err)
	}

	rootCertTmpl, err := r.createCACertTemplate(notAfter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating CA cert: %w", err)
	}

	rootCert, rootCertPEM, err := createCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error signing the CA cert: %w", err)
	}
	return rootKey, rootCert, rootCertPEM, nil
}

// CreateCerts creates and returns a CA certificate and certificate and
// key for the server. serverKey and serverCert are used by the server
// to establish trust for clients, CA certificate is used by the
// client to verify the server authentication chain. notAfter specifies
// the expiration date.
func (r *reconciler) createCerts(notAfter time.Time) (serverKey, serverCert, caCert []byte, err error) {
	// First create a CA certificate and private key
	caKey, caCertificate, caCertificatePEM, err := r.createCA(notAfter)
	if err != nil {
		return nil, nil, nil, err
	}

	// Then create the private key for the serving cert
	servKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error generating random key: %w", err)
	}
	servCertTemplate, err := r.createServerCertTemplate(notAfter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create the server certificate template: %w", err)
	}

	// create a certificate which wraps the server's public key, sign it with the CA private key
	_, servCertPEM, err := createCert(servCertTemplate, caCertificate, &servKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error signing server certificate template: %w", err)
	}
	servKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(servKey),
	})
	return servKeyPEM, servCertPEM, caCertificatePEM, nil
}
