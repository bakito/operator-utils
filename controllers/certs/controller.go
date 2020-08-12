package certs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"time"

	"github.com/go-logr/logr"
	"github.com/bakito/operator-utils/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles a ClusterRole object
type Reconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	NamespacedName types.NamespacedName
}

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
	// Time used for updating a certificate before it expires.
	oneWeek = 7 * 24 * time.Hour
)

func (r *Reconciler) logger() logr.Logger {
	return r.Log.WithValues("certs", r.NamespacedName)
}

// +kubebuilder:rbac:groups=,resources=secret,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations;validatingwebhookconfigurations,verbs=get;list;watch;create;update;patch

func (r *Reconciler) Reconcile(ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	certLog := r.logger()

	// Fetch the ClusterRole instance
	secret := &corev1.Secret{}
	err := r.Get(ctx, r.NamespacedName, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			certLog.Error(err, "could not find cert secret")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	recreate := true

	if _, haskey := secret.Data[ServerKey]; !haskey {
		certLog.WithValues("cert", ServerKey).Info("Certificate secret is missing key")
	} else if _, haskey := secret.Data[ServerCert]; !haskey {
		certLog.WithValues("cert", ServerCert).Info("Certificate secret is missing key")
	} else if _, haskey := secret.Data[CACert]; !haskey {
		certLog.WithValues("cert", CACert).Info("Certificate secret is missing key")
	} else {
		// Check the expiration date of the certificate to see if it needs to be updated
		cert, err := tls.X509KeyPair(secret.Data[ServerCert], secret.Data[ServerKey])
		if err != nil {
			certLog.Error(err, "Error creating pem from certificate and key")
		} else {
			certData, err := x509.ParseCertificate(cert.Certificate[0])
			if err != nil {
				certLog.Error(err, "Error parsing certificate")
			} else if time.Now().Add(oneWeek).Before(certData.NotAfter) {
				recreate = false
			}
		}
	}

	if recreate {
		l := log.With(certLog, secret)
		l.Info("Recreating certificates")
		serverKey, serverCert, caCert, err := r.createCerts(time.Now().AddDate(1, 0, 0))
		if err != nil {
			return reconcile.Result{}, err
		}

		secret.Data = map[string][]byte{
			ServerKey:  serverKey,
			ServerCert: serverCert,
			CACert:     caCert,
		}
		err = r.patchSecret(ctx, secret)
	}
	return reconcile.Result{}, err
}

func (r *Reconciler) patchSecret(ctx context.Context, secret *corev1.Secret) error {
	patch := map[string]interface{}{
		"data": secret.Data,
	}

	mergePatch, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	return r.Patch(ctx, secret, client.RawPatch(types.MergePatchType, mergePatch))
}
