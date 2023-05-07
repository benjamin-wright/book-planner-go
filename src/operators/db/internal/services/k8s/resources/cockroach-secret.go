package resources

import (
	"encoding/base64"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type CockroachSecret struct {
	Name     string
	DB       string
	Database string
	User     string
}

func encode(format string, args ...interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf(format, args...),
	))
}

func (s *CockroachSecret) ToUnstructured(namespace string) *unstructured.Unstructured {
	secret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name": s.Name,
				"labels": k8s_generic.Merge(map[string]interface{}{
					"app": s.Name,
				}, LABEL_FILTERS),
			},
			"data": map[string]interface{}{
				"POSTGRES_HOST": encode("%s.%s.svc.cluster.local", s.DB, namespace),
				"POSTGRES_PORT": encode("26257"),
				"POSTGRES_USER": encode(s.User),
				"POSTGRES_NAME": encode(s.Database),
			},
		},
	}

	return secret
}

func (s *CockroachSecret) FromUnstructured(obj *unstructured.Unstructured) error {
	s.Name = obj.GetName()

	hostname, err := k8s_generic.GetEncodedProperty(obj, "data", "POSTGRES_HOST")
	if err != nil {
		return fmt.Errorf("failed to get POSTGRES_HOST: %+v", err)
	}
	s.DB = strings.Split(hostname, ".")[0]

	s.User, err = k8s_generic.GetEncodedProperty(obj, "data", "POSTGRES_USER")
	if err != nil {
		return fmt.Errorf("failed to get POSTGRES_USER: %+v", err)
	}

	s.Database, err = k8s_generic.GetEncodedProperty(obj, "data", "POSTGRES_NAME")
	if err != nil {
		return fmt.Errorf("failed to get POSTGRES_NAME: %+v", err)
	}

	return nil
}

func (s *CockroachSecret) GetName() string {
	return s.Name
}

var CockroachSecretSchema = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "secrets",
}

func NewCockroachSecretClient(namespace string) (*k8s_generic.Client[CockroachSecret, *CockroachSecret], error) {
	return k8s_generic.New[CockroachSecret](CockroachSecretSchema, namespace, LABEL_FILTERS)
}
