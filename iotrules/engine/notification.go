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
	mylog.Debugf("enter NewNotif %q", data)
	defer func() { mylog.Debugf("exit NewNotif %+v, %+v", n, err) }()

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
	mylog.Debugf("enter Notif.GetNumber %q", exp)
	defer func() { mylog.Debugf("exit Notif.GetNumber %+v, %+v", number, err) }()

	strval, err := n.GetString(exp)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(strval, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (n *Notif) GetString(exp string) (str string, err error) {
	mylog.Debugf("enter Notif.GetString %q", exp)
	defer func() { mylog.Debugf("exit Notif.GetString %+v, %+v", str, err) }()

	fields := strings.Split(exp, ".")
	d := n.Data
	for _, f := range fields[:len(fields)-1] {
		i, ok := d[f]
		if !ok {
			return "", fmt.Errorf("%q not found in notification", f)
		}
		d, ok = i.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("%q is not an object in %q", f, exp)
		}
	}
	last := fields[len(fields)-1]
	i, ok := d[last]
	if !ok {
		return "", fmt.Errorf("%q not found in notification", last)
	}
	switch i := i.(type) {
	case string:
		return i, nil
	case int, float64, bool:
		return fmt.Sprint(i), nil
	default:
		return "", fmt.Errorf("%q is not an basic type in %q", last, exp)
	}
}
