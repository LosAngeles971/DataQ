package main

import (
	"reflect"
	"testing"

	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
)

type Level2 struct {
	Ypsilon int
	Omega   string
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

func TestGetVars(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	vars, err := s.GetVars(l1)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range vars {
		_, ok := Vars[name]
		if !ok {
			t.Errorf("variable %v is not expected", name)
		}
	}
	for name := range Vars {
		ok := false
		for i := range vars {
			if vars[i] == name {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("variable %v is missing", name)
		}
	}
}

func TestGetFlatDatas(t *testing.T) {
	ff := float64(5.0)
	//vv := reflect.ValueOf(ff)
	tt := reflect.ValueOf(ff).Kind()
	log.Println(tt)
	log.SetLevel(log.TraceLevel)
	l1 := getData()
	s := NewSurfer()
	vars, err := s.GetFlatData(l1)
	if err != nil {
		t.Fatal(err)
	}
	for name, value := range vars {
		vv, ok := Vars[name]
		if !ok {
			t.Errorf("variable %v is not expected", name)
		}
		ok, err := Compare(vv, value)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("variable %v got %v instead of %v", name, value, vv)
		}
	}
	for name, value := range Vars {
		vv, ok := vars[name]
		if !ok {
			t.Errorf("variable %v is missing", name)
		}
		ok, err := Compare(vv, value)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("variable %v got %v instead of %v", name, vv, value)
		}
	}
}

func TestMath(t *testing.T) {
	l1 := getData()
	s := NewSurfer(WithSep("_"))
	expr, err := govaluate.NewEvaluableExpression(Expression)
	if err != nil {
		t.Fatal(err)
	}
	data, err := s.GetFlatData(l1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := expr.Evaluate(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.(float64) != Expression_result {
		t.Errorf("result should be %v not %v", Expression_result, result)
	}
}
