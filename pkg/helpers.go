// helpers.go includes all utility (not exported) functions
package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"unicode"
)

// custom standardization for supported data types
const (
	T_PTR           = 0
	T_STRUCT        = 1
	T_INT           = 2
	T_INT64         = 3
	T_FLOAT32       = 4
	T_FLOAT64       = 5
	T_STRING        = 6
	T_BOOL          = 7
	T_MAP           = 8
	T_NOT_SUPPORTED = -1
)

// datatype returns the type of an interface using a custom standardization
func datatype(i interface{}) int {
	tt := reflect.ValueOf(i).Kind()
	switch tt {
	case reflect.Ptr:
		return T_PTR
	case reflect.Struct:
		return T_STRUCT
	case reflect.Int:
		return T_INT
	case reflect.Int64:
		return T_INT64
	case reflect.Float32:
		return T_FLOAT32
	case reflect.Float64:
		return T_FLOAT64
	case reflect.String:
		return T_STRING
	case reflect.Bool:
		return T_BOOL
	case reflect.Map:
		return T_MAP
	default:
		return T_NOT_SUPPORTED
	}
}

// checkFieldName checks if a relative field's name is syntattically valid
func checkFieldName(name string) bool {
	if len(name) < 0 {
		log.Errorf("field's name cannot be null")
		return false
	}
	if !unicode.IsUpper([]rune(name)[0]) {
		log.Errorf("field must be exported (first letter of the name capitalized) [%v]", name[0])
		return false
	}
	return true
}

// getFieldsFromMap returns the list of keys from a map
func getFieldsFromMap(m interface{}) []string {
	fields := []string{}
	tt := datatype(m)
	if tt != T_MAP {
		log.Errorf("skipped fields recognizing because input is not a map but %v", tt)
		return fields
	}
	kk := reflect.TypeOf(m).Key().Kind()
	if kk != reflect.String {
		log.Errorf("skipped fields recognizing because map's keys are not string but %v", kk)
		return fields
	}
	vv := reflect.TypeOf(m).Elem().Kind()
	switch vv {
	case reflect.Int:
		for k := range m.(map[string]int) {
			fields = append(fields, k)
		}
	case reflect.Int64:
		for k := range m.(map[string]int64) {
			fields = append(fields, k)
		}
	case reflect.Float32:
		for k := range m.(map[string]float32) {
			fields = append(fields, k)
		}
	case reflect.Float64:
		for k := range m.(map[string]float64) {
			fields = append(fields, k)
		}
	case reflect.Bool:
		for k := range m.(map[string]bool) {
			fields = append(fields, k)
		}
	case reflect.String:
		for k := range m.(map[string]string) {
			fields = append(fields, k)
		}
	case reflect.Interface:
		for k := range m.(map[string]interface{}) {
			fields = append(fields, k)
		}
	default:
		log.Errorf("skipped fields recognizing because the type of map's values is unsupported: %v", vv)
	}
	return fields
}

// getValueFromMap returs the value associated to the key "field" from a given map in the form of interface{}
func getValueFromMap(field string, i interface{}) (interface{}, error) {
	tt := datatype(i)
	if tt != T_MAP {
		return nil, fmt.Errorf("skipped fields recognizing because input is not a map but code: %v", tt)
	}
	m := reflect.ValueOf(i)
	// in case, to get type of fields -> datatype(reflect.TypeOf(i).Elem())
	for _, e := range m.MapKeys() {
		if e.Interface().(string) == field {
			return m.MapIndex(e).Interface(), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("map does not contain field %v", field)
}

// getValueOf returns the value of a given variable, recursively browsing the given data in the form of an interface{}
func getValueOf(name string, source interface{}, sep string) (interface{}, error) {
	fields := strings.Split(name, sep)
	field_name := fields[0]
	if !checkFieldName(field_name) {
		return nil, fmt.Errorf("field %v is not valid", field_name)
	}
	var obj reflect.Value
	if reflect.ValueOf(source).Kind() == reflect.Ptr {
		// taking the object from the pointer
		obj = reflect.ValueOf(source).Elem()
	} else {
		obj = reflect.ValueOf(source)
	}
	switch obj.Kind() {
	case reflect.Struct:
		f_value := obj.FieldByName(field_name)
		// f must not be a (struct) zero value
		if f_value.IsValid() {
			switch f_value.Kind() {
			case reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Float32:
				if len(fields) == 1 {
					// positive exit: reached the target field
					return f_value.Interface(), nil
				} else {
					return nil, fmt.Errorf("field [%v] is primitive, cannot be a sublevel ", field_name)
				}
			case reflect.Struct, reflect.Ptr:
				if len(fields) == 1 {
					// positive exit: reached the target field
					return nil, fmt.Errorf("requested field [%v] points to a struct or a pointer", field_name)
				} else {
					if f_value.IsNil() {
						return nil, fmt.Errorf("surfing stopped by nil field [%v]", field_name)
					} else {
						// going to the sublevel (struct) or getting the object from the pointer
						return getValueOf(strings.Join(fields[1:], sep), f_value.Interface(), sep)
					}
				}
			case reflect.Map:
				if len(fields) == 1 {
					// positive exit: reached the target field
					return nil, fmt.Errorf("requested field [%v] points to a map", field_name)
				} else {
					if f_value.IsNil() {
						return nil, fmt.Errorf("surfing stopped by nil field [%v]", field_name)
					} else {
						// only one level of mapping is supported
						// map of complex objects is not supported
						value, err := getValueFromMap(strings.Join(fields[1:], sep), f_value.Interface())
						if err != nil {
							return nil, err
						}
						return value, err
					}
				}
			default:
				// error: field is not a struct or pointer (deep dive not possible)
				return nil, fmt.Errorf("field %v is a not supported type", field_name)
			}
		}
		return nil, fmt.Errorf("missing field %v", field_name)
	default:
		return nil, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

// get returns the value of the given field from the given data in the form of an interface{}
func get(name string, source interface{}, sep string) (interface{}, int, error) {
	f, err := getValueOf(name, source, sep)
	if err != nil {
		return nil, T_NOT_SUPPORTED, err
	}
	t := datatype(f)
	if t == T_NOT_SUPPORTED {
		return f, T_NOT_SUPPORTED, fmt.Errorf("type of data not supported: %v", t)
	}
	return f, t, nil
}