package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s/generic_client"
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

func (db *CockroachDB) FromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()

	spec := obj.Object["spec"]
	storage := spec.(map[string]interface{})["storage"]
	db.Storage = storage.(string)
}

func (db *CockroachDB) GetName() string {
	return db.Name
}

var CockroachDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachdbs",
}

func NewCockroachDBClient(namespace string) (*generic_client.Client[CockroachDB, *CockroachDB], error) {
	return generic_client.New[CockroachDB](CockroachDBSchema, namespace)
}
