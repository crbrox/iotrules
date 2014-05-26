package engine

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestSimpleMatch(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var cases = []struct {
		Op     Op
		E1, E2 string
		IsN    bool
	}{
		{EQ, "$uno", "1", true},
		{GT, "$uno", "0", true},
		{LT, "$uno", "3", true},
		{NE, "$uno", "2", true},
	}
	for _, c := range cases {
		var rule = Rule{Op: c.Op, Exp1: c.E1, Exp2: c.E2, IsNumber: c.IsN}
		b, e := rule.Matched(&n)
		if e != nil {
			t.Fatal(e)
		}
		if !b {
			t.Fatalf("%+v should match %+v\n", n, rule)
		}

	}

}
func TestSimpleNotMatch(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var rule = Rule{Op: EQ, Exp1: "$uno", Exp2: "2", IsNumber: true}
	b, e := rule.Matched(&n)
	if e != nil {
		t.Fatal(e)
	}
	if b {
		t.Fatalf("%+v should NOT match %+v\n", n, rule)
	}
}
func TestStringMatch(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var rule = Rule{Op: EQ, Exp1: "$cadena.otro.mas", Exp2: "alfanumérica"}
	b, e := rule.Matched(&n)
	if e != nil {
		t.Fatal(e)
	}
	if !b {
		t.Fatalf("%+v should match %+v\n", n, rule)
	}
}
func TestFromDataSimple(t *testing.T) {
	data := map[string]interface{}{
		"op":       ">",
		"ls":       "$a",
		"rs":       "3",
		"isNumber": true,
		"type":     "SMS",
	}
	r, err := FromMap(data)
	if err != nil {
		t.Fatal(err)
	}

	//Use gocheck !!!!!
	if r.Op != GT {
		t.Fatalf("operator %d should be SMS(%d)", r.Op, GT)
	}
}
