// server.go
package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"iotrules/engine"
	"iotrules/mylog"
)

var Server WebServer

func Mux() *http.ServeMux {
	Server := &WebServer{&engine.Engine{}}
	m := http.NewServeMux()
	m.HandleFunc("/rules/", handlePanic(Server.Rules))
	m.HandleFunc("/notif/", handlePanic(Server.Notif))
	m.HandleFunc("/notifCB/", handlePanic(Server.NotifCB))
	return m

}

type WebServer struct {
	engine *engine.Engine
}

func (ws *WebServer) Rules(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
		go func() {
			ws.engine.AddRule(rule)
		}()
		fmt.Fprintf(w, "%s %s\n", r.Method, r.URL)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return

	}
}
func (ws *WebServer) Notif(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
			ws.engine.Process(n)
		}()
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return

	}
}
func (ws *WebServer) NotifCB(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
			ws.engine.Process(n)
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
