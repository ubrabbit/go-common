package tests

import (
	"fmt"
	. "github.com/ubrabbit/go-common/common"
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
	"testing"
	"time"
)

func Receiver(b []byte) (bool, error) {
	LogInfo("Receiver: %s", string(b))
	return true, nil
}

func Receiver2(b []byte) (bool, error) {
	LogInfo("Receiver2: %s", string(b))
	return true, nil
}

func TestRabbitMQ(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRabbitMQ  =====================\n")

	config.InitConfig("config_test.conf")

	session := lib.NewRabbitMQSession("test")
	session.AddReceiver("test", Receiver)
	session.AddReceiver("test2", Receiver2)

	go func() {
		session.ConsumeMsg()
	}()
	go func() {

	}()
	time.Sleep(3 * time.Second)

}
