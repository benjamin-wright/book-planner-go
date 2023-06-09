package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type RedisPVCComparable struct {
	Name     string
	Database string
	Storage  string
}

type RedisPVC struct {
	RedisPVCComparable
	UID             string
	ResourceVersion string
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
	p.UID = string(obj.GetUID())
	p.ResourceVersion = obj.GetResourceVersion()
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

func (p *RedisPVC) GetName() string {
	return p.Name
}

func (p *RedisPVC) GetUID() string {
	return p.UID
}

func (p *RedisPVC) GetResourceVersion() string {
	return p.ResourceVersion
}

func (p *RedisPVC) Equal(obj RedisPVC) bool {
	return p.RedisPVCComparable == obj.RedisPVCComparable
}

func NewRedisPVCClient(namespace string) (*k8s_generic.Client[RedisPVC, *RedisPVC], error) {
	return k8s_generic.New[RedisPVC](
		schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "persistentvolumeclaims",
		},
		"PersistentVolumeClaim",
		namespace,
		k8s_generic.Merge(map[string]interface{}{
			"ponglehub.co.uk/resource-type": "redis",
		}, LABEL_FILTERS),
	)
}
