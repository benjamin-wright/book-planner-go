package k8s

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type basicType interface {
	string | int | map[string]string
}

func getProperty[T basicType](u *unstructured.Unstructured, args ...string) (T, error) {
	var current interface{} = u.Object
	var empty T

	for _, arg := range args {
		property, ok := current.(map[string]interface{})
		if !ok {
			return empty, fmt.Errorf("property %s was not an object in %s", arg, strings.Join(args, "."))
		}

		value, ok := property[arg]
		if !ok {
			return empty, fmt.Errorf("object doesn't have property %s in %s", arg, strings.Join(args, "."))
		}

		current = value
	}

	value, ok := current.(T)
	if !ok {
		return empty, fmt.Errorf("property %s is not the right type: %T", strings.Join(args, "."), value)
	}

	return value, nil
}
