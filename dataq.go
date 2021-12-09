// Main and only package of DataQ
// dataq.go defines the Surfer object and its methods
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
)

const (
	Default_sep = "."
)

type Surfer struct {
	sep string
}

type SurferOption func(*Surfer)

// WithSep sets the separation string for the fully qualified name of the fields
func WithSep(sep string) SurferOption {
	return func(s *Surfer) {
		s.sep = sep
	}
}

// GetVars extracts from the source a list of all exportable fields using their fully qualified names
func (s Surfer) GetVars(source interface{}) ([]string, error) {
	fields := []string{}
	var obj reflect.Value
	if reflect.ValueOf(source).Kind() == reflect.Ptr {
		// this is the case of passing a pointer to a struct because you wanna update a field
		obj = reflect.ValueOf(source).Elem()
	} else {
		obj = reflect.ValueOf(source)
	}
	switch obj.Kind() {
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			f := obj.Field(i)
			f_name := obj.Type().Field(i).Name
			err := checkFieldName(f_name)
			// f must not be a (struct) zero value
			if f.IsValid() && err == nil {
				switch f.Kind() {
				case reflect.Struct, reflect.Ptr:
					// the field is a struct ready for a sublevel search
					ff, err := s.GetVars(f.Interface())
					if err != nil {
						return fields, err
					}
					for _, name := range ff {
						fields = append(fields, f_name+s.sep+name)
					}
				case reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Float32:
					fields = append(fields, f_name)
				case reflect.Map:
					// only one level of mapping is supported
					// map of complex objects is not supported
					ff := getFieldsFromMap(f.Interface())
					for _, name := range ff {
						fields = append(fields, f_name+s.sep+name)
					}
				default:
					log.Printf("field %v got a not supported type %v", f_name, f.Kind())
				}
			} else {
				log.Printf("field %v is not valid", f_name)
			}
		}
		return fields, nil
	default:
		return fields, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

// GetFlatData returns a map of interface{} including all fields extracted from the source
func (s Surfer) GetFlatData(source interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	var obj reflect.Value
	if reflect.ValueOf(source).Kind() == reflect.Ptr {
		// this is the case of passing a pointer to a struct because you wanna update a field
		obj = reflect.ValueOf(source).Elem()
	} else {
		obj = reflect.ValueOf(source)
	}
	switch obj.Kind() {
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			f := obj.Field(i)
			f_name := obj.Type().Field(i).Name
			err := checkFieldName(f_name)
			// f must not be a (struct) zero value
			if f.IsValid() && err == nil {
				switch f.Kind() {
				case reflect.Struct, reflect.Ptr:
					// the struct's field is a struct ready for a sublevel search
					subdata, err := s.GetFlatData(f.Interface())
					if err != nil {
						return data, err
					}
					for name, value := range subdata {
						data[f_name+s.sep+name] = value
					}
				case reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Float32:
					// the struct's field is a supported primitive data
					data[f_name] = f.Interface()
				case reflect.Map:
					// the struct's field is a map
					for _, k := range getFieldsFromMap(f.Interface()) {
						value, err := getValueFromMap(k, f.Interface())
						if err != nil {
							return data, err
						}
						data[f_name + s.sep + k] = value
					}
				default:
					log.Printf("field %v got a not supported type %v", f_name, f.Kind())
				}
			} else {
				log.Printf("field %v is not valid", f_name)
			}
		}
		return data, nil
	case reflect.Map:
		// starting object is not struct but map
		keys := getFieldsFromMap(obj)
		for _, k := range keys {
			var err error
			data[k], err = getValueFromMap(k, obj)
			if err != nil {
				return data, err
			}
		}
		return data, nil
	default:
		return data, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

// NewSurfer creates a pointer to a new Surfer object with default configuration
func NewSurfer(opts ...SurferOption) *Surfer {
	s := &Surfer{
		sep: Default_sep,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
