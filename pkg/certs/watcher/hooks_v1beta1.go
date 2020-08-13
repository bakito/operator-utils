package watcher

import (
	"bytes"
	"context"

	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

func (w *watcher) patchHooksBeta1V1(ctx context.Context, caCert []byte) error {
	mwc := &arv1beta1.MutatingWebhookConfiguration{}
	err := w.client.Get(ctx, types.NamespacedName{Name: w.opts.MutatingWebhookConfigName}, mwc)
	if err != nil {
		return err
	}

	var mNames []string
	for i := range mwc.Webhooks {
		if !bytes.Equal(mwc.Webhooks[i].ClientConfig.CABundle, caCert) {
			mNames = append(mNames, mwc.Webhooks[i].Name)
		}
		mwc.Webhooks[i].ClientConfig.CABundle = caCert
	}

	err = w.patchWebhookConfig(ctx, mwc, mNames, caCert)
	if err != nil {
		return err
	}

	vwc := &arv1beta1.ValidatingWebhookConfiguration{}
	err = w.client.Get(ctx, types.NamespacedName{Name: w.opts.ValidatingWebhookConfigName}, vwc)
	if err != nil {
		return err
	}

	var vNames []string
	for i := range vwc.Webhooks {
		if !bytes.Equal(vwc.Webhooks[i].ClientConfig.CABundle, caCert) {
			vNames = append(vNames, vwc.Webhooks[i].Name)
		}
		vwc.Webhooks[i].ClientConfig.CABundle = caCert
	}

	err = w.patchWebhookConfig(ctx, vwc, vNames, caCert)
	if err != nil {
		return err
	}

	return nil
}
