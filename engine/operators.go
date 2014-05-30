// operators
package engine

import "fmt"

type Op int

const (
	EQ = Op(iota + 1)
	NE
	LT
	GT
)

func (o Op) String() string {
	var s string
	switch o {
	case EQ:
		s = "="
	case NE:
		s = "!="
	case LT:
		s = "<"
	case GT:
		s = ">"
	}
	return s
}
func parseOp(op string) (Op, error) {
	var finalOp Op
	switch op {
	case "<":
		finalOp = LT
	case ">":
		finalOp = GT
	case "=":
		finalOp = EQ
	case "!=":
		finalOp = NE
	default:
		return 0, fmt.Errorf("unknown operator %q", op)
	}
	return finalOp, nil
}
