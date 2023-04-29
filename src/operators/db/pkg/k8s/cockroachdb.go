package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

var CockroachDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachdbs",
}

func (c *Client) CockroachDBCreate(ctx context.Context, db CockroachDB) error {
	_, err := c.client.Resource(CockroachDBSchema).Namespace(db.Namespace).Create(ctx, db.toUnstructured(), v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cockroachdb: %+v", err)
	}

	return nil
}

func (c *Client) CockroachDBWatch(ctx context.Context, cancel context.CancelFunc, namespace string) (<-chan map[string]CockroachDB, error) {
	return watchResource[CockroachDB](c.client, ctx, cancel, namespace, CockroachDBSchema)
}
