package validate

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var UUID_REGEX = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func formatArray(t *testing.T, typ reflect.StructField, value reflect.Value, expected bool) {
	for i := 0; i < value.Len(); i++ {
		element := value.Index(i)

		switch element.Type().Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			formatArray(t, typ, element, expected)
		case reflect.Struct:
			if element.Type() == reflect.TypeOf(time.Time{}) {
				formatTimeField(t, typ, element, expected)
			} else {
				formatStruct(t, element, expected)
			}
		case reflect.String:
			formatStringField(t, typ, element, expected)
		}
	}
}

func formatStruct(t *testing.T, value reflect.Value, expected bool) {
	if !value.IsValid() {
		return
	}

	typ := value.Type()
	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)

		switch field.Type.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			formatArray(t, field, value.FieldByName(field.Name), expected)
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				formatTimeField(t, field, value.FieldByName(field.Name), expected)
			} else {
				formatStruct(t, value.FieldByName(field.Name), expected)
			}
		case reflect.String:
			formatStringField(t, field, value.FieldByName(field.Name), expected)
		}
	}
}

func formatStringField(t *testing.T, field reflect.StructField, value reflect.Value, expected bool) {
	tag := field.Tag.Get("validate")

	if tag == "" {
		return
	}

	keys := strings.Split(tag, ",")
	for _, key := range keys {
		switch key {
		case "uuid":
			valueStr := value.String()
			if expected || UUID_REGEX.Match([]byte(valueStr)) {
				value.SetString("a valid UUID")
			}
		}
	}
}

func formatTimeField(t *testing.T, field reflect.StructField, value reflect.Value, expected bool) {
	tag := field.Tag.Get("validate")

	if tag == "" {
		return
	}

	keys := strings.Split(tag, ",")
	for _, key := range keys {
		switch key {
		case "ignore":
			value.Set(reflect.ValueOf(time.Time{}))
		}
	}
}

func Equal[T any](t *testing.T, expected T, actual T) bool {
	ptrTyp := reflect.TypeOf(expected)
	if !assert.Equal(t, ptrTyp.Kind(), reflect.Pointer, "expecting a pointer to a struct") {
		return false
	}

	typ := ptrTyp.Elem()
	if !assert.Equal(t, typ.Kind(), reflect.Struct, "expecting a pointer to a struct") {
		return false
	}

	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()

	formatStruct(t, expectedVal, true)
	formatStruct(t, actualVal, false)

	return assert.Equal(t, expected, actual)
}
