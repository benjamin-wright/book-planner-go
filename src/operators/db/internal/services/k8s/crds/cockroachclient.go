package crds

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachClient struct {
	Name       string
	Deployment string
	Database   string
	Username   string
	Secret     string
}

func (cli *CockroachClient) ToUnstructured(namespace string) *unstructured.Unstructured {
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
			"secret":     cli.Secret,
		},
	})

	return result
}

func (db *CockroachClient) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	db.Name = obj.GetName()
	db.Deployment, err = k8s_generic.GetProperty[string](obj, "spec", "deployment")
	if err != nil {
		return fmt.Errorf("failed to get deployment: %+v", err)
	}

	db.Database, err = k8s_generic.GetProperty[string](obj, "spec", "database")
	if err != nil {
		return fmt.Errorf("failed to get database: %+v", err)
	}

	db.Username, err = k8s_generic.GetProperty[string](obj, "spec", "username")
	if err != nil {
		return fmt.Errorf("failed to get username: %+v", err)
	}

	db.Secret, err = k8s_generic.GetProperty[string](obj, "spec", "secret")
	if err != nil {
		return fmt.Errorf("failed to get secret: %+v", err)
	}

	return nil
}

func (db *CockroachClient) GetName() string {
	return db.Name
}

func (db *CockroachClient) GetTarget() string {
	return db.Deployment
}

var CockroachClientSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachclients",
}

func NewCockroachClientClient(namespace string) (*k8s_generic.Client[CockroachClient, *CockroachClient], error) {
	return k8s_generic.New[CockroachClient](CockroachClientSchema, namespace, nil)
}
