/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Author: LosAngeles971
Date: November 2021
File: helpers.go
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package main

import (
	"fmt"
	"reflect"
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
func getValueOf(name string, source interface{}, sep string) (reflect.Value, error) {
	fields := strings.Split(name, sep)
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
					return getValueOf(strings.Join(fields[1:], "."), f.Interface(), sep)
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