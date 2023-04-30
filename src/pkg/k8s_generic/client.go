package k8s_generic

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Resource[T any] interface {
	*T
	GetName() string
	ToUnstructured() *unstructured.Unstructured
	FromUnstructured(obj *unstructured.Unstructured) error
}

type Client[T any, PT Resource[T]] struct {
	client    dynamic.Interface
	namespace string
	schema    schema.GroupVersionResource
}

func Watch[T any, PT Resource[T]](ctx context.Context, cancel context.CancelFunc, schema schema.GroupVersionResource, namespace string) (<-chan map[string]T, error) {
	cli, err := New[T, PT](schema, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %+v", err)
	}

	stream, err := cli.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to start watch: %+v", err)
	}

	return stream, nil
}

func New[T any, PT Resource[T]](schema schema.GroupVersionResource, namespace string) (*Client[T, PT], error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client[T, PT]{
		namespace: namespace,
		schema:    schema,
		client:    dynClient,
	}, nil
}

func (c *Client[T, PT]) Create(ctx context.Context, resource T) error {
	ptr := PT(&resource)

	_, err := c.client.Resource(c.schema).Namespace(c.namespace).Create(ctx, ptr.ToUnstructured(), v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create %T: %+v", resource, err)
	}

	return nil
}

func (c *Client[T, PT]) Delete(ctx context.Context, name string) error {
	err := c.client.Resource(c.schema).Namespace(c.namespace).Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete %T: %+v", name, err)
	}

	return nil
}

func (c *Client[T, PT]) DeleteAll(ctx context.Context) error {
	err := c.client.Resource(c.schema).Namespace(c.namespace).DeleteCollection(ctx, v1.DeleteOptions{}, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete all resources: %+v", err)
	}

	return nil
}

func (c *Client[T, PT]) Watch(ctx context.Context, cancel context.CancelFunc) (<-chan map[string]T, error) {
	resources := map[string]T{}
	output := make(chan map[string]T, 1)
	convert := func(obj interface{}) (string, T) {
		var res T
		ptr := PT(&res)
		ptr.FromUnstructured(obj.(*unstructured.Unstructured))
		return ptr.GetName(), res
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.client, time.Minute, c.namespace, nil)
	informer := factory.ForResource(c.schema).Informer()

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
