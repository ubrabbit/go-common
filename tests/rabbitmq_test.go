package tests

import (
	"fmt"
	. "github.com/ubrabbit/go-common/common"
	config "github.com/ubrabbit/go-common/config"
	lib "github.com/ubrabbit/go-common/lib"
	"testing"
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

func TestRabbitMQ_Consumer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRabbitMQ_Consumer  =====================\n")

	config.InitConfig("config_test.conf")
	config.InitRabbitMQConfig()
	cfg := config.GetRabbitMQConfig()

	handle := &RabbitMQCmd{Name: "Consumer"}
	session := lib.NewRabbitMQSession("Consumer", "test", "test", "fanout", "")
	session.Init(cfg.Account, cfg.Password, cfg.Host, cfg.Port, cfg.HostName, handle)
	go session.StartConsumer()

	go func() {
		time.Sleep(5 * time.Second)
		LogInfo("Close %s", session)
		session.Close()
	}()
}

func TestRabbitMQ_Producer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestRabbitMQ_Producer  =====================\n")

	config.InitConfig("config_test.conf")
	config.InitRabbitMQConfig()
	cfg := config.GetRabbitMQConfig()

	handle := &RabbitMQCmd{Name: "Producer"}
	session := lib.NewRabbitMQSession("Producer", "test", "test", "fanout", "")
	session.Init(cfg.Account, cfg.Password, cfg.Host, cfg.Port, cfg.HostName, handle)
	session.StartProducer()

	go func() {
		i := 0
		for {
			msg := fmt.Sprintf("data_%d", i)
			i++
			session.PushMsg([]byte(msg))
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(10 * time.Second)
}
