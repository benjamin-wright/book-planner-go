package k8s

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type watchOptions struct {
	namespace  string
	schema     schema.GroupVersionResource
	addFunc    func(obj *unstructured.Unstructured)
	updateFunc func(oldObj, newObj *unstructured.Unstructured)
	deleteFunc func(obj *unstructured.Unstructured)
}

func (c *Client) watchResource(ctx context.Context, cancel context.CancelFunc, options watchOptions) error {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.client, time.Minute, options.namespace, nil)
	informer := factory.ForResource(options.schema).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) { options.addFunc(obj.(*unstructured.Unstructured)) },
		UpdateFunc: func(oldObj, newObj interface{}) {
			options.updateFunc(oldObj.(*unstructured.Unstructured), newObj.(*unstructured.Unstructured))
		},
		DeleteFunc: func(obj interface{}) { options.deleteFunc(obj.(*unstructured.Unstructured)) },
	})

	go func() {
		informer.Run(ctx.Done())
		if ctx.Err() == nil {
			cancel()
		}
	}()

	return nil
}
