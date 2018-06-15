package tests

import (
	"fmt"
	. "github.com/ubrabbit/go-common/common"
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
	"testing"
	"time"
)

func ReceiverFunc1(b []byte) (bool, error) {
	LogInfo("ReceiverFunc1: %s", string(b))
	return true, nil
}

func ReceiverFunc2(b []byte) (bool, error) {
	LogInfo("ReceiverFunc2: %s", string(b))
	return true, nil
}

func OnConfirm(args ...interface{}) {
	par1 := args[0].(int)
	par2 := args[1].(string)
	tag := args[2].(uint64)
	ack := args[3].(bool)

	LogInfo("OnConfirm : %d %s", tag, ack, par1, par2)
}

func TestRabbitMQ_Consumer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRabbitMQ_Consumer  =====================\n")

	config.InitConfig("config_test.conf")

	session := lib.NewRabbitMQSession("test")
	session.AddReceiver("test", ReceiverFunc1)
	session.AddReceiver("test2", ReceiverFunc2)

	go func() {
		session.ConsumeMsg()
	}()
	time.Sleep(1 * time.Second)

}

func TestRabbitMQ_Producer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRabbitMQ_Producer  =====================\n")

	config.InitConfig("config_test.conf")
	session := lib.NewRabbitMQSession("test")
	session.SetPushConfirm(NewFunctor("callback", OnConfirm, 1, "2"))
	go func() {
		i := 0
		for {
			msg := fmt.Sprintf("data_%d", i)
			i++
			session.PushMsg([]byte(msg))
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(1 * time.Second)
}
