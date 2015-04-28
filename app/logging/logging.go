// Logging for AppEngine with line numbers, etc.
// Should be able to drop-in replacement for a non-AppEngine implementation
package logging

import (
	"appengine"
	"net/http"
	"runtime"
	"strconv"
)

type logFunc func(string, ...interface{})

func Infof(req *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(req)
	f := c.Infof
	logMsg(f, format, args...)
}

func Debugf(req *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(req)
	f := c.Debugf
	logMsg(f, format, args...)
}

func Warningf(req *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(req)
	f := c.Warningf
	logMsg(f, format, args...)
}

func Errorf(req *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(req)
	f := c.Errorf
	logMsg(f, format, args...)
}

func Criticalf(req *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(req)
	f := c.Criticalf
	logMsg(f, format, args...)
}

func logMsg(f logFunc, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	format = "[" + file + ":" + strconv.Itoa(line) + "] " + format
	f(format, args...)
}
