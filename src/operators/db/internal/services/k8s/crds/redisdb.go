package crds

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/utils"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type RedisDB struct {
	Name    string
	Storage string
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
	db.Storage, err = utils.GetProperty[string](obj, "spec", "storage")
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

func NewRedisDBClient(namespace string) (*k8s_generic.Client[RedisDB, *RedisDB], error) {
	return k8s_generic.New[RedisDB](RedisDBSchema, namespace)
}
