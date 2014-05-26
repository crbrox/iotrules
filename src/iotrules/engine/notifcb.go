package engine

import (
	"encoding/json"
	"iotrules/mylog"
)

type NotifyContextRequest struct {
	SubscriptionId   string
	Originator       string
	ContextResponses []struct{ ContextElement ContextElement }
}
type ContextElement struct {
	Id         string
	IsPattern  string
	Type       string
	Attributes []Attribute
}
type Attribute struct {
	Name  string
	Type  string
	Value string
}

func NewNotifFromCB(ngsi []byte, service int) (n *Notif, err error) {
	mylog.Debugf("enter NewNotifFromCB(%s,%d)\n", ngsi, service)
	defer func() { mylog.Debugf("exit NewNotifFromCB (%+v,%v)\n", n, err) }()

	n = &Notif{Data: map[string]interface{}{}}

	var ncr NotifyContextRequest
	err = json.Unmarshal(ngsi, &ncr)
	if err != nil {
		return nil, err
	}
	mylog.Debugf("in NewNotifFromCB NotifyContextRequest: %+v\n", ncr)

	n.Data["id"] = ncr.ContextResponses[0].ContextElement.Id
	n.Data["type"] = ncr.ContextResponses[0].ContextElement.Type
	for _, attr := range ncr.ContextResponses[0].ContextElement.Attributes {
		n.Data[attr.Name] = attr.Value
	}

	return n, nil
}
