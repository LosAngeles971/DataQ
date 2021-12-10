package pkg

import (
	"reflect"
	"testing"

	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
)

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
