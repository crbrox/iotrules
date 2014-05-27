package engine

import (
	"sync"

	"iotrules/mylog"
)

type Engine struct {
	Rules map[string]*Rule //Too simple, better data structure/s is/are required. Only for starting
	sync.RWMutex
}

func NewEngine() (eng *Engine, err error) {
	mylog.Debugf("enter NewEngine")
	defer func() { mylog.Debugf("exit NewEngine %+v, %+v", eng, err) }()

	eng = &Engine{Rules: map[string]*Rule{}}
	return eng, nil
}
func (e *Engine) Process(n *Notif) (err error) {
	mylog.Debugf("enter Engine.Process %+v, %+v", e, n)
	defer func() { mylog.Debugf("exit Engine.Process %+v", err) }()

	e.RLock()
	defer e.RUnlock()
	for _, rule := range e.Rules {
		rule.Do(n)
	}
	return nil
}
func (e *Engine) AddRule(r *Rule) (err error) {
	mylog.Debugf("enter Engine.AddRule %+v", r)
	defer func() { mylog.Debugf("exit Engine.AddRule %+v", err) }()

	e.Lock()
	defer e.Unlock()
	e.Rules[r.ID] = r
	return nil
}
func (e *Engine) DeleteRule(id string) (err error) {
	mylog.Debugf("enter Engine.DeleteRule %q", id)
	defer func() { mylog.Debugf("exit Engine.DeleteRule %+v", err) }()

	e.Lock()
	defer e.Unlock()
	delete(e.Rules, id)
	return nil
}
