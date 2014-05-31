// expression_test.go
package engine

import (
	"testing"
)

func TestGetNumber(t *testing.T) {
	type datacase struct {
		e      Expression
		m      map[string]interface{}
		result float64
	}
	var data = []datacase{
		{Expression{Reference: "a"}, map[string]interface{}{"a": 12.34}, 12.34},
		{Expression{Reference: "a"}, map[string]interface{}{"a": "-12.32"}, -12.32},
		{Expression{Reference: "a.b"}, map[string]interface{}{"a.b": "1"}, 1},
		{Expression{Reference: "a.b.c"}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": 19}}}, 19},
		{Expression{Reference: "a", Number: 32.1}, map[string]interface{}{"a": 12.34}, 12.34},
		{Expression{Number: -12.32}, map[string]interface{}{"a": "5"}, -12.32},
		{Expression{Number: 4}, nil, 4},
		{Expression{Text: "21", Number: 17}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": 19}}}, 17},
	}

	for _, d := range data {
		n := Notif{Data: d.m}
		f, err := d.e.getNumber(&n)
		if err != nil {
			t.Fatalf("%v %+v %+v", err, d.e, n)
		}
		if f != d.result {
			t.Fatalf("expression should be %f instead of %f, %+v", d.result, f, d)
		}
	}
}
func TestGetString(t *testing.T) {
	type datacase struct {
		e      Expression
		m      map[string]interface{}
		result string
	}
	var data = []datacase{
		{Expression{Reference: "a"}, map[string]interface{}{"a": "12.34"}, "12.34"},
		{Expression{Reference: "a"}, map[string]interface{}{"a": "zz zz.x"}, "zz zz.x"},
		{Expression{Reference: "a.b"}, map[string]interface{}{"a.b": "1"}, "1"},
		{Expression{Reference: "a.b.c"}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": "hawaii"}}}, "hawaii"},
		{Expression{Reference: "a", Text: "zx"}, map[string]interface{}{"a": "rrr"}, "rrr"},
		{Expression{Text: "qwe"}, map[string]interface{}{"a": "asd"}, "qwe"},
		{Expression{Text: "four"}, nil, "four"},
		{Expression{Text: "downton abbey", Number: 17}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": "the shield"}}}, "downton abbey"},
	}

	for _, d := range data {
		n := Notif{Data: d.m}
		s, err := d.e.getString(&n)
		if err != nil {
			t.Fatalf("%v %+v %+v", err, d.e, n)
		}
		if s != d.result {
			t.Fatalf("expression should be %q instead of %q, %+v", d.result, s, d)
		}
	}
}

func TestToRuleJSON(t *testing.T) {
	type datacase struct {
		e      Expression
		result interface{}
	}
	var data = []datacase{
		{e: Expression{Number: 8.123}, result: 8.123},
		{e: Expression{Number: -32.234}, result: -32.234},
		{e: Expression{Number: 0.0}, result: 0.0},
		{e: Expression{Text: "x-32.234a"}, result: "x-32.234a"},
		{e: Expression{Text: "a.b.c.d"}, result: "a.b.c.d"},
		{e: Expression{Reference: "a"}, result: "$a"},
		{e: Expression{Reference: "a.b.c.d"}, result: "$a.b.c.d"},
	}
	for _, d := range data {
		if res := d.e.ToRuleJSON(); res != d.result {
			t.Fatalf("%v should be %v %+v", res, d.result, d)
		}
	}
}

func TestMakeExpressionFromJSON(t *testing.T) {
	type datacase struct {
		i      interface{}
		isNum  bool
		result Expression
	}
	var data = []datacase{
		{result: Expression{Number: 8.123}, i: 8.123, isNum: true},
		{result: Expression{Number: -32.234}, i: -32.234, isNum: true},
		{result: Expression{Number: 0.0}, i: 0.0, isNum: true},
		{result: Expression{Number: 8.0}, i: 8, isNum: true},
		{result: Expression{Number: -32.234}, i: -32.234, isNum: true},
		{result: Expression{Number: 0.0}, i: 0.0, isNum: true},
		{result: Expression{Text: "x-32.234a"}, i: "x-32.234a", isNum: false},
		{result: Expression{Text: "a.b.c.d"}, i: "a.b.c.d", isNum: false},
		{result: Expression{Reference: "a"}, i: "$a", isNum: false},
		{result: Expression{Reference: "a.b.c.d"}, i: "$a.b.c.d", isNum: true},
	}
	for _, d := range data {
		res, err := makeExpressionFromJSON(d.i, d.isNum)
		if err != nil {
			t.Fatalf("%v, %+v", err, d)
		}
		if res != d.result {
			t.Fatalf("%v should be %v, %+v", res, d)
		}
	}
}
