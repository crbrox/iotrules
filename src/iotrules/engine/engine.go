package engine

import (
	"sync"

	"iotrules/mylog"
)

type Engine struct {
	Rules []*Rule //Too simple, better data structure is required. Only for starting
	sync.RWMutex
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
	mylog.Debugf("enter Engine.AddRule %+v, %+v", e, r)
	defer func() { mylog.Debugf("exit Engine.AddRule %+v", err) }()

	e.Lock()
	defer e.Unlock()
	e.Rules = append(e.Rules, r)
	return nil
}
