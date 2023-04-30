package k8s

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
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

func (m *CockroachMigration) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	m.Name = obj.GetName()
	m.Namespace = obj.GetNamespace()
	m.Deployment, err = getProperty[string](obj, "spec", "deployment")
	if err != nil {
		return fmt.Errorf("failed to get deployment: %+v", err)
	}
	m.Database, err = getProperty[string](obj, "spec", "database")
	if err != nil {
		return fmt.Errorf("failed to get database: %+v", err)
	}
	m.Migration, err = getProperty[string](obj, "spec", "migration")
	if err != nil {
		return fmt.Errorf("failed to get migration: %+v", err)
	}
	m.Index, err = getProperty[int](obj, "spec", "index")
	if err != nil {
		return fmt.Errorf("failed to get index: %+v", err)
	}

	return nil
}

func (db *CockroachMigration) GetName() string {
	return db.Name
}

var CockroachMigrationSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachmigrations",
}

func NewCockroachMigrationClient(namespace string) (*k8s_generic.Client[CockroachMigration, *CockroachMigration], error) {
	return k8s_generic.New[CockroachMigration](CockroachMigrationSchema, namespace)
}

func WatchCockroachMigrations(ctx context.Context, cancel context.CancelFunc, namespace string) (<-chan map[string]CockroachMigration, error) {
	return k8s_generic.Watch[CockroachMigration](ctx, cancel, CockroachMigrationSchema, namespace)
}
