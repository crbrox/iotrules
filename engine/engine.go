package engine

import (
	"fmt"
	"sync"

	"iotrules/mylog"
)

type Engine struct {
	Rules map[string]*Rule //Too simple, better data structure/s is/are required. Only for starting
	sync.RWMutex
}

func NewEngine() (eng *Engine, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter NewEngine")
		defer func() { mylog.Debugf("exit NewEngine %+v, %+v", eng, err) }()
	}

	eng = &Engine{Rules: map[string]*Rule{}}
	return eng, nil
}
func (e *Engine) Process(n *Notif) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Engine.Process %+v, %+v", e, n)
		defer func() { mylog.Debugf("exit Engine.Process %+v", err) }()
	}

	e.RLock()
	defer e.RUnlock()
	for _, rule := range e.Rules {
		rule.Do(n)
	}
	return nil
}
func (e *Engine) AddRule(r *Rule) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Engine.AddRule %+v", r)
		defer func() { mylog.Debugf("exit Engine.AddRule %+v", err) }()
	}

	e.Lock()
	defer e.Unlock()
	e.Rules[r.ID] = r
	return nil
}
func (e *Engine) DeleteRule(id string) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Engine.DeleteRule %q", id)
		defer func() { mylog.Debugf("exit Engine.DeleteRule %+v", err) }()
	}

	e.Lock()
	defer e.Unlock()
	delete(e.Rules, id)
	return nil
}
func (e *Engine) GetRule(id string) (r *Rule, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Engine.GetRule %q", id)
		defer func() { mylog.Debugf("exit Engine.GetRule %+v, %+v", r, err) }()
	}

	e.RLock()
	defer e.RUnlock()
	r, ok := e.Rules[id]
	if !ok {
		return nil, fmt.Errorf("rule id %q not found", id)
	}
	return r, nil

}
func (e *Engine) GetAllRules() (rs []*Rule, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter Engine.GetAllRules")
		defer func() { mylog.Debugf("exit Engine.GetAllRules %+v, %+v", rs, err) }()
	}

	e.RLock()
	defer e.RUnlock()
	for _, r := range e.Rules {
		rs = append(rs, r)
	}
	return rs, nil

}
