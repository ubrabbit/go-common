package tests

import (
	"fmt"
	log "github.com/ubrabbit/go-public/log"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Printf("\n\n=====================  TestLog  =====================\n")

	log.Debug("Debug Log 111")
	log.Release("Release Log 111")
	log.Error("Error Log 111")

	log.InitLogger("release", "")
	log.Debug("Debug Log 222")
	log.Release("Release Log 222")

	log.Debug("Debug Log 333")
	log.Release("Release Log 333")

	log.InitLogger("debug", "")
}
