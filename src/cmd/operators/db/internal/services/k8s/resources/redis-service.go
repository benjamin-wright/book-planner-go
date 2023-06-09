package resources

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type RedisServiceComparable struct {
	Name string
}

type RedisService struct {
	RedisServiceComparable
	UID             string
	ResourceVersion string
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
	s.UID = string(obj.GetUID())
	s.ResourceVersion = obj.GetResourceVersion()
	return nil
}

func (s *RedisService) GetName() string {
	return s.Name
}

func (s *RedisService) GetUID() string {
	return s.UID
}

func (s *RedisService) GetResourceVersion() string {
	return s.ResourceVersion
}

func (s *RedisService) Equal(obj RedisService) bool {
	return s.RedisServiceComparable == obj.RedisServiceComparable
}

func NewRedisServiceClient(namespace string) (*k8s_generic.Client[RedisService, *RedisService], error) {
	return k8s_generic.New[RedisService](
		schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "services",
		},
		"Service",
		namespace,
		k8s_generic.Merge(map[string]interface{}{
			"ponglehub.co.uk/resource-type": "redis",
		}, LABEL_FILTERS),
	)
}
