package engine

import (
	"code.google.com/p/go-uuid/uuid"

	"iotrules/mylog"
)

type Rule struct {
	ID     string
	Conds  []Condition
	Action Action
}

func NewRule(body []byte) (r *Rule, err error) {
	mylog.Debugf("enter NewRule %q", body)
	defer func() { mylog.Debugf("exit NewRule %+v, %+v", r, err) }()

	rj, err := ParseRuleJSON(body)
	if err != nil {
		return nil, err
	}
	r, err = rj.Rule()
	if err != nil {
		return nil, err
	}
	r.ID = uuid.New()
	return r, nil
}

func (r *Rule) Do(n *Notif) (err error) {
	mylog.Debugf("enter Rule.Do %+v %+v", r, n)
	defer func() { mylog.Debugf("exit Rule.Do %+v", err) }()

	matched, err := r.Matched(n)
	if err == nil && matched {
		err = r.Action.Do(n)
	}
	return err
}
func (r *Rule) Matched(n *Notif) (matched bool, err error) {
	mylog.Debugf("enter Rule.Matched %+v %+v", r, n)
	defer func() { mylog.Debugf("exit Rule.Matched %+v  %+v", matched, err) }()

	for _, cond := range r.Conds {
		matched, err := cond.Matched(n)
		if !matched || err != nil {
			return false, err
		}
	}
	return true, nil
}
func (r *Rule) RuleJSON() (rj *RuleJSON) {
	rj = &RuleJSON{ID: r.ID}
	cjs := make([]CondJSON, 0, len(r.Conds))
	for _, c := range r.Conds {
		e1 := c.Exp1.ToRuleJSON()
		e2 := c.Exp2.ToRuleJSON()
		opStr := c.Op.String()
		typeStr := "string"
		if c.IsNumber {
			typeStr = "number"
		}
		cJSON := CondJSON{Type: typeStr, Expr: [3]interface{}{e1, opStr, e2}}
		cjs = append(cjs, cJSON)
	}
	rj.And = cjs
	axnR := r.Action.Data()
	rj.Action = ActionJSON{
		Type:       axnR.Type.String(),
		Template:   axnR.TemplateText,
		Parameters: axnR.Parameters}
	return rj

}
