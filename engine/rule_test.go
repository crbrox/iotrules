package engine

import (
	"testing"
	"time"
)

func TestSimpleMatch(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0,
		"numero": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "2.0"}}}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var (
		condAlfa = Condition{Op: EQ, Exp1: Expression{Reference: "cadena.otro.mas"}, Exp2: Expression{Text: "alfanumérica"}, IsNumber: false}
		condOne  = Condition{Op: EQ, Exp1: Expression{Reference: "uno"}, Exp2: Expression{Number: 1.0}, IsNumber: true}
		condTwo  = Condition{Op: EQ, Exp1: Expression{Reference: "numero.otro.mas"}, Exp2: Expression{Number: 2.0}, IsNumber: true}
	)

	var cases = [][]Condition{
		{condAlfa}, {condOne}, {condTwo},
		{condAlfa, condOne}, {condOne, condAlfa}, {condTwo, condOne}, {condOne, condTwo},
		{condAlfa, condOne, condTwo}, {condOne, condAlfa, condTwo}, {condTwo, condOne, condAlfa}, {condOne, condTwo, condAlfa},
	}

	for _, c := range cases {
		var rule = Rule{Conds: c}
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
		"uno": 1.0,
		"numero": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "2.0"}}}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var (
		condFalse = Condition{Op: EQ, Exp1: Expression{Reference: "uno"}, Exp2: Expression{Number: 1.0000001}, IsNumber: true}
		condAlfa  = Condition{Op: EQ, Exp1: Expression{Reference: "cadena.otro.mas"}, Exp2: Expression{Text: "alfanumérica"}, IsNumber: false}
		condOne   = Condition{Op: EQ, Exp1: Expression{Reference: "uno"}, Exp2: Expression{Number: 1.0}, IsNumber: true}
		condTwo   = Condition{Op: EQ, Exp1: Expression{Reference: "numero.otro.mas"}, Exp2: Expression{Number: 2.0}, IsNumber: true}
	)

	var cases = [][]Condition{
		{condFalse},
		{condAlfa, condFalse}, {condOne, condFalse}, {condTwo, condFalse},
		{condFalse, condAlfa}, {condFalse, condOne}, {condFalse, condTwo},
		{condFalse, condFalse},
		{condAlfa, condOne, condTwo, condFalse},
		{condOne, condFalse, condAlfa, condTwo},
		{condTwo, condOne, condFalse, condAlfa},
		{condOne, condTwo, condAlfa, condFalse},
	}

	for _, c := range cases {
		var rule = Rule{Conds: c}
		b, e := rule.Matched(&n)
		if e != nil {
			t.Fatal(e)
		}
		if b {
			t.Fatalf("%+v should not match %+v\n", n, rule)
		}

	}

}
func TestSimpleNotFound(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0,
		"numero": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "2.0"}}}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	var (
		condInexistent = Condition{Op: EQ, Exp1: Expression{Reference: "cadena.otro.jalapeño"}, Exp2: Expression{Number: 1.0}, IsNumber: true}
		condAlfa       = Condition{Op: EQ, Exp1: Expression{Reference: "cadena.otro.mas"}, Exp2: Expression{Text: "alfanumérica"}, IsNumber: false}
		condOne        = Condition{Op: EQ, Exp1: Expression{Reference: "uno"}, Exp2: Expression{Number: 1.0}, IsNumber: true}
		condTwo        = Condition{Op: EQ, Exp1: Expression{Reference: "numero.otro.mas"}, Exp2: Expression{Number: 2.0}, IsNumber: true}
	)

	var cases = [][]Condition{
		{condInexistent},
		{condAlfa, condInexistent}, {condOne, condInexistent}, {condTwo, condInexistent},
		{condInexistent, condAlfa}, {condInexistent, condOne}, {condInexistent, condTwo},
		{condInexistent, condInexistent},
		{condAlfa, condOne, condTwo, condInexistent},
		{condOne, condInexistent, condAlfa, condTwo},
		{condTwo, condOne, condInexistent, condAlfa},
		{condOne, condTwo, condAlfa, condInexistent},
	}

	for _, c := range cases {
		var rule = Rule{Conds: c}
		b, e := rule.Matched(&n)
		if b {
			t.Fatalf("%+v should not match %+v\n", n, rule)
		}
		if e != nil {
			if _, ok := e.(*ErrorNotFound); !ok {
				t.Fatal(e)
			}
		}

	}

}
