package formatter

import (
	"encoding/json"
	"reflect"
)

// filterFields filters data to only include the specified fields.
// If fields is nil or empty, data is returned unchanged.
func filterFields(data interface{}, fields []string) interface{} {
	if len(fields) == 0 {
		return data
	}

	v := reflect.ValueOf(data)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return data
		}
		v = v.Elem()
	}

	// Handle slices: filter each element
	if v.Kind() == reflect.Slice {
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = filterFields(v.Index(i).Interface(), fields)
		}
		return result
	}

	// For structs and maps, marshal to map then filter keys
	return filterMap(data, fields)
}

// filterMap marshals data to a map and keeps only requested fields
func filterMap(data interface{}, fields []string) interface{} {
	b, err := json.Marshal(data)
	if err != nil {
		return data
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return data
	}

	keep := make(map[string]bool, len(fields))
	for _, f := range fields {
		keep[f] = true
	}

	filtered := make(map[string]interface{}, len(fields))
	for k, v := range m {
		if keep[k] {
			filtered[k] = v
		}
	}
	return filtered
}
