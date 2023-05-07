package resources

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachService struct {
	Name string
}

func (s *CockroachService) ToUnstructured(namespace string) *unstructured.Unstructured {
	statefulset := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": s.Name,
				"labels": map[string]interface{}{
					"app": s.Name,
				},
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
	return nil
}

func (s *CockroachService) GetName() string {
	return s.Name
}

var CockroachServiceSchema = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "services",
}

func NewCockroachServiceClient(namespace string) (*k8s_generic.Client[CockroachService, *CockroachService], error) {
	return k8s_generic.New[CockroachService](CockroachServiceSchema, namespace)
}
