package controller

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"time"

	"github.com/bakito/operator-utils/pkg/certs"
	"github.com/bakito/operator-utils/pkg/log"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler interface
type Reconciler interface {
	SetupWithManager(globalMgr, namespacedMgr ctrl.Manager) error
}

// reconciler reconciles a ClusterRole object
type reconciler struct {
	client.Client
	log  logr.Logger
	nn   types.NamespacedName
	opts certs.Options
}

func (r *reconciler) logger() logr.Logger {
	return r.log.WithValues("certs", r.nn)
}

// +kubebuilder:rbac:groups=,resources=secret,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations;validatingwebhookconfigurations,verbs=get;list;watch;create;update;patch

func (r *reconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	certLog := r.logger()

	// Fetch the ClusterRole instance
	secret := &corev1.Secret{}
	err := r.Get(ctx, r.nn, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			certLog.Error(err, "could not find cert secret")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	recreate := true

	if _, haskey := secret.Data[r.opts.ServerKey]; !haskey {
		certLog.WithValues("cert", r.opts.ServerKey).Info("Certificate secret is missing key")
	} else if _, haskey := secret.Data[r.opts.ServerCert]; !haskey {
		certLog.WithValues("cert", r.opts.ServerCert).Info("Certificate secret is missing key")
	} else if _, haskey := secret.Data[r.opts.CACert]; !haskey {
		certLog.WithValues("cert", r.opts.CACert).Info("Certificate secret is missing key")
	} else {
		// Check the expiration date of the certificate to see if it needs to be updated
		cert, err := tls.X509KeyPair(secret.Data[r.opts.ServerCert], secret.Data[r.opts.ServerKey])
		if err != nil {
			certLog.Error(err, "Error creating pem from certificate and key")
		} else {
			certData, err := x509.ParseCertificate(cert.Certificate[0])
			if err != nil {
				certLog.Error(err, "Error parsing certificate")
			} else if time.Now().Add(r.opts.UpdateBefore).Before(certData.NotAfter) {
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
			r.opts.ServerKey:  serverKey,
			r.opts.ServerCert: serverCert,
			r.opts.CACert:     caCert,
		}
		err = r.patchSecret(ctx, secret)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *reconciler) patchSecret(ctx context.Context, secret *corev1.Secret) error {
	patch := map[string]interface{}{
		"data": secret.Data,
	}

	mergePatch, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	return r.Patch(ctx, secret, client.RawPatch(types.MergePatchType, mergePatch))
}
