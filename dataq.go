/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Author: LosAngeles971
Date: November 2021
File: dataq.go
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// check if the field name is valid
func checkFieldName(name string) error {
	if len(name) < 0 {
		return fmt.Errorf("field's name cannot be null")
	}
	if !unicode.IsUpper([]rune(name)[0]) {
		return fmt.Errorf("field must be exported (first letter of the name capitalized) [%v]", name[0])
	}
	return nil
}

// recursive browsing of a struct till to the desired target field
func getValueOf(name string, source interface{}) (reflect.Value, error) {
	fields := strings.Split(name, ".")
	field_name := fields[0]
	if err := checkFieldName(field_name); err != nil {
		return reflect.Value{}, err
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
		f := obj.FieldByName(field_name)
		// f must not be a (struct) zero value
		if f.IsValid() {
			if len(fields) == 1 {
				// positive exit: reached the target field
				return f, nil
			} else {
				switch f.Kind() {
				case reflect.Struct, reflect.Ptr:
					// going to the sublevel (struct) or getting the object from the pointer
					return getValueOf(strings.Join(fields[1:], "."), f.Interface())
				default:
					// error: field is not a struct or pointer (deep dive not possible)
					return reflect.Value{}, fmt.Errorf("field %v is a not supported type", f)
				}
			}				
		}
		return reflect.Value{}, fmt.Errorf("missing field %v", field_name)
	default:
		return reflect.Value{}, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

// return the given field from an interface{}
func Get(name string, source interface{}) (interface{}, error) {
	f, err := getValueOf(name, source)
	if err != nil {
		return nil, err
	}
	return f.Interface(), nil
}

// return the given field as Float64 from an interface{}
func GetFloat64(name string, source interface{}) (float64, error) {
	value, err := Get(name, source)
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
func GetInt64(name string, source interface{}) (int64, error) {
	value, err := Get(name, source)
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
func GetString(name string, source interface{}) (string, error) {
	value, err := Get(name, source)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

// return the given field as Bool from an interface{}
func GetBool(name string, source interface{}) (bool, error) {
	value, err := Get(name, source)
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
func SetString(name string, value string, source interface{}) error {
	f, err := getValueOf(name, source)
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
func SetInt64(name string, value int64, source interface{}) error {
	v, err := getValueOf(name, source)
	if err != nil {
		return err
	}
	v.SetInt(value)
	return nil
}

// update the given field of the interface{} with a Float64
func SetFloat64(name string, value float64, source interface{}) error {
	v, err := getValueOf(name, source)
	if err != nil {
		return err
	}
	v.SetFloat(value)
	return nil
}

// update the given field of the interface{} with a Bool
func SetBool(name string, value bool, source interface{}) error {
	v, err := getValueOf(name, source)
	if err != nil {
		return err
	}
	v.SetBool(value)
	return nil
}

// compare two fields without knowing their types
func Compare(f1 interface{}, f2 interface{}) (bool, error) {
	k1 := reflect.TypeOf(f1).Kind()
	k2 := reflect.TypeOf(f2).Kind()
	if k1 != k2 {
		return false, nil
	}
	switch k1 {
	case reflect.Int:
		if f1.(int) != f2.(int) {
			return false, nil
		}
	case reflect.Int64:
		if f1.(int64) != f2.(int64) {
			return false, nil
		}
	case reflect.Float32:
		if f1.(float32) != f2.(float32) {
			return false, nil
		}
	case reflect.Float64:
		if f1.(float64) != f2.(float64) {
			return false, nil
		}
	case reflect.Bool:
		if f1.(bool) != f2.(bool) {
			return false, nil
		}
	case reflect.String:
		if f1.(string) != f2.(string) {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unsupported type %v", k1)
	}
	return true, nil
}

// return the list of exported fields using their fully qualified names from the interface{}
func GetVars(source interface{}) ([]string, error) {
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
					ff, err := GetVars(f.Interface())
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
func GetFlatData(source interface{}) (map[string]interface{}, error) {
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
					subdata, err := GetFlatData(f.Interface())
					if err != nil {
						return data, err
					}
					for name, value := range subdata {
						data[f_name + "." + name] = value
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