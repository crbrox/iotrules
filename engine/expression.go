// expression.go
package engine

import (
	"fmt"

	"iotrules/mylog"
)

type Expression struct {
	Reference string
	Text      string
	Number    float64
}

func (e Expression) getNumber(n *Notif) (number float64, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Expression.getNumber %+v, %+v", e, n)
		defer func() { mylog.Debugf("exit Expression.getNumber %+v, %+v", number, err) }()
	}

	if e.Reference == "" {
		return e.Number, nil
	} else {
		return n.GetNumber(e.Reference)
	}
}

func (e Expression) getString(n *Notif) (str string, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Expression.getString %+v, %+v", e, n)
		defer func() { mylog.Debugf("exit Expression.getString %+v, %+v", str, err) }()
	}

	if e.Reference == "" {
		return e.Text, nil
	} else {
		return n.GetString(e.Reference)
	}
}

func makeExpressionFromJSON(i interface{}, isNumber bool) (exp Expression, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter makeExpressionFromJSON %+v %+v", i, isNumber)
		defer func() { mylog.Debugf("exit makeExpressionFromJSON %+v  %+v", exp, err) }()
	}

	switch i := i.(type) {
	case string:
		if i[0] == '$' {
			exp.Reference = i[1:]
		} else {
			if !isNumber {
				exp.Text = i
			} else {
				return exp, fmt.Errorf("not numerical value %q in numerical condition")
			}
		}
	case int:
		if isNumber {
			exp.Number = float64(i)
		} else {
			return exp, fmt.Errorf("numerical value %v in not numerical condition")
		}
	case float64:
		if isNumber {
			exp.Number = i
		} else {
			return exp, fmt.Errorf("numerical value %v in not numerical condition")
		}
	default:
		return exp, fmt.Errorf("invalid type for expression %T", i)
	}
	return exp, nil
}

func (e Expression) ToRuleJSON() interface{} {
	if e.Reference != "" {
		return "$" + e.Reference
	} else if e.Text != "" {
		return e.Text
	} else {
		return e.Number
	}
}
