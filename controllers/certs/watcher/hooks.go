package watcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bakito/operator-utils/log"
	"k8s.io/apimachinery/pkg/runtime"
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
	for _, g := range apiList.Groups {
		if g.Name == apiGroup {
			for _, v := range g.Versions {
				if v.Version == "v1" {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, fmt.Errorf("could not find api group %q", apiGroup)
}

func (w *watcher) patchWebhookConfig(ctx context.Context, whc runtime.Object, webhookNames []string, cert []byte) error {
	if len(webhookNames) == 0 {
		return nil
	}

	log.With(w.logger, whc).Info("Updating CA cert")
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
