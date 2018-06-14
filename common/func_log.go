package common

import (
	log "github.com/ubrabbit/go-common/log"
)

func LogError(format string, v ...interface{}) {
	log.Error(format, v...)
}

func LogDebug(format string, v ...interface{}) {
	log.Debug(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	log.Release(format, v...)
}

func LogFatal(format string, v ...interface{}) {
	log.Fatal(format, v...)
}
