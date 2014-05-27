// iotrules project main.go
package main

import (
	"fmt"
	"net/http"

	"iotrules/mylog"
	"iotrules/web"
)

func main() {
	var err error

	fmt.Println("Hello World!")
	mylog.SetLevel("debug")

	mux, err := web.Mux()
	if err != nil {
		mylog.Alert(err)
	}
	http.Handle("/", mux)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		mylog.Alert(err)
	}

}
