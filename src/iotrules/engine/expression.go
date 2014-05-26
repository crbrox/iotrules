// expression.go
package engine

import (
	"iotrules/mylog"
)

type Expression struct {
	Reference string
	Text      string
	Number    float64
}

func (e Expression) getNumber(n *Notif) (number float64, err error) {
	mylog.Debugf("enter Expression.getNumber %+v, %+v", e, n)
	defer func() { mylog.Debugf("exit Expression.getNumber %+v, %+v", number, err) }()

	if e.Reference == "" {
		return e.Number, nil
	} else {
		return n.GetNumber(e.Reference)
	}
}

func (e Expression) getString(n *Notif) (str string, err error) {
	mylog.Debugf("enter Expression.getString %+v, %+v", e, n)
	defer func() { mylog.Debugf("exit Expression.getString %+v, %+v", str, err) }()

	if e.Reference == "" {
		return e.Text, nil
	} else {
		return n.GetString(e.Reference)
	}
}
