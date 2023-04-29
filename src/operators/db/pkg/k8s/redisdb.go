package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type RedisDB struct {
	Name      string
	Namespace string
}

func (db *RedisDB) toUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "RedisDB",
		"metadata": map[string]interface{}{
			"name": db.Name,
		},
		"spec": map[string]interface{}{},
	})

	return result
}

func (db *RedisDB) fromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
}

func (db *RedisDB) getName() string {
	return db.Name
}

var RedisDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "redisdbs",
}

func (c *Client) RedisDBCreate(ctx context.Context, db CockroachDB) error {
	_, err := c.client.Resource(RedisDBSchema).Namespace(db.Namespace).Create(ctx, db.toUnstructured(), v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create redisdb: %+v", err)
	}

	return nil
}

func (c *Client) RedisDBWatch(ctx context.Context, cancel context.CancelFunc, namespace string) (<-chan map[string]RedisDB, error) {
	return watchResource[RedisDB](c.client, ctx, cancel, namespace, RedisDBSchema)
}
