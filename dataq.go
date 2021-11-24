/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Author: LosAngeles971
Date: November 2021
File: dataq.go
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// check if the field name is valid
func checkFieldName(name string) error {
	if len(name) < 0 {
		return fmt.Errorf("field name cannot be null")
	}
	if !unicode.IsUpper([]rune(name)[0]) {
		return fmt.Errorf("field name must be exported (first letter capitalized) [%v]", name[0])
	}
	return nil
}

// recursive function to browse a struct till to the desired target field
func getValueOf(name string, source interface{}) (reflect.Value, error) {
	fields := strings.Split(name, ".")
	field_name := fields[0]
	if err := checkFieldName(field_name); err != nil {
		return reflect.Value{}, err
	}
	var obj reflect.Value
	if reflect.ValueOf(source).Kind() == reflect.Ptr {
		// this is the case of passing a pointer to a struct because you wanna update a field
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
				// reached the desired field
				return f, nil
			} else {
				switch f.Kind() {
				case reflect.Struct:
					// the field is a struct ready for a sublevel search
					return getValueOf(strings.Join(fields[1:], "."), f.Interface())
				case reflect.Ptr:
					// the field is a pointer to a struct (it should be a struct) ready for a sublevel search
					return getValueOf(strings.Join(fields[1:], "."), f.Interface())
				default:
					// the field is not a struct, a sublevel search is not possible
					return reflect.Value{}, fmt.Errorf("field %v is a not supported type", f)
				}
			}				
		}
		return reflect.Value{}, fmt.Errorf("missing field %v", field_name)
	default:
		return reflect.Value{}, fmt.Errorf("unhandled type of data %v", obj.Kind())
	}
}

// return the asked fully qualified field from source in form of interface{}
func Get(name string, source interface{}) (interface{}, error) {
	f, err := getValueOf(name, source)
	if err != nil {
		return nil, err
	}
	return f.Interface(), nil
}

// return the asked fully qualified field from source in form of float64
func GetFloat64(name string, source interface{}) (float64, error) {
	value, err := Get(name, source)
	if err != nil {
		return 0.0, err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.Float64:
		return value.(float64), nil
	case reflect.String:
		return strconv.ParseFloat(value.(string), 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not float64 but %v", name, t)
	}
}

// return the asked fully qualified field from source in form of int64
func GetInt64(name string, source interface{}) (int64, error) {
	value, err := Get(name, source)
	if err != nil {
		return 0.0, err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.Int:
		return value.(int64), nil
	case reflect.String:
		return strconv.ParseInt(value.(string), 0, 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not int64 but %v", name, t)
	}
}

// return the asked fully qualified field from source in form of string
func GetString(name string, source interface{}) (string, error) {
	value, err := Get(name, source)
	if err != nil {
		return "", err
	}
	t := reflect.TypeOf(value).Kind()
	switch t {
	case reflect.String:
		return value.(string), nil
	default:
		return "", fmt.Errorf("variable %v is not string but %v", name, t)
	}
}

// update the asked fully qualified field from source in form of string
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

// update the asked fully qualified field from source in form of int64
func SetInt64(name string, value int64, source interface{}) error {
	v, err := getValueOf(name, source)
	if err != nil {
		return err
	}
	v.SetInt(value)
	return nil
}

// update the asked fully qualified field from source in form of float64
func SetFloat64(name string, value float64, source interface{}) error {
	v, err := getValueOf(name, source)
	if err != nil {
		return err
	}
	v.SetFloat(value)
	return nil
}