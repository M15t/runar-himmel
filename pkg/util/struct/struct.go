package structutil

import (
	"reflect"

	"github.com/imdatngo/mergo"
)

// ToMap converts a struct into a map, respect the json tag name
func ToMap[T any](in T) map[string]interface{} {
	out := make(map[string]interface{})
	mergo.Map(&out, in, mergo.WithJSONTagLookup)
	return out
}

// Convert converts requests struct to repo struct
func Convert[U any, T any](requestStruct *U, targetStruct *T) {
	requestStructType := reflect.TypeOf(*requestStruct)
	requestStructValue := reflect.ValueOf(*requestStruct)

	for i := 0; i < requestStructType.NumField(); i++ {
		field := requestStructType.Field(i)
		value := requestStructValue.Field(i).Interface()

		if field.Type.Kind() == reflect.Ptr && !reflect.ValueOf(value).IsNil() {
			targetField := reflect.ValueOf(targetStruct).Elem().FieldByName(field.Name)
			if targetField.IsValid() { // Ensure target field exists
				targetField.Set(reflect.ValueOf(value).Elem())
			}
		}
	}
}
