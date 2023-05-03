package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/utils"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachStatefulSet struct {
	Name    string
	Storage string
	CPU     string
	Memory  string
}

func (s *CockroachStatefulSet) ToUnstructured() *unstructured.Unstructured {
	statefulset := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "StatefulSet",
			"metadata": map[string]interface{}{
				"name": s.Name,
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": s.Name,
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": s.Name,
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "database",
								"image": "cockroachdb/cockroach:v20.2.8",
								"command": []string{
									"cockroach",
								},
								"args": []string{
									"--logtostderr",
									"start-single-node",
									"--insecure",
								},
								"resources": map[string]interface{}{
									"requests": map[string]interface{}{
										"cpu":    s.CPU,
										"memory": s.Memory,
									},
									"limits": map[string]interface{}{
										"cpu":    s.CPU,
										"memory": s.Memory,
									},
								},
								"volumes": []map[string]interface{}{
									{
										"name": "datadir",
										"persistentVolumeClaim": map[string]interface{}{
											"claimName": "datadir",
										},
									},
								},
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 8080,
									},
									{
										"name":          "grpc",
										"protocol":      "TCP",
										"containerPort": 26257,
									},
								},
								"readinessProbe": map[string]interface{}{
									"httpGet": map[string]interface{}{
										"path":   "/health?ready=1",
										"port":   "http",
										"scheme": "HTTPS",
									},
									"initialDelaySeconds": 10,
									"periodSeconds":       5,
									"failureThreshold":    2,
								},
							},
						},
					},
				},
				"volumeClaimTemplates": []map[string]interface{}{
					{
						"metadata": map[string]interface{}{
							"name": "datadir",
						},
						"spec": map[string]interface{}{
							"accessModes": []string{
								"ReadWriteOnce",
							},
							"resources": map[string]interface{}{
								"requests": map[string]interface{}{
									"storage": s.Storage,
								},
							},
						},
					},
				},
			},
		},
	}

	return statefulset
}

func (s *CockroachStatefulSet) FromUnstructured(obj *unstructured.Unstructured) error {
	var err error

	s.Name = obj.GetName()
	s.Storage, err = utils.GetProperty[string](obj, "spec", "volumeClaimTemplates", "0", "spec", "resources", "requests", "storage")
	if err != nil {
		return fmt.Errorf("failed to get storage: %+v", err)
	}
	s.CPU, err = utils.GetProperty[string](obj, "spec", "template", "spec", "containers", "0", "resources", "requests", "cpu")
	if err != nil {
		return fmt.Errorf("failed to get cpu: %+v", err)
	}
	s.Memory, err = utils.GetProperty[string](obj, "spec", "template", "spec", "containers", "0", "resources", "requests", "memory")
	if err != nil {
		return fmt.Errorf("failed to get memory: %+v", err)
	}

	return nil
}

func (db *CockroachStatefulSet) GetName() string {
	return db.Name
}

var CockroachStatefulSetSchema = schema.GroupVersionResource{
	Group:    "apps",
	Version:  "v1",
	Resource: "statefulsets",
}

func NewCockroachStatefulSetClient(namespace string) (*k8s_generic.Client[CockroachStatefulSet, *CockroachStatefulSet], error) {
	return k8s_generic.New[CockroachStatefulSet](CockroachStatefulSetSchema, namespace)
}
