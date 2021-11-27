// Main and only package of DataQ
// dataq.go defines the Surfer object and its methods
package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

const (
	Default_sep = "."
)

type Surfer struct {
	sep string
}

type SurferOption func(*Surfer)

func WithSep(sep string) SurferOption {
	return func(s *Surfer) {
		s.sep = sep
	}
}

// return the given field from an interface{}
func (s Surfer) Get(name string, source interface{}) (interface{}, error) {
	f, err := getValueOf(name, source, s.sep)
	if err != nil {
		return nil, err
	}
	return f.Interface(), nil
}

// return the given field as Float64 from an interface{}
func (s Surfer) GetFloat64(name string, source interface{}) (float64, error) {
	value, err := s.Get(name, source)
	if err != nil {
		return 0.0, err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.Float64, reflect.Float32, reflect.Int64, reflect.Int:
		return value.(float64), nil
	case reflect.String:
		return strconv.ParseFloat(value.(string), 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not float64 but %v", name, t)
	}
}

// return the given field as Int64 from an interface{}
func (s Surfer) GetInt64(name string, source interface{}) (int64, error) {
	value, err := s.Get(name, source)
	if err != nil {
		return 0.0, err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.Int64, reflect.Int:
		return value.(int64), nil
	case reflect.String:
		return strconv.ParseInt(value.(string), 0, 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not int64 but %v", name, t)
	}
}

// return the given field as String from an interface{}
func (s Surfer) GetString(name string, source interface{}) (string, error) {
	value, err := s.Get(name, source)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

// return the given field as Bool from an interface{}
func (s Surfer) GetBool(name string, source interface{}) (bool, error) {
	value, err := s.Get(name, source)
	if err != nil {
		return false, err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.Bool:
		return value.(bool), nil
	default:
		return false, fmt.Errorf("variable %v is not bool but %v", name, t)
	}
}

// update the given field of the interface{} with a string
func (s Surfer) SetString(name string, value string, source interface{}) error {
	f, err := getValueOf(name, source, s.sep)
	if err != nil {
		return err
	}
	if f.IsValid() {
		if f.CanSet() {
			f.SetString(value)
			return nil
		}
		return fmt.Errorf("field %v cannot be updated", f)
	}
	return fmt.Errorf("field %v not valid for changing", f)
}

// update the given field of the interface{} with a Int64
func (s Surfer) SetInt64(name string, value int64, source interface{}) error {
	v, err := getValueOf(name, source, s.sep)
	if err != nil {
		return err
	}
	v.SetInt(value)
	return nil
}

// update the given field of the interface{} with a Float64
func (s Surfer) SetFloat64(name string, value float64, source interface{}) error {
	v, err := getValueOf(name, source, s.sep)
	if err != nil {
		return err
	}
	v.SetFloat(value)
	return nil
}

// update the given field of the interface{} with a Bool
func (s Surfer) SetBool(name string, value bool, source interface{}) error {
	v, err := getValueOf(name, source, s.sep)
	if err != nil {
		return err
	}
	v.SetBool(value)
	return nil
}

// return the list of exported fields using their fully qualified names from the interface{}
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
						fields = append(fields, f_name + "." + name)
					}
				case reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Float32:
					fields = append(fields, f_name)
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

// return all fields and their values as a map from the interface{}
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
					// the field is a struct ready for a sublevel search
					subdata, err := s.GetFlatData(f.Interface())
					if err != nil {
						return data, err
					}
					for name, value := range subdata {
						data[f_name + s.sep + name] = value
					}
				case reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Float32:
					data[f_name] = f.Interface()
				default:
					log.Printf("field %v got a not supported type %v", f_name, f.Kind())
				}
			} else {
				log.Printf("field %v is not valid", f_name)
			}
		}
		return data, nil
	default:
		return data, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

func NewSurfer(opts ...SurferOption) *Surfer {
	return &Surfer{
		sep: Default_sep,
	}
}