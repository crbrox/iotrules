// action
package engine

import (
	"fmt"
	"text/template"

	"iotrules/mylog"
)

type ActionType int

const (
	SMS = ActionType(iota + 1)
	EMAIL
	UPDATE
	HTTP
)

func (t ActionType) String() string {
	var s string
	switch t {
	case SMS:
		s = "SMS"
	case EMAIL:
		s = "email"
	case UPDATE:
		s = "update"
	case HTTP:
		s = "HTTP"
	}
	return s
}
func ParseActionType(s string) (at ActionType, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter ParseActionType %q", s)
		defer func() { mylog.Debugf("exit ParseActionType %+v, %+v ", at, err) }()
	}

	var finalType ActionType
	switch s {
	case "SMS":
		finalType = SMS
	case "email":
		finalType = EMAIL
	case "update":
		finalType = UPDATE
	default:
		return 0, fmt.Errorf("unknown action type %q", s)
	}
	return finalType, nil
}

type Action interface {
	Do(n *Notif) error
	Data() *ActionData
}

type ActionData struct {
	Type         ActionType
	TemplateText string
	Parameters   map[string]string
	template     *template.Template
}

func (ad *ActionData) Data() *ActionData {
	return ad
}

func NewAction(actionType ActionType, templateString string, parameters map[string]string) (axn Action, err error) {
	if mylog.Debugging {
		mylog.Debugf("enter NewAction %+v, %q, %+v", actionType, templateString, parameters)
		defer func() { mylog.Debugf("exit NewAction %+v, %+v ", axn, err) }()
	}

	var ad = &ActionData{Type: actionType, TemplateText: templateString, Parameters: parameters}
	t, err := template.New("").Parse(templateString)
	if err != nil {
		return nil, err
	}
	ad.template = t

	switch actionType {
	case SMS:
		axn = &SMSAction{ad}
	case EMAIL:
		axn = &EmailAction{ad}
	case UPDATE:
		axn = &UpdateAction{ad}
	case HTTP:
		axn = &HTTPAction{ad}
	default:
		return nil, fmt.Errorf("unknown action type %v", t)
	}
	return axn, nil
}

type HTTPAction struct {
	*ActionData
}

func (a *HTTPAction) Do(n *Notif) error {
	return nil
}
