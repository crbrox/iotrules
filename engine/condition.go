package engine

import (
	"fmt"

	"iotrules/mylog"
)

type Condition struct {
	Op       Op
	Exp1     Expression
	Exp2     Expression
	IsNumber bool
}

func (c *Condition) Matched(n *Notif) (matched bool, err error) {
	mylog.Debugf("enter Condition.Matched %+v, %+v", c, n)
	defer func() { mylog.Debugf("exit Condition.Matched %+v, %+v ", matched, err) }()

	if c.IsNumber {
		value1, err := c.Exp1.getNumber(n)
		if err != nil {
			return false, err
		}
		value2, err := c.Exp2.getNumber(n)
		if err != nil {
			return false, err
		}
		switch c.Op {
		case EQ:
			if value1 == value2 {
				matched = true
			}
		case NE:
			if value1 != value2 {
				matched = true
			}
		case LT:
			if value1 < value2 {
				matched = true
			}
		case GT:
			if value1 > value2 {
				matched = true
			}
		default:
			return false, fmt.Errorf("unknown relational operator %v", c.Op)
		}
	} else {
		value1, err := c.Exp1.getString(n)
		if err != nil {
			return false, err
		}
		value2, err := c.Exp2.getString(n)
		if err != nil {
			return false, err
		}
		switch c.Op {
		case EQ:
			if value1 == value2 {
				matched = true
			}
		case NE:
			if value1 != value2 {
				matched = true
			}
		case LT:
			if value1 < value2 {
				matched = true
			}
		case GT:
			if value1 > value2 {
				matched = true
			}
		default:
			return false, fmt.Errorf("unknown relational operator %v", c.Op)
		}
	}
	return matched, nil
}
