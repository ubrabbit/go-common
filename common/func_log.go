package common

import (
	"fmt"
	log "github.com/ubrabbit/go-common/log"
	debug "github.com/ubrabbit/go-debug"
)

func LogError(format string, v ...interface{}) {
	debug.Print(v...)
	log.Error(format, v...)
}

func LogDebug(format string, v ...interface{}) {
	log.Debug(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	log.Release(format, v...)
}

func LogPanic(format string, v ...interface{}) {
	st := debug.StackTrace(0)
	log.Error(st.String("    "))
	panic(fmt.Sprintf(format, v...))
}

func LogFatal(format string, v ...interface{}) {
	st := debug.StackTrace(0)
	log.Error(st.String("    "))
	log.Fatal(format, v...)
}
