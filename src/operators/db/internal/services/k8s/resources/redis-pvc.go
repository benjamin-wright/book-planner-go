package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type RedisPVC struct {
	Name     string
	Database string
	Storage  string
}

func (p *RedisPVC) ToUnstructured(namespace string) *unstructured.Unstructured {
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}

	return pvc
}

func (p *RedisPVC) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	p.Name = obj.GetName()
	p.Storage, err = k8s_generic.GetProperty[string](obj, "spec", "resources", "requests", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}

	p.Database, err = k8s_generic.GetProperty[string](obj, "metadata", "labels", "app")
	if err != nil {
		return fmt.Errorf("failed to get database from app label: %+v", err)
	}

	return nil
}

func (db *RedisPVC) GetName() string {
	return db.Name
}

var Redisema = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "persistentvolumeclaims",
}

func NewRedisPVCClient(namespace string) (*k8s_generic.Client[RedisPVC, *RedisPVC], error) {
	return k8s_generic.New[RedisPVC](
		Redisema,
		namespace,
		k8s_generic.Merge(map[string]interface{}{
			"ponglehub.co.uk/resource-type": "redis",
		}, LABEL_FILTERS),
	)
}
