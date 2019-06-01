package gornir

import (
	"reflect"
	"runtime"
)

type Logger interface {
	Info(...interface{})
	Debug(...interface{})
	Error(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
	WithField(string, interface{}) Logger
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
