package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s/generic_client"
)

type CockroachMigration struct {
	Name       string
	Namespace  string
	Deployment string
	Database   string
	Migration  string
	Index      int
}

func (cm *CockroachMigration) ToUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "CockroachMigration",
		"metadata": map[string]interface{}{
			"name": cm.Name,
		},
		"spec": map[string]interface{}{
			"deployment": cm.Deployment,
			"database":   cm.Database,
			"migration":  cm.Migration,
			"index":      cm.Index,
		},
	})

	return result
}

func (db *CockroachMigration) FromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
}

func (db *CockroachMigration) GetName() string {
	return db.Name
}

var CockroachMigrationSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachmigrations",
}

func NewCockroachMigrationClient(namespace string) (*generic_client.Client[CockroachMigration, *CockroachMigration], error) {
	return generic_client.New[CockroachMigration](CockroachMigrationSchema, namespace)
}
