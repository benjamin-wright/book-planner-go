package crds

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachMigration struct {
	Name       string
	Deployment string
	Database   string
	Migration  string
	Index      int64
}

func (cm *CockroachMigration) ToUnstructured(namespace string) *unstructured.Unstructured {
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
	m.Deployment, err = k8s_generic.GetProperty[string](obj, "spec", "deployment")
	if err != nil {
		return fmt.Errorf("failed to get deployment: %+v", err)
	}
	m.Database, err = k8s_generic.GetProperty[string](obj, "spec", "database")
	if err != nil {
		return fmt.Errorf("failed to get database: %+v", err)
	}
	m.Migration, err = k8s_generic.GetProperty[string](obj, "spec", "migration")
	if err != nil {
		return fmt.Errorf("failed to get migration: %+v", err)
	}
	m.Index, err = k8s_generic.GetProperty[int64](obj, "spec", "index")
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
	return k8s_generic.New[CockroachMigration](CockroachMigrationSchema, namespace, nil)
}
