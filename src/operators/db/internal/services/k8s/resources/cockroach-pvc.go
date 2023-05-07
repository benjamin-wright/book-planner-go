package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/utils"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachPVC struct {
	Name     string
	Database string
	Storage  string
}

func (p *CockroachPVC) ToUnstructured(namespace string) *unstructured.Unstructured {
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}

	return pvc
}

func (p *CockroachPVC) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	p.Name = obj.GetName()
	p.Storage, err = utils.GetProperty[string](obj, "spec", "resources", "requests", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}

	p.Database, err = utils.GetProperty[string](obj, "metadata", "labels", "app")
	if err != nil {
		return fmt.Errorf("failed to get database from app label: %+v", err)
	}

	return nil
}

func (db *CockroachPVC) GetName() string {
	return db.Name
}

var CockroachPVCSchema = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "persistentvolumeclaims",
}

func NewCockroachPVCClient(namespace string) (*k8s_generic.Client[CockroachPVC, *CockroachPVC], error) {
	return k8s_generic.New[CockroachPVC](CockroachPVCSchema, namespace)
}
