// iotrules project main.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"iotrules/config"
	"iotrules/mylog"
	"iotrules/web"
)

func main() {
	var err error

	fmt.Println("Hello World!")
	mylog.SetLevel("debug")
	err = config.LoadConfig("iotrules.conf")
	if err != nil {
		mylog.Alert(err)
		os.Exit(-1)
	}

	mux, err := web.Mux()
	if err != nil {
		mylog.Alert(err)
	}
	http.Handle("/", mux)
	err = http.ListenAndServe(config.Port(), nil)
	if err != nil {
		mylog.Alert(err)
	}

}
