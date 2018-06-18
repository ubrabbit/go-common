package tests

import (
	"fmt"
	"testing"
)

import (
	. "github.com/ubrabbit/go-common/common"
)

func TestCommonFunction(t *testing.T) {
	fmt.Printf("\n\n=====================  TestCommonFunction  =====================\n")

	fmt.Println("JoinString", JoinString("/", "aaa", "ccc", "dddd"))

	decodeStr := EncodeBase64("aAbBcCdDHello")
	fmt.Println("EncodeBase64 : ", decodeStr)
	fmt.Println("DecodeBase64 : ", DecodeBase64(decodeStr))
}
