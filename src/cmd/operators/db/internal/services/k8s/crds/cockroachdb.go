package crds

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachDBComparable struct {
	Name    string
	Storage string
}

type CockroachDB struct {
	CockroachDBComparable
	UID             string
	ResourceVersion string
}

func (db *CockroachDB) ToUnstructured(namespace string) *unstructured.Unstructured {
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
	db.UID = string(obj.GetUID())
	db.ResourceVersion = obj.GetResourceVersion()
	db.Storage, err = k8s_generic.GetProperty[string](obj, "spec", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}

	return nil
}

func (db *CockroachDB) GetName() string {
	return db.Name
}

func (db *CockroachDB) GetStorage() string {
	return db.Storage
}

func (db *CockroachDB) GetUID() string {
	return db.UID
}

func (db *CockroachDB) GetResourceVersion() string {
	return db.ResourceVersion
}

func (db *CockroachDB) Equal(obj CockroachDB) bool {
	return db.CockroachDBComparable == obj.CockroachDBComparable
}

func NewCockroachDBClient(namespace string) (*k8s_generic.Client[CockroachDB, *CockroachDB], error) {
	return k8s_generic.New[CockroachDB](
		schema.GroupVersionResource{
			Group:    "ponglehub.co.uk",
			Version:  "v1alpha1",
			Resource: "cockroachdbs",
		},
		"CockroachDB",
		namespace,
		nil,
	)
}
