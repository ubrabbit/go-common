package main

import (
	"fmt"
	. "github.com/ubrabbit/go-common/common"
	lib "github.com/ubrabbit/go-common/lib"
	"sync"
	"time"
)

func main() {
	fmt.Printf("\n\n=====================  Producer  =====================\n")

	session := lib.NewRabbitMQSession("Producer", "test", "test", "fanout", "")
	session.Init("rabbitmq", "rabbitmq", "127.0.0.1", 5672, "rabbitmq", nil)
	session.StartProducer()

	i := 0
	go func() {
		for {
			msg := fmt.Sprintf("data_%d", i)
			i++
			session.PushMsg([]byte(msg))
			time.Sleep(1 * time.Microsecond)
			//time.Sleep(1 * time.Second)
		}
	}()

	g := new(sync.WaitGroup)
	g.Add(1)
	go func() {
		time.Sleep(10 * time.Second)
		LogInfo("Close %s", session)
		session.Close()
		g.Done()
	}()
	g.Wait()
}
