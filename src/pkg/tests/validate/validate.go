package validate

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var UUID_REGEX = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func Equal[T any](t *testing.T, expected T, actual T) bool {
	typ := reflect.TypeOf(expected)

	typ.Kind()
	if !assert.Equal(t, typ.Kind(), reflect.Struct, "can't validate non-structs") {
		return false
	}

	expectedVal := reflect.ValueOf(&expected).Elem()
	actualVal := reflect.ValueOf(&actual).Elem()

	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			assert.Equal(t, expectedVal.FieldByName(field.Name).Interface(), actualVal.FieldByName(field.Name).Interface(), "fields should match: %s", field.Name)
			continue
		}

		keys := strings.Split(tag, ",")
		for _, key := range keys {
			switch key {
			case "uuid":
				value := actualVal.FieldByName(field.Name)
				if value.Kind() != reflect.String {
					assert.Fail(t, "non-string field %s (%s) marked as a UUID for validation", field.Name, value.Kind())
					continue
				}

				expectedVal.FieldByName(field.Name).SetString("a valid UUID")

				valueStr := value.String()
				if UUID_REGEX.Match([]byte(valueStr)) {
					value.SetString("a valid UUID")
				}
			}
		}
	}

	return assert.Equal(t, expected, actual)
}
