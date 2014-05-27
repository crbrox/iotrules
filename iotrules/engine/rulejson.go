// rulejson.go
package engine

import (
	"encoding/json"
	"fmt"

	"iotrules/mylog"
)

type RuleJSON struct {
	ID     string     `json:"id,omitempty"`
	And    []CondJSON `json:"and"`
	Action ActionJSON `json:"axn"`
}
type CondJSON struct {
	Type string         `json:"type"`
	Expr [3]interface{} `json:"expr"`
}
type ActionJSON struct {
	Type       string            `json:"type"`
	Template   string            `json:"template"`
	Parameters map[string]string `json:"params"`
}

func (rj *RuleJSON) Rule() (pr *Rule, err error) {
	mylog.Debugf("enter RuleJSON.Rule %+v", rj)
	defer func() { mylog.Debugf("exit RuleJSON.Rule %+v, %+v", pr, err) }()

	var r Rule
	condListR := make([]Condition, 0, len(rj.And))
	for _, condJ := range rj.And {
		var condR Condition
		switch t := condJ.Type; t {
		case "number":
			condR.IsNumber = true
		case "string":
			condR.IsNumber = false
		default:
			return nil, fmt.Errorf("unknown condition type %q", t)
		}

		opStr, ok := condJ.Expr[1].(string)
		if !ok {
			return nil, fmt.Errorf("operator should be a string")
		}
		condR.Op, err = parseOp(opStr)
		if err != nil {
			return nil, err
		}

		condR.Exp1, err = makeExpressionFromJSON(condJ.Expr[0], condR.IsNumber)
		if err != nil {
			return nil, err
		}

		condR.Exp2, err = makeExpressionFromJSON(condJ.Expr[2], condR.IsNumber)
		if err != nil {
			return nil, err
		}
		condListR = append(condListR, condR)
	}
	r.Conds = condListR

	at, err := ParseActionType(rj.Action.Type)
	if err != nil {
		return nil, err
	}
	axn, err := NewAction(at, rj.Action.Template, rj.Action.Parameters)
	if err != nil {
		return nil, err
	}
	r.Action = axn
	return &r, nil
}

func ParseRuleJSON(data []byte) (rj *RuleJSON, err error) {
	mylog.Debugf("enter ParseRuleJSON %s", data)
	defer func() { mylog.Debugf("exit ParseRuleJSON %+v  %+v", rj, err) }()

	rj = &RuleJSON{}
	err = json.Unmarshal(data, rj)
	if err != nil {
		return nil, err
	}
	return rj, nil
}
func (rj *RuleJSON) ToJSON() ([]byte, error) {
	return json.MarshalIndent(rj, "", " ")
}
