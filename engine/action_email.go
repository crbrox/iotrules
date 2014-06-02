// email
package engine

import (
	"bytes"
	"net/smtp"

	"iotrules/config"
	"iotrules/mylog"
)

type EmailAction struct {
	*ActionData
}

func (a *EmailAction) Do(n *Notif) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter EmailAction.Do %+v %+v", a, n)
		defer func() { mylog.Debugf("exit EmailAction.Do  %+v", err) }()
	}

	var buffer bytes.Buffer
	err = a.ActionData.template.Execute(&buffer, n.Data)
	err = smtp.SendMail(config.SMTPServer(), nil,
		a.Parameters["from"],
		[]string{a.Parameters["to"]},
		buffer.Bytes())
	return err
}
