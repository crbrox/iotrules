// server.go
package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"iotrules/engine"
	"iotrules/mylog"
)

var Server WebServer

func Mux(services []string) (mux *http.ServeMux, err error) {
	mylog.Debugf("enter Mux")
	defer func() { mylog.Debugf("exit Mux %+v %v", mux, err) }()

	engines := map[string]*engine.Engine{}
	for _, service := range services {
		engines[service], err = engine.NewEngine()
		if err != nil {
			return nil, err
		}
	}
	Server := &WebServer{engines}
	mux = http.NewServeMux()
	mux.HandleFunc("/rules/", handlePanic(Server.Rules))
	mux.HandleFunc("/notif/", handlePanic(Server.Notif))
	mux.HandleFunc("/notifCB/", handlePanic(Server.NotifCB))
	return mux, nil

}

type WebServer struct {
	engines map[string]*engine.Engine
}

func (ws *WebServer) Rules(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := strings.TrimPrefix(r.URL.Path, "/rules/")
	service := r.Header.Get("fiware-service")
	serviceEngine := ws.engines[service]
	if serviceEngine == nil {
		http.Error(w, "Service does not exist", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rule, err := engine.NewRule(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = serviceEngine.AddRule(rule)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "{\"id\":%q}\n", rule.ID)
	case "DELETE":
		err := serviceEngine.DeleteRule(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "{\"id\":%q}\n", id)
	case "GET":
		if id != "" {
			r, err := serviceEngine.GetRule(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data, err := r.RuleJSON().ToJSON()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "%s\n", data)
		} else {
			rs, err := serviceEngine.GetAllRules()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var texts [][]byte
			for _, r := range rs {
				data, err := r.RuleJSON().ToJSON()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				texts = append(texts, data)
			}
			fmt.Fprintf(w, "[%s]\n", bytes.Join(texts, []byte{',', '\n'}))
		}
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return

	}
}
func (ws *WebServer) Notif(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	service := r.Header.Get("fiware-service")
	serviceEngine := ws.engines[service]
	if serviceEngine == nil {
		http.Error(w, "Service does not exist", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		n, err := engine.NewNotif(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go func() {
			serviceEngine.Process(n)
		}()
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return

	}
}
func (ws *WebServer) NotifCB(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	service := r.Header.Get("fiware-service")
	serviceEngine := ws.engines[service]
	if serviceEngine == nil {
		http.Error(w, "Service does not exist", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		n, err := engine.NewNotifFromCB(data, 600)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go func() {
			serviceEngine.Process(n)
		}()
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return

	}
}
func handlePanic(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if excp := recover(); excp != nil {
					mylog.Alert(excp)
				}
			}()
			f(w, r)
		})
}
