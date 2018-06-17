package main

import (
	"fmt"
	. "github.com/ubrabbit/go-common/common"
	lib "github.com/ubrabbit/go-common/lib"
	"sync"
	"time"
)

type RabbitMQCmd struct {
	Name string
}

func (self *RabbitMQCmd) OnConfirmMessage(v uint64, b bool) {
	LogInfo("%s OnConfirmMessage: %d %v,", self.Name, v, b)
}

func (self *RabbitMQCmd) OnReceiveMessage(b []byte) (bool, bool) {
	LogInfo("%s OnReceiveMessage: %s", self.Name, string(b))
	return true, false
}

func main() {
	fmt.Printf("\n\n=====================  Consumer  =====================\n")

	handle := &RabbitMQCmd{Name: "Consumer"}
	session := lib.NewRabbitMQSession("Consumer", "test", "test", "fanout", "")
	session.Init("rabbitmq", "rabbitmq", "127.0.0.1", 5672, "rabbitmq", handle)
	go session.StartConsumer()

	g := new(sync.WaitGroup)
	g.Add(1)
	go func() {
		time.Sleep(30 * time.Second)
		LogInfo("Close %s", session)
		session.Close()
		g.Done()
	}()
	g.Wait()
}
