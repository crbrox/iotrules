// notification.go
package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"iotrules/mylog"
)

type Notif struct {
	ID       string
	Received time.Time
	Data     map[string]interface{}
}

func NewNotif(data []byte) (n *Notif, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter NewNotif %q", data)
		defer func() { mylog.Debugf("exit NewNotif %+v, %+v", n, err) }()
	}

	n = &Notif{}
	n.ID = uuid.New()
	n.Received = time.Now()
	err = json.Unmarshal(data, &n.Data)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Notif) GetNumber(exp string) (number float64, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Notif.GetNumber %q", exp)
		defer func() { mylog.Debugf("exit Notif.GetNumber %+v, %+v", number, err) }()
	}

	//We don't care if it comes as a string or as a number.
	i, err := n.getElement(exp)
	if err != nil {
		return 0, err
	}
	switch i := i.(type) {
	case string:
		num, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	case float64:
		return i, nil
	case int:
		return float64(i), nil
	default:
		return 0, fmt.Errorf("%q is not valid as number (%T)", exp, i)
	}
}

func (n *Notif) GetString(exp string) (str string, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Notif.GetString %q", exp)
		defer func() { mylog.Debugf("exit Notif.GetString %+v, %+v", str, err) }()
	}

	//We don't care if it comes as a string or as a number.
	i, err := n.getElement(exp)
	if err != nil {
		return "", err
	}
	switch i := i.(type) {
	case string:
		return i, nil
	case int, float64, bool:
		return fmt.Sprint(i), nil
	default:
		return "", fmt.Errorf("%q is not a basic type (%T)", exp, i)
	}
}

func (n *Notif) getElement(exp string) (str interface{}, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Notif.GetString %q", exp)
		defer func() { mylog.Debugf("exit Notif.GetString %+v, %+v", str, err) }()
	}
	d := n.Data
	// First check the complete field. Allow "a.b.c" as a valid name
	// and a litle faster for non-nested fields.
	if i, ok := d[exp]; ok {
		return i, nil
	}
	// Else we check nested objects
	fields := strings.Split(exp, ".")
	for _, f := range fields[:len(fields)-1] {
		i, ok := d[f]
		if !ok {
			return "", &ErrorNotFound{Field: exp, Part: f, Notif: n}
		}
		d, ok = i.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("%q is not an object in %q", f, exp)
		}
	}
	last := fields[len(fields)-1]
	i, ok := d[last]
	if !ok {
		return "", &ErrorNotFound{Field: exp, Part: last, Notif: n}
	}
	return i, nil

}

type ErrorNotFound struct {
	Field string
	Part  string
	Notif *Notif
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("%q in %q not found in notification", e.Part, e.Field)
}
