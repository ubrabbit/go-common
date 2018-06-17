package tests

import (
	"fmt"
	"testing"
)

import (
	. "github.com/ubrabbit/go-common/common"
	debug "github.com/ubrabbit/go-debug"
)

func Test_DebugStackTrace(t *testing.T) {
	fmt.Printf("\n\n=====================  Test_DebugStackTrace  =====================\n")

	si := debug.StackTrace(0)
	LogInfo("\n" + si.String("    "))
}

func Test_DebugPrint(t *testing.T) {
	fmt.Printf("\n\n=====================  Test_DebugPrint  =====================\n")

	debug.Print("abc", 123)
}

func Test_DebugPause(t *testing.T) {
	fmt.Printf("\n\n=====================  Test_DebugPause  =====================\n")

	debug.Pause(true)
}
