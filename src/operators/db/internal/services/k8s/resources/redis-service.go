package resources

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type RedisService struct {
	Name string
}

func (s *RedisService) ToUnstructured(namespace string) *unstructured.Unstructured {
	statefulset := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": s.Name,
				"labels": k8s_generic.Merge(map[string]interface{}{
					"app":                           s.Name,
					"ponglehub.co.uk/resource-type": "redis",
				}, LABEL_FILTERS),
			},
			"spec": map[string]interface{}{
				"ports": []map[string]interface{}{
					{
						"name":       "tcp",
						"port":       6379,
						"protocol":   "TCP",
						"targetPort": "tcp",
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

func (s *RedisService) FromUnstructured(obj *unstructured.Unstructured) error {
	s.Name = obj.GetName()
	return nil
}

func (s *RedisService) GetName() string {
	return s.Name
}

var RedisServiceSchema = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "services",
}

func NewRedisServiceClient(namespace string) (*k8s_generic.Client[RedisService, *RedisService], error) {
	return k8s_generic.New[RedisService](
		RedisServiceSchema,
		namespace,
		k8s_generic.Merge(map[string]interface{}{
			"ponglehub.co.uk/resource-type": "redis",
		}, LABEL_FILTERS),
	)
}
