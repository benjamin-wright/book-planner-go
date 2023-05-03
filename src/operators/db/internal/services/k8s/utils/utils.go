package utils

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type BasicType interface {
	string | int | map[string]string
}

func GetProperty[T BasicType](u *unstructured.Unstructured, args ...string) (T, error) {
	var current interface{} = u.Object
	var empty T

	for _, arg := range args {
		switch c := current.(type) {
		case map[string]interface{}:
			value, ok := c[arg]
			if !ok {
				return empty, fmt.Errorf("object doesn't have property %s in %s", arg, strings.Join(args, "."))
			}

			current = value
		case []interface{}:
			idx, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				return empty, fmt.Errorf("attempt to index array with non-numeric key %s in %s", arg, strings.Join(args, "."))
			}

			if idx < 0 || idx >= int64(len(c)) {
				return empty, fmt.Errorf("attempt to index array with non-numeric key %s in %s", arg, strings.Join(args, "."))
			}

			current = c[idx]
		default:
			return empty, fmt.Errorf("property %s was not an object in %s", arg, strings.Join(args, "."))
		}
	}

	value, ok := current.(T)
	if !ok {
		return empty, fmt.Errorf("property %s is not the right type: %T", strings.Join(args, "."), value)
	}

	return value, nil
}
