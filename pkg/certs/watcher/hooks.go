package watcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bakito/operator-utils/pkg/log"
	arv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (w *watcher) supportsV1() (bool, error) {
	dcl, err := discovery.NewDiscoveryClientForConfig(w.config)
	if err != nil {
		return false, err
	}
	apiList, err := dcl.ServerGroups()
	if err != nil {
		return false, err
	}
	return w.supportsV1Internal(apiList)
}

func (w *watcher) supportsV1Internal(apiList *metav1.APIGroupList) (bool, error) {
	for _, g := range apiList.Groups {
		if g.Name == arv1.GroupName {
			for _, v := range g.Versions {
				if v.Version == arv1.SchemeGroupVersion.Version {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, fmt.Errorf("could not find api group %q", arv1.GroupName)
}

func (w *watcher) patchWebhookConfig(ctx context.Context, whc client.Object, webhookNames []string, cert []byte) error {
	if len(webhookNames) == 0 {
		return nil
	}

	log.With(w.logger, whc).Info("Updating webhook ca cert")
	var webhooks []interface{}
	for _, name := range webhookNames {
		webhooks = append(webhooks, map[string]interface{}{
			"name": name,
			"clientConfig": map[string][]byte{
				"caBundle": cert,
			},
		})
	}

	patch := map[string]interface{}{
		"webhooks": webhooks,
	}

	mergePatch, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	err = w.client.Patch(ctx, whc, client.RawPatch(types.StrategicMergePatchType, mergePatch))
	return err
}
