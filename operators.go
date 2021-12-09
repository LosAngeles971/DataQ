// Main and only package of DataQ
// dataq.go defines the Surfer object and its methods
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

// GetBool returns the float64 value of the given field
func (s Surfer) GetFloat64(name string, source interface{}) (float64, error) {
	i, t, err := s.get(name, source)
	if err != nil {
		return 0.0, err
	}
	switch t {
	case T_FLOAT64, T_FLOAT32, T_INT64, T_INT:
		return i.(float64), nil
	case T_STRING:
		return strconv.ParseFloat(i.(string), 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not float64 but %v", name, t)
	}
}

// GetInt64 returns the int64 value of the given field
func (s Surfer) GetInt64(name string, source interface{}) (int64, error) {
	i, t, err := s.get(name, source)
	if err != nil {
		return 0.0, err
	}
	switch t {
	case T_INT64, T_INT:
		return i.(int64), nil
	case T_STRING:
		return strconv.ParseInt(i.(string), 0, 64)
	default:
		return 0.0, fmt.Errorf("variable %v is not int64 but %v", name, t)
	}
}

// GetString returns the string value of the given field
func (s Surfer) GetString(name string, source interface{}) (string, error) {
	i, t, err := s.get(name, source)
	if err != nil {
		return "", err
	}
	switch t {
	case T_PTR, T_STRUCT, T_MAP:
		return "", fmt.Errorf("not supported type for string: %v", t)
	default:
		return i.(string), nil
	}
}

// GetBool returns the bool value of the given field
func (s Surfer) GetBool(name string, source interface{}) (bool, error) {
	i, t, err := s.get(name, source)
	if err != nil {
		return false, err
	}
	switch t {
	case T_BOOL:
		return i.(bool), nil
	case T_STRING:
		if strings.ToUpper(i.(string)) == "TRUE" {
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("variable %v is not bool but %v", name, t)
	}
}

// SetString updates the given string field
func (s Surfer) SetString(name string, value string, source interface{}) error {
	return s.set(name, value, source)
}

// SetInt64 updates the given int64 field
func (s Surfer) SetInt64(name string, value int64, source interface{}) error {
	return s.set(name, value, source)
}

// SetFloat64 updates the given float64 field
func (s Surfer) SetFloat64(name string, value float64, source interface{}) error {
	return s.set(name, value, source)
}

// SetBool updates the given bool field
func (s Surfer) SetBool(name string, value bool, source interface{}) error {
	return s.set(name, value, source)
}

// Two fields comparison without first knowing their types
func Compare(f1 interface{}, f2 interface{}) (bool, error) {
	k1 := reflect.ValueOf(f1).Kind()
	k2 := reflect.ValueOf(f2).Kind()
	if k1 != k2 {
		log.Tracef("different types %v and %v for %v and %v", k1, k2, f1, f2)
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