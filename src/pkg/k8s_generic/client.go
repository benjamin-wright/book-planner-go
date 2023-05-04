package k8s_generic

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
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

func (c *Client[T, PT]) Update(ctx context.Context, resource T) error {
	ptr := PT(&resource)

	_, err := c.client.Resource(c.schema).Namespace(c.namespace).Update(ctx, ptr.ToUnstructured(), v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed up update resource %s: %+v", ptr.GetName(), err)
	}

	return nil
}

type Update[T any] struct {
	ToAdd    []T
	ToRemove []T
}

func (c *Client[T, PT]) Watch(ctx context.Context, cancel context.CancelFunc) (<-chan Update[T], error) {
	output := make(chan Update[T], 1)
	convert := func(obj interface{}) T {
		var res T
		ptr := PT(&res)
		err := ptr.FromUnstructured(obj.(*unstructured.Unstructured))
		if err != nil {
			zap.S().Errorf("Failed to parse unstructured obj for %T: %+v\n %+v", res, err, obj)
		}
		return res
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.client, time.Minute, c.namespace, nil)
	informer := factory.ForResource(c.schema).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			res := convert(obj)

			output <- Update[T]{
				ToAdd:    []T{res},
				ToRemove: []T{},
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldRes := convert(oldObj)
			newRes := convert(newObj)

			output <- Update[T]{
				ToAdd:    []T{newRes},
				ToRemove: []T{oldRes},
			}
		},
		DeleteFunc: func(obj interface{}) {
			res := convert(obj)

			output <- Update[T]{
				ToAdd:    []T{},
				ToRemove: []T{res},
			}
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
