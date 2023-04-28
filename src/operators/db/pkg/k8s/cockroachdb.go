package k8s

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type CockroachDB struct {
	Name      string
	Namespace string
}

func (db *CockroachDB) toUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "CockroachDB",
		"metadata": map[string]interface{}{
			"name": db.Name,
		},
		"spec": map[string]interface{}{},
	})

	return result
}

func (db *CockroachDB) fromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
}

func (db *CockroachDB) getName() string {
	return db.Name
}

func (c *Client) CockroachDBCreate(ctx context.Context, db CockroachDB) error {
	_, err := c.client.Resource(CockroachDBSchema).Namespace(db.Namespace).Create(ctx, db.toUnstructured(), v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cockroachdb: %+v", err)
	}

	return nil
}

type Resources[T watchable] struct {
	resources map[string]T
}

func (r *Resources[T]) Add(obj interface{}) {
	var res T
	res.fromUnstructured(obj.(*unstructured.Unstructured))
	r.resources[res.getName()] = res
}

func (r *Resources[T]) Delete(obj interface{}) {
	var res T
	res.fromUnstructured(obj.(*unstructured.Unstructured))
	delete(r.resources, res.getName())
}

type watchable interface {
	getName() string
	toUnstructured() *unstructured.Unstructured
	fromUnstructured(obj *unstructured.Unstructured)
}

func watchResource[T watchable](client dynamic.Interface, ctx context.Context, cancel context.CancelFunc, namespace string, schema schema.GroupVersionResource) (<-chan map[string]T, error) {
	resources := Resources[T]{}
	output := make(chan map[string]T, 1)

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, time.Minute, namespace, nil)
	informer := factory.ForResource(schema).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			resources.Add(obj)
			output <- resources.resources
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			resources.Add(newObj)
			output <- resources.resources
		},
		DeleteFunc: func(obj interface{}) {
			resources.Delete(obj)
			output <- resources.resources
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

func (c *Client) CockroachDBWatch(ctx context.Context, cancel context.CancelFunc, namespace string) (<-chan map[string]*CockroachDB, error) {
	return watchResource[*CockroachDB](c.client, ctx, cancel, namespace, CockroachDBSchema)
}
