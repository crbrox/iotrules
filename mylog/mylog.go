package mylog

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"sync"
)

var (
	dbgLog   *log.Logger
	alertLog *syslog.Writer
	level    int
	mtx      sync.Mutex
)

const (
	ALERT = 0
	INFO  = 6
	DEBUG = 7
)

func init() {
	dbgLog = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	var err error
	alertLog, err = syslog.New(syslog.LOG_EMERG|syslog.LOG_USER, "iotrules")
	if err != nil {
		log.Fatal(err)
	}
}

func SetLevelInt(lvl int) {
	mtx.Lock()
	defer mtx.Unlock()
	level = lvl
}
func SetLevel(lvl string) {
	switch lvl {
	case "alert":
		SetLevelInt(ALERT)
	case "info":
		SetLevelInt(INFO)
	case "debug":
		SetLevelInt(DEBUG)
	default:
		panic("unknown log level: " + lvl)
	}
}
func Logger() *log.Logger {
	return dbgLog
}

func Debugf(f string, args ...interface{}) {
	if level >= DEBUG {
		str := fmt.Sprintf(f, args...)
		dbgLog.Output(2, str)
	}
}

func Debug(args ...interface{}) {
	if level >= DEBUG {
		str := fmt.Sprintln(args...)
		dbgLog.Output(2, str)
	}
}
func Infof(f string, args ...interface{}) {
	if level >= INFO {
		str := fmt.Sprintf(f, args...)
		dbgLog.Output(2, str)
	}
}

func Info(args ...interface{}) {
	if level >= INFO {
		str := fmt.Sprintln(args...)
		dbgLog.Output(2, str)
	}
}
func Alertf(f string, args ...interface{}) {
	str := fmt.Sprintf(f, args...)
	alert(str)
}

func Alert(args ...interface{}) {
	str := fmt.Sprintln(args...)
	alert(str)

}
func alert(str string) {
	err := alertLog.Alert(str)
	if err != nil {
		dbgLog.Println(err) // Panic would be better??
	}
	dbgLog.Output(3, str)
}
