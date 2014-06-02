// sms
package engine

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"iotrules/config"
	"iotrules/mylog"
)

type SMSAction struct {
	*ActionData
}

func (a *SMSAction) Do(n *Notif) (err error) {
	if mylog.Debugging {
		defer func() { mylog.Debugf("exit SMSAction.Do  %+v", err) }()
	}

	var buffer bytes.Buffer
	err = a.ActionData.template.Execute(&buffer, n.Data)
	fmt.Println(buffer)
	fmt.Println(a.Parameters)
	msg := fmt.Sprintf(`{"to":["tel:%s"], "message": %q}`, a.Parameters["to"], buffer.String())
	fmt.Println(msg)
	client := &http.Client{} //reusar entre acciones ??
	req, err := http.NewRequest("POST", config.SMSEndpoint(), strings.NewReader(msg))
	req.Header.Add("API_KEY", config.APIKey())
	req.Header.Add("API_SECRET", config.APISecret())
	req.Header.Add("Content-Type", "application/json")
	fmt.Printf("\n****** REQUEST\n%s %v\n******\n\n", msg, err)
	resp, err := client.Do(req)
	respDump, err := httputil.DumpResponse(resp, true)
	fmt.Printf("\n****** RESPONSE\n%s %v\n******\n\n", respDump, err)
	fmt.Println(resp)
	return err
}
