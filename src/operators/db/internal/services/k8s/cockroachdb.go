package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachDB struct {
	Name      string
	Namespace string
	Storage   string
}

func (db *CockroachDB) ToUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "CockroachDB",
		"metadata": map[string]interface{}{
			"name": db.Name,
		},
		"spec": map[string]interface{}{
			"storage": db.Storage,
		},
	})

	return result
}

func (db *CockroachDB) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
	db.Storage, err = getProperty[string](obj, "spec", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}

	return nil
}

func (db *CockroachDB) GetName() string {
	return db.Name
}

var CockroachDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachdbs",
}

func NewCockroachDBClient(namespace string) (*k8s_generic.Client[CockroachDB, *CockroachDB], error) {
	return k8s_generic.New[CockroachDB](CockroachDBSchema, namespace)
}
