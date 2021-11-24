package main

import (
	"testing"
)

type Level2 struct {
	Ypsilon int
	Omega   string
}

type Level1 struct {
	Alfa  float64
	beta  string
	Gamma *Level2
}

const (
	Alfa_name = "Alfa"
	Alfa_value = 1.0
	Alfa_update = 2.0
	Ypsilon_name = "Ypsilon"
	Ypsilon_value = 10
	Omega_name = "Omega"
	Omega_value = "test2"
	Omega_update = "updated"
	Beta_name = "beta"
	Beta_value = "test1"
)

var Vars = map[string]interface{}{
	Alfa_name: Alfa_value,
	"Gamma." + Ypsilon_name: Ypsilon_value,
	"Gamma." + Omega_name: Omega_value,
}

func getData() Level1 {
	l2 := Level2{
		Ypsilon: Ypsilon_value,
		Omega: Omega_value,
	}
	
	return Level1{
		Alfa: Alfa_value,
		beta: Beta_value,
		Gamma: &l2,
	}
}
func TestReadL1ExistentField(t *testing.T) {
	l1 := getData()
	alfa, err := GetFloat64(Alfa_name, l1)
	if err != nil {
		t.Fatal(err)
	}
	if alfa != l1.Alfa {
		t.Errorf("Alfa should be %v not %v", l1.Alfa, alfa)
	}
}
func TestReadL1UnexportedField(t *testing.T) {
	l1 := getData()
	_, err := Get("beta", l1)
	if err == nil {
		t.Error("beta cannot be accessed")
	}
}

func TestReadL2NotExistentField(t *testing.T) {
	l1 := getData()
	_, err2 := Get("Alfa.Omega", l1)
	if err2 == nil {
		t.Errorf("varibale Alfa.omega should not exist")
	}
}

func TestReadL2ExistentField(t *testing.T) {
	l1 := getData()
	omega, err2 := GetString("Gamma.Omega", l1)
	if err2 != nil {
		t.Fatal(err2)
	}
	if omega != l1.Gamma.Omega {
		t.Errorf("Omega should be %v not %v", l1.Gamma.Omega, omega)
	}
}

func TestUpdateL1Field(t *testing.T) {
	l1 := getData()
	err := SetFloat64(Alfa_name, Alfa_update, &l1)
	if err != nil {
		t.Fatal(err)
	}
	alfa, err := GetFloat64(Alfa_name, l1)
	if err != nil {
		t.Fatal(err)
	}
	if alfa != Alfa_update {
		t.Errorf("%v should be %v not %v", Alfa_name, Alfa_update, alfa)
	}
}

func TestUpdateL2Field(t *testing.T) {
	l1 := getData()
	err2 := SetString("Gamma.Omega", Omega_update, &l1)
	if err2 != nil {
		t.Fatal(err2)
	}
	omega, err := GetString("Gamma.Omega", l1)
	if err != nil {
		t.Fatal(err)
	}
	if omega != Omega_update {
		t.Errorf("%v should be %v not %v", Omega_name, Omega_update, omega)
	}
}

func TestGetVars(t *testing.T) {
	l1 := getData()
	vars, err := GetVars(l1)
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
	l1 := getData()
	vars, err := GetFlatData(l1)
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
