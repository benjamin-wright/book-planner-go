package resources

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachServiceComparable struct {
	Name string
}

type CockroachService struct {
	CockroachServiceComparable
	UID             string
	ResourceVersion string
}

func (s *CockroachService) ToUnstructured(namespace string) *unstructured.Unstructured {
	statefulset := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": s.Name,
				"labels": k8s_generic.Merge(map[string]interface{}{
					"app":                           s.Name,
					"ponglehub.co.uk/resource-type": "cockroachdb",
				}, LABEL_FILTERS),
			},
			"spec": map[string]interface{}{
				"ports": []map[string]interface{}{
					{
						"name":       "grpc",
						"port":       26257,
						"protocol":   "TCP",
						"targetPort": "grpc",
					},
				},
				"selector": map[string]interface{}{
					"app": s.Name,
				},
			},
		},
	}

	return statefulset
}

func (s *CockroachService) FromUnstructured(obj *unstructured.Unstructured) error {
	s.Name = obj.GetName()
	s.UID = string(obj.GetUID())
	s.ResourceVersion = obj.GetResourceVersion()
	return nil
}

func (s *CockroachService) GetName() string {
	return s.Name
}

func (s *CockroachService) GetUID() string {
	return s.UID
}

func (s *CockroachService) GetResourceVersion() string {
	return s.ResourceVersion
}

func (s *CockroachService) Equal(obj CockroachService) bool {
	return s.CockroachServiceComparable == obj.CockroachServiceComparable
}

func NewCockroachServiceClient(namespace string) (*k8s_generic.Client[CockroachService, *CockroachService], error) {
	return k8s_generic.New[CockroachService](
		schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "services",
		},
		"Service",
		namespace,
		k8s_generic.Merge(map[string]interface{}{
			"ponglehub.co.uk/resource-type": "cockroachdb",
		}, LABEL_FILTERS),
	)
}
