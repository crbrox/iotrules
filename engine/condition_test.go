package engine

import (
	"testing"
)

func TestNumericCondTrue(t *testing.T) {
	type datacase struct {
		o  Op
		e1 Expression
		e2 Expression
		m  map[string]interface{}
	}
	var data = []datacase{
		{EQ, Expression{Reference: "a"}, Expression{Number: 1}, map[string]interface{}{"a": 1}},
		{EQ, Expression{Reference: "a"}, Expression{Number: 1}, map[string]interface{}{"a": "1"}},
		{EQ, Expression{Reference: "a.b"}, Expression{Number: 10}, map[string]interface{}{"a.b": "10"}},
		{EQ, Expression{Reference: "a.b.c"}, Expression{Number: 99}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": 99}}}},
		{LT, Expression{Reference: "a"}, Expression{Number: 1.1}, map[string]interface{}{"a": 1}},
		{LT, Expression{Reference: "a"}, Expression{Number: 22}, map[string]interface{}{"a": "21.5"}},
		{LT, Expression{Reference: "a.b"}, Expression{Number: 10}, map[string]interface{}{"a.b": "-3"}},
		{LT, Expression{Reference: "a.b.c"}, Expression{Number: 99}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": 0}}}},
		{GT, Expression{Reference: "a"}, Expression{Number: 1.1}, map[string]interface{}{"a": 1.2}},
		{GT, Expression{Reference: "a"}, Expression{Number: 22}, map[string]interface{}{"a": "22.5"}},
		{GT, Expression{Reference: "a.b"}, Expression{Number: -10}, map[string]interface{}{"a.b": "-9"}},
		{GT, Expression{Reference: "a.b.c"}, Expression{Number: 99}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": 100}}}},
	}

	for _, d := range data {
		c := Condition{Op: d.o, Exp1: d.e1, Exp2: d.e2, IsNumber: true}
		n := Notif{Data: d.m}
		m, e := c.Matched(&n)
		if e != nil {
			t.Fatalf("%v %+v %+v", e, c, n)
		}
		if !m {
			t.Fatalf("condition should be true %+v %+v", c, n)
		}
	}

}
func TestStringCondTrue(t *testing.T) {
	type datacase struct {
		o  Op
		e1 Expression
		e2 Expression
		m  map[string]interface{}
	}
	var data = []datacase{
		{EQ, Expression{Reference: "a"}, Expression{Text: "jamarkanda"}, map[string]interface{}{"a": "jamarkanda"}},
		{EQ, Expression{Reference: "a"}, Expression{Text: "1"}, map[string]interface{}{"a": "1"}},
		{EQ, Expression{Reference: "a.b"}, Expression{Text: "c.d"}, map[string]interface{}{"a.b": "c.d"}},
		{EQ, Expression{Reference: "a.b.c"}, Expression{Text: "z"}, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": "z"}}}},
	}

	for _, d := range data {
		c := Condition{Op: d.o, Exp1: d.e1, Exp2: d.e2, IsNumber: false}
		n := Notif{Data: d.m}
		m, e := c.Matched(&n)
		if e != nil {
			t.Fatalf("%v %+v %+v", e, c, n)
		}
		if !m {
			t.Fatalf("condition should be true %+v %+v", c, n)
		}
	}

}
