package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s/generic_client"
)

type RedisDB struct {
	Name      string
	Namespace string
}

func (db *RedisDB) ToUnstructured() *unstructured.Unstructured {
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

func (db *RedisDB) FromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
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
