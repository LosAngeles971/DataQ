// Main and only package of DataQ
// dataq.go defines the Surfer object and its methods
package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// GetBool returns the float64 value of the given field
func (s Surfer) GetFloat64(name string, source interface{}) (float64, error) {
	i, t, err := get(name, source, s.sep)
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
	i, t, err := get(name, source, s.sep)
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
	i, t, err := get(name, source, s.sep)
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
	i, t, err := get(name, source, s.sep)
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

// Two fields comparison without first knowing their types
func Compare(f1 interface{}, f2 interface{}) (bool, error) {
	k1 := datatype(f1)
	k2 := datatype(f2)
	if k1 != k2 {
		log.Tracef("different types %v and %v for %v and %v", k1, k2, f1, f2)
		return false, nil
	}
	switch k1 {
	case T_INT:
		if f1.(int) != f2.(int) {
			return false, nil
		}
	case T_INT64:
		if f1.(int64) != f2.(int64) {
			return false, nil
		}
	case T_FLOAT32:
		if f1.(float32) != f2.(float32) {
			return false, nil
		}
	case T_FLOAT64:
		if f1.(float64) != f2.(float64) {
			return false, nil
		}
	case T_BOOL:
		if f1.(bool) != f2.(bool) {
			return false, nil
		}
	case T_STRING:
		if f1.(string) != f2.(string) {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unsupported type %v", k1)
	}
	return true, nil
}
