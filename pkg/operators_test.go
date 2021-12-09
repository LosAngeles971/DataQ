package pkg

import (
	"encoding/json"
	"testing"
)

const (
	JJ = `
{
	"key1": 1.0,
	"key2": true,
	"key3": 10,
	"key4": "ciao"
}
`
)

var JJ_translate = map[string]interface{}{
	"key1": float64(1.0),
	"key2": true,
	"key3": float64(10), // unmarshal of number  convert int -> float64 (?)
	"key4": "ciao",
}

type JData struct {
	Data map[string]interface{}
}

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
	_, _, err := get(Beta_name, l1, s.sep)
	if err == nil {
		t.Error("beta cannot be accessed")
	}
}

func TestReadL2NotExistentField(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	_, _, err2 := get("Alfa.Omega", l1, s.sep)
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

func TestReadValidFieldFromMap(t *testing.T) {
	l1 := getData()
	s := NewSurfer()
	z1, err := s.GetFloat64(Zeta+"."+Zeta_field1, l1)
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

func TestReadFromJson(t *testing.T) {
	dd := JData{
		Data: map[string]interface{}{},
	}
	err := json.Unmarshal([]byte(JJ), &dd.Data)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range JJ_translate {
		vv, ok := dd.Data[k]
		if !ok {
			t.Fatalf("missing key %v", k)
		}
		eq, err := Compare(vv, v)
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Errorf("key %v must be %v not %v", k, v, vv)
		}
	}
	s := NewSurfer()
	mm, err := s.GetFlatData(dd)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range JJ_translate {
		vv, ok := mm["Data."+k]
		if !ok {
			t.Fatalf("missing key %v", k)
		}
		eq, err := Compare(vv, v)
		if err != nil {
			t.Fatal(err)
		}
		if !eq {
			t.Errorf("key %v must be %v not %v", k, v, vv)
		}
	}
}
