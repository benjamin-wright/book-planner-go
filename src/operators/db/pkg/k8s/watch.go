package k8s

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type watchable[T any] interface {
	*T
	getName() string
	toUnstructured() *unstructured.Unstructured
	fromUnstructured(obj *unstructured.Unstructured)
}

func watchResource[T any, PT watchable[T]](client dynamic.Interface, ctx context.Context, cancel context.CancelFunc, namespace string, schema schema.GroupVersionResource) (<-chan map[string]T, error) {
	resources := map[string]T{}
	output := make(chan map[string]T, 1)
	convert := func(obj interface{}) (string, T) {
		var res T
		ptr := PT(&res)
		ptr.fromUnstructured(obj.(*unstructured.Unstructured))
		return ptr.getName(), res
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, time.Minute, namespace, nil)
	informer := factory.ForResource(schema).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			name, res := convert(obj)
			resources[name] = res

			output <- resources
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			name, res := convert(newObj)
			resources[name] = res

			output <- resources
		},
		DeleteFunc: func(obj interface{}) {
			name, _ := convert(obj)
			delete(resources, name)

			output <- resources
		},
	})

	go func() {
		informer.Run(ctx.Done())
		if ctx.Err() == nil {
			cancel()
		}
	}()

	return output, nil
}
