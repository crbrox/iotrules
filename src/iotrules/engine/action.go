// action
package engine

import (
	"fmt"
	"os"
	"text/template"

	"iotrules/mylog"
)

type ActionType int

// QUE TODO VAYA POR HTTP CON DIFERENTES PLANTILLAS INTERNAS?????
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
		s = "EMAIL"
	case UPDATE:
		s = "UPDATE"
	case HTTP:
		s = "HTTP"
	}
	return fmt.Sprintf("%s(%d)", s, t)
}
func ParseActionType(s string) (at ActionType, err error) {
	mylog.Debugf("enter ParseActionType %q", s)
	defer func() { mylog.Debugf("exit ParseActionType %+v, %+v ", at, err) }()

	var finalType ActionType
	switch s {
	case "SMS":
		finalType = SMS
	default:
		return 0, fmt.Errorf("unknown action type %q", s)
	}
	return finalType, nil
}

type Action interface {
	Do(n *Notif) error
}

type ActionData struct {
	Type         ActionType
	TemplateText string
	Parameters   map[string]string
	template     *template.Template
}

func NewAction(actionType ActionType, templateString string, parameters map[string]string) (axn Action, err error) {
	mylog.Debugf("enter NewAction %+v, %q, %+v", actionType, templateString, parameters)
	defer func() { mylog.Debugf("exit NewAction %+v, %+v ", axn, err) }()

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

type SMSAction struct {
	*ActionData
}

func (a *SMSAction) Do(n *Notif) (err error) {
	mylog.Debugf("enter SMSAction.Do %+v %+v", a, n)
	defer mylog.Debugf("exit SMSAction.Do  %+v", err)

	err = a.ActionData.template.Execute(os.Stdout, n.Data)
	return err
}

type EmailAction struct {
	*ActionData
}

func (a *EmailAction) Do(n *Notif) error {
	return nil
}

type UpdateAction struct {
	*ActionData
}

func (a *UpdateAction) Do(n *Notif) error {
	return nil
}

type HTTPAction struct {
	*ActionData
}

func (a *HTTPAction) Do(n *Notif) error {
	return nil
}
