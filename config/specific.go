package config

import (
	"fmt"
	"iotrules/mylog"
)

var port string
var loaded bool = false

func LoadConfig(filename string) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter config.LoadConfig %q", filename)
		defer func() { mylog.Debugf("exit config.LoadConfig %+v", err) }()
	}

	err = loadFile(filename)
	if err != nil {
		return err
	}

	var found bool
	port, found = str("port")
	if !found {
		return fmt.Errorf("config port is mandatory")
	}
	return nil
}

func Port() string { return port }
