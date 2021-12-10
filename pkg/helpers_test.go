package pkg

import (
	"testing"
)

type Level3 struct {
	Delta int
}

type Level2 struct {
	Ypsilon int
	Omega   string
	Epsilon *Level3
}

type Level1 struct {
	Alfa  float64
	beta  string
	Gamma *Level2
	Zeta  map[string]float64
}

const (
	SEP                   = "."
	Alfa_name             = "Alfa"
	Alfa_value            = 1.0
	Alfa_update           = 2.0
	Gamma                 = "Gamma"
	Epsilon  			  = "Epsilon"
	Ypsilon_name          = "Ypsilon"
	Ypsilon_value         = 10
	Omega_name            = "Omega"
	Omega_value           = "test2"
	Omega_update          = "updated"
	Beta_name             = "beta"
	Beta_value            = "test1"
	Expression            = "Alfa + Gamma_Ypsilon"
	Expression_result     = 11.0
	Zeta                  = "Zeta"
	Zeta_field1           = "zeta1"
	Zeta_field1_value     = 1.0
	Zeta_field2           = "zeta2"
	Zeta_field2_value     = 2.0
	Zeta_supported_fields = 2
)

var Vars = map[string]interface{}{
	Alfa_name:                  Alfa_value,
	Gamma + SEP + Ypsilon_name: Ypsilon_value,
	Gamma + SEP + Omega_name:   Omega_value,
	Zeta + SEP + Zeta_field1:   Zeta_field1_value,
	Zeta + SEP + Zeta_field2:   Zeta_field2_value,
}

func getData() Level1 {
	l2 := Level2{
		Ypsilon: Ypsilon_value,
		Omega:   Omega_value,
		Epsilon: nil,
	}
	return Level1{
		Alfa:  Alfa_value,
		beta:  Beta_value,
		Gamma: &l2,
		Zeta: map[string]float64{
			Zeta_field1: Zeta_field1_value,
			Zeta_field2: Zeta_field2_value,
		},
	}
}
func TestDatatype(t *testing.T) {
	aa := map[string]interface{}{
		"a": float32(1.0),
		"b": float64(5.0),
		"c": "ciao",
		"d": int(5),
		"e": int64(12),
		"f": true,
	}
	bb := map[string]int{
		"a": T_FLOAT32,
		"b": T_FLOAT64,
		"c": T_STRING,
		"d": T_INT,
		"e": T_INT64,
		"f": T_BOOL,
	}
	for k := range aa {
		d := datatype(aa[k]) 
		if d != bb[k] {
			t.Errorf("expected type of %v is %v not %v", k, d, bb[k])
		}
	}
	if datatype(aa) != T_MAP {
		t.Errorf("expected type of aa is map not %v", datatype(aa))
	}
	l1 := getData()
	if datatype(l1) != T_STRUCT {
		t.Errorf("wrong type for %v - %v", l1, datatype(l1))
	}
	if datatype(l1.Gamma) != T_PTR {
		t.Errorf("wrong type for %v - %v", l1, datatype(l1.Gamma))
	}
	if datatype(l1.Zeta) != T_MAP {
		t.Errorf("wrong type for %v - %v", l1, datatype(l1.Zeta))
	}
	if datatype(l1.Alfa) != T_FLOAT64 {
		t.Errorf("wrong type for %v - %v", l1, datatype(l1.Alfa))
	}
	if datatype(l1.Gamma.Ypsilon) != T_INT {
		t.Errorf("wrong type for %v - %v", l1, datatype(l1.Gamma.Ypsilon))
	}
}

func TestGetFieldsFromMap(t *testing.T) {
	l1 := getData()
	fields := getFieldsFromMap(l1.Zeta)
	if len(fields) != Zeta_supported_fields {
		t.Errorf("number of Zeta supported fields must be %v not %v", Zeta_supported_fields, len(fields))
	}
	c := 0
	for i := range fields {
		if fields[i] == Zeta_field1 || fields[i] == Zeta_field2 {
			c++
		}
	}
	if c != Zeta_supported_fields {
		t.Errorf("number of Zeta recognized fields must be %v not %v", Zeta_supported_fields, c)
	}
}

func TestGetValueFromMap(t *testing.T) {
	l1 := getData()
	vv, err := getValueFromMap(Zeta_field1, l1.Zeta)
	if err != nil {
		t.Fatal(err)
	}
	v := vv.(float64)
	if v != Zeta_field1_value {
		t.Errorf("value of %v must be %v not %v", Zeta_field1, Zeta_field1_value, vv)
	}
}

func TestGetValueOf(t *testing.T) {
	l1 := getData()
	vv, err := getValueOf("Alfa", l1, SEP)
	if err != nil {
		t.Fatal(err)
	}
	if vv.(float64) != Alfa_value {
		t.Errorf("variable %v must be %v not %v", Alfa_name, Alfa_value, vv)
	}
}

func TestGet(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	vv, tt, err := get("Alfa", l1, s.sep)
	if err != nil {
		t.Fatal(err)
	}
	if tt != T_FLOAT64 {
		t.Errorf("variable %v must be float64 not %v", Alfa_name, tt)
	}
	if vv.(float64) != Alfa_value {
		t.Errorf("variable %v must be %v not %v", Alfa_name, Alfa_value, vv)
	}
}

func TestGetOfNil(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	_, _, err := get(Gamma + s.sep + Epsilon, l1, s.sep)
	if err == nil {
		t.Fatal("accessing a nil data cannot be done")
	}
}