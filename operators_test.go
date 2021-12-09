package main

import (
	"testing"
)

func TestCompare(t *testing.T) {
	aa := map[string]interface{}{
		"a": float32(1.0),
		"b": float64(5.0),
		"c": "ciao",
		"d": int(5),
		"e": int64(12),
		"f": true,
	}
	bb := map[string]interface{}{
		"a": float32(1.0),
		"b": float64(5.0),
		"c": "ciao",
		"d": int(5),
		"e": int64(12),
		"f": true,
	}
	for k := range aa {
		v1, ok := aa[k]
		if !ok {
			t.Fatalf("missing %v from aa", k)
		}
		v2, ok := bb[k]
		if !ok {
			t.Fatalf("missing %v from bb", k)
		}
		ok, err := Compare(v1, v2)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("different values %v and %v", v1, v2)
		}
	}
}

func TestReadL1ExistentField(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	alfa, err := s.GetFloat64(Alfa_name, l1)
	if err != nil {
		t.Fatal(err)
	}
	if alfa != l1.Alfa {
		t.Errorf("Alfa should be %v not %v", l1.Alfa, alfa)
	}
}
func TestReadL1UnexportedField(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	_, _, err := s.get(Beta_name, l1)
	if err == nil {
		t.Error("beta cannot be accessed")
	}
}

func TestReadL2NotExistentField(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	_, _, err2 := s.get("Alfa.Omega", l1)
	if err2 == nil {
		t.Errorf("varibale Alfa.omega should not exist")
	}
}

func TestReadL2ExistentField(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	omega, err2 := s.GetString("Gamma.Omega", l1)
	if err2 != nil {
		t.Fatal(err2)
	}
	if omega != l1.Gamma.Omega {
		t.Errorf("Omega should be %v not %v", l1.Gamma.Omega, omega)
	}
}

/*
func TestUpdateL1Field(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	err := s.SetFloat64(Alfa_name, Alfa_update, &l1)
	if err != nil {
		t.Fatal(err)
	}
	alfa, err := s.GetFloat64(Alfa_name, l1)
	if err != nil {
		t.Fatal(err)
	}
	if alfa != Alfa_update {
		t.Errorf("%v should be %v not %v", Alfa_name, Alfa_update, alfa)
	}
}

func TestUpdateL2Field(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	err2 := s.SetString("Gamma.Omega", Omega_update, &l1)
	if err2 != nil {
		t.Fatal(err2)
	}
	omega, err := s.GetString("Gamma.Omega", l1)
	if err != nil {
		t.Fatal(err)
	}
	if omega != Omega_update {
		t.Errorf("%v should be %v not %v", Omega_name, Omega_update, omega)
	}
}
*/

func TestReadValidFieldFromMap(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	z1, err := s.GetFloat64(Zeta + "." + Zeta_field1, l1)
	if err != nil {
		t.Fatal(err)
	}
	if z1 != l1.Zeta[Zeta_field1] {
		t.Errorf("field %v must be %v not %v", Zeta_field1, l1.Zeta[Zeta_field1], z1)
	}
	ok, err := Compare(z1, l1.Zeta[Zeta_field1])
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("field %v must be %v not %v", Zeta_field1, l1.Zeta[Zeta_field1], z1)
	}
}