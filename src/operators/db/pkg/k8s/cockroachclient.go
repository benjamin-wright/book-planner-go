package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s/generic_client"
)

type CockroachClient struct {
	Name       string
	Namespace  string
	Deployment string
	Database   string
	Username   string
}

func (cli *CockroachClient) ToUnstructured() *unstructured.Unstructured {
	result := &unstructured.Unstructured{}
	result.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "ponglehub.co.uk/v1alpha1",
		"kind":       "CockroachClient",
		"metadata": map[string]interface{}{
			"name": cli.Name,
		},
		"spec": map[string]interface{}{
			"deployment": cli.Deployment,
			"database":   cli.Database,
			"username":   cli.Username,
		},
	})

	return result
}

func (db *CockroachClient) FromUnstructured(obj *unstructured.Unstructured) {
	db.Name = obj.GetName()
	db.Namespace = obj.GetNamespace()
}

func (db *CockroachClient) GetName() string {
	return db.Name
}

var CockroachClientSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachclients",
}

func NewCockroachClientClient(namespace string) (*generic_client.Client[CockroachClient, *CockroachClient], error) {
	return generic_client.New[CockroachClient](CockroachClientSchema, namespace)
}
