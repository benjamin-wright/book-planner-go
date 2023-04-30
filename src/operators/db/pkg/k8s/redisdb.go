package k8s

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s/generic_client"
)

type RedisDB struct {
	Name      string
	Namespace string
	Storage   string
}

func (db *RedisDB) ToUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "RedisDB",
		"metadata": map[string]interface{}{
			"name": db.Name,
		},
		"spec": map[string]interface{}{
			"storage": db.Storage,
		},
	})

	return result
}

func (db *RedisDB) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
	db.Storage, err = getProperty[string](obj, "spec", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}

	return nil
}

func (db *RedisDB) GetName() string {
	return db.Name
}

var RedisDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "redisdbs",
}

func NewRedisDBClient(namespace string) (*generic_client.Client[RedisDB, *RedisDB], error) {
	return generic_client.New[RedisDB](RedisDBSchema, namespace)
}

func WatchRedisDBs(ctx context.Context, cancel context.CancelFunc, namespace string) (<-chan map[string]RedisDB, error) {
	return generic_client.Watch[RedisDB](ctx, cancel, RedisDBSchema, namespace)
}
