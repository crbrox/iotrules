// update
package engine

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"text/template"

	"iotrules/config"
)

const updateTemplateText = `{
    "contextElements": [
        {
            "type": "{{.type}}",
            "isPattern": "{{.isPattern}}",
            "id": "{{.id}}",
            "attributes": [
            {
                "name": "{{.__attrName}}",
                "type": "{{.__attrType}}",
                "value": "{{.__attrValue}}"
            }
            ]
        }
    ],
    "updateAction": "APPEND"
}`

var updateTemplate = template.Must(template.New("updateTemplate").Parse(updateTemplateText))

type UpdateAction struct {
	*ActionData
}

func (a *UpdateAction) Do(n *Notif) (err error) {

	var buffer bytes.Buffer

	// A litle (or very) dirty. Add "hidden" parameters as notification data
	// and remove them after executing template (better copy first level fields in a new map?)
	n.Data["__attrName"] = a.ActionData.Parameters["name"]
	n.Data["__attrValue"] = a.ActionData.Parameters["value"]
	n.Data["__attrType"] = a.ActionData.Parameters["type"]
	err = updateTemplate.Execute(&buffer, n.Data)
	fmt.Println(buffer.String())
	fmt.Println(a.Parameters)
	delete(n.Data, "__attrName")
	delete(n.Data, "__attrValue")
	delete(n.Data, "__attrType")

	req, err := http.NewRequest("POST", config.UpdateEndpoint(), &buffer)
	req.Header.Add("Content-Type", "application/json")
	fmt.Printf("\n****** REQUEST\n%s %v\n******\n\n", buffer.String(), err)
	client := &http.Client{} // reusar entre acciones ??
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	respDump, err := httputil.DumpResponse(resp, true)
	fmt.Printf("\n****** RESPONSE\n%s %v\n******\n\n", respDump, err)
	fmt.Println(resp)
	return err
}
