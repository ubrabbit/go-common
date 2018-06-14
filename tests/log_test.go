package tests

import (
	log "github.com/ubrabbit/go-common/log"
)

func TestLog(t *testing.T) {
	log.Debug("Debug Log 111")
	log.Release("Release Log 111")
	log.Error("Error Log 111")

	log.InitLogger("release", "")
	log.Debug("Debug Log 222")
	log.Release("Release Log 222")

	log.Debug("Debug Log 333")
	log.Release("Release Log 333")
}
