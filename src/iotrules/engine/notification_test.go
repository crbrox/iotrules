// notification_test.go
package engine

import (
	"testing"
	"time"
)

func TestNotifGetString(t *testing.T) {
	var campos = map[string]interface{}{"uno": 1, "cadena": "alfanumerica"}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	s, err := n.GetString("$uno")
	if err != nil {
		t.Error(err)
	}
	if s != "1" {
		t.Errorf("%q != \"1\"", s)
	}
	s, err = n.GetString("$cadena")
	if err != nil {
		t.Fatal(err)
	}
	if s != "alfanumerica" {
		t.Fatalf("%q != \"1\"", s)
	}

}

func TestNotifGetNumber(t *testing.T) {
	var campos = map[string]interface{}{"uno": 1, "cadena": "alfanumerica"}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}

	i, err := n.GetNumber("$uno")
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatalf("%f != 1", i)
	}
}

func TestNotifGetNestedNumber(t *testing.T) {
	type dic map[string]interface{}
	var campos = map[string]interface{}{
		"uno": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": 1}},
		"cadena": "alfanumérica"}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}
	i, err := n.GetNumber("$uno.otro.mas")
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatalf("%f != 1", i)
	}
}

func TestNotifGetNestedString(t *testing.T) {
	var campos = map[string]interface{}{
		"cadena": map[string]interface{}{
			"otro": map[string]interface{}{
				"mas": "alfanumérica"}},
		"uno": 1.0}
	var n = Notif{ID: "example_id", Received: time.Now(), Data: campos}
	i, err := n.GetString("$cadena.otro.mas")
	if err != nil {
		t.Fatal(err)
	}
	if i != "alfanumérica" {
		t.Fatalf("%s != alfanumérica", i)
	}
}
