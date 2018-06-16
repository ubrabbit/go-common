package lib

//https://segmentfault.com/a/1190000010516906
//http://www.damonyi.cc/%E5%9F%BA%E4%BA%8Egolang%E5%AE%9E%E7%8E%B0%E7%9A%84rabbitmq-%E8%BF%9E%E6%8E%A5%E6%B1%A0/

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

import (
	rabbitmq "github.com/streadway/amqp"

	. "github.com/ubrabbit/go-common/common"
)

type RabbitMsgReceiver interface {
	OnConfirmMessage(uint64, bool)
	OnReceiveMessage([]byte) (bool, bool)
}

type RabbitMQSession struct {
	sync.Mutex

	Account  string
	Password string
	Host     string
	HostName string
	Port     int

	Name         string
	QueueName    string
	Exchange     string
	ExchangeType string
	RouterKey    string
	Conn         *rabbitmq.Connection
	Channel      *rabbitmq.Channel
	Queue        *rabbitmq.Queue

	onMessageHandle interface{}
	isClosed        bool
}

func NewRabbitMQSession(name string, qname string, ex string, extype string, router string) *RabbitMQSession {
	session := new(RabbitMQSession)
	session.Name = name
	session.QueueName = qname
	session.Exchange = ex
	session.ExchangeType = extype
	if router != "" {
		session.RouterKey = router
	} else {
		session.RouterKey = ex
	}

	session.Conn = nil
	session.Channel = nil
	session.Queue = nil

	session.onMessageHandle = nil
	session.isClosed = false
	return session
}

func (self *RabbitMQSession) String() string {
	return fmt.Sprintf("[RabbitMQ][%s-%s-%s-%s]", self.Name, self.QueueName, self.Exchange, self.ExchangeType)
}

func (self *RabbitMQSession) Init(acct string, pwd string, host string, port int, hostname string, handle interface{}) {
	self.Account = acct
	self.Password = pwd
	self.Host = host
	self.HostName = hostname
	self.Port = port
	self.onMessageHandle = handle
}

func (self *RabbitMQSession) Close() (err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("Close Error: %v", p)
			LogError("%s : %v", self, err)
		}
	}()

	//关闭阶段不处理断线的异常
	self.disconnect()
	self.isClosed = true
	LogInfo("%s Closed", self)
	return nil
}

func (self *RabbitMQSession) IsClosed() bool {
	return self.isClosed
}

func (self *RabbitMQSession) Ping() error {
	if self.Conn == nil || self.Channel == nil {
		return rabbitmq.ErrClosed
	}

	channel := self.Channel
	err := channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
	if err != nil {
		return err
	}

	msgContent := "ping.ping"
	err = channel.Publish("ping.ping", "ping.ping", false, false, rabbitmq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})
	if err != nil {
		return err
	}

	err = channel.ExchangeDelete("ping.ping", false, false)
	return err
}

func (self *RabbitMQSession) ResetExchange(name string, t string) {
	last, last_t := self.Exchange, self.ExchangeType
	self.Exchange = name
	self.ExchangeType = name
	if last != name || last_t != t {
		LogInfo("reconnect by exchange changed")
		self.reconnect()
	}
}

func (self *RabbitMQSession) disconnect() (err error) {
	self.Lock()
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("Disconnect Error: %v", p)
			LogError("%s : %v", self, err)
		}
		self.Unlock()
	}()

	if self.Channel != nil {
		// will close() the deliveries channel
		if err := self.Channel.Cancel(self.QueueName, false); err != nil {
			LogError("%s Cancel %s Error: %v", self, self.QueueName, err)
			return fmt.Errorf("Consumer cancel failed: %s", err)
		}
		self.Channel.Close()
		self.Channel = nil
	}
	if self.Conn != nil {
		if err := self.Conn.Close(); err != nil {
			LogError("%s Close Conn Error: %v", self, err)
			return fmt.Errorf("AMQP connection close error: %s", err)
		}
		self.Conn = nil
		LogInfo("%s Disconnect", self)
	}
	return nil
}

func (self *RabbitMQSession) connect() error {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("Connect Err:  %v", err)
		}
		self.Unlock()
	}()
	if self.IsClosed() {
		return rabbitmq.ErrClosed
	}
	account := self.Account
	password := self.Password
	host := self.Host
	hostname := self.HostName
	port := self.Port

	LogInfo("%s Connect RabbitMQ %s:%d:%s", self, host, port, hostname)
	// 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
	server := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", account, password, host, port, hostname)
	conn, err := rabbitmq.Dial(server)
	if err != nil {
		LogError("Connect RabbitMQ Error: %v", err)
		return err
	}
	LogInfo("%s Connect RabbitMQ Success", self)

	channel, err := conn.Channel()
	if err != nil {
		LogDebug("Channel Error: %v", err)
		return err
	}

	ex, extype := self.Exchange, self.ExchangeType
	if err = channel.ExchangeDeclare(
		ex,     // name of the exchange
		extype, // type
		true,   // durable
		false,  // delete when complete
		false,  // internal
		false,  // noWait
		nil,    // arguments
	); err != nil {
		LogDebug("ExchangeDeclare '%s' '%s' Error: %v", ex, extype, err)
		return err
	}

	//channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
	queue, err := channel.QueueDeclare(
		self.QueueName, // name
		true,           // durable    持久化标识
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		LogDebug("QueueDeclare '%s' Error: %v", self.QueueName, err)
		return err
	}

	self.Conn = conn
	self.Channel = channel
	self.Queue = &queue

	go self.confirmPushMsg(channel.NotifyPublish(make(chan rabbitmq.Confirmation, 1)))
	return nil
}

func (self *RabbitMQSession) reconnect() error {
	for {
		self.disconnect()
		if self.IsClosed() {
			return rabbitmq.ErrClosed
		}
		time.Sleep(time.Second * 1)
		LogInfo("%s reconnect", self)
		err := self.connect()
		if err != nil {
			LogError("%s reconnect error: %v", self, err)
			continue
		}
		err = self.Ping()
		if err != nil {
			LogError("%s reconnect ping error: %v", self, err)
			continue
		}
		LogInfo("%s reconnect success", self)
		return nil
	}
	return errors.New(fmt.Sprintf("%s fail reconnect", self))
}

func (self *RabbitMQSession) publish(data []byte, save bool) error {
	if self.Channel == nil {
		return fmt.Errorf("publishMsg Error by channel is nil")
	}
	exchange, exchangeType, router := self.Exchange, self.ExchangeType, self.RouterKey
	if err := self.Channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	mod := rabbitmq.Transient //temp save
	if save {
		mod = rabbitmq.Persistent // 持久化标记
	}
	err := self.Channel.Publish(
		exchange, // exchange
		router,   // routing key
		false,    // mandatory
		false,    // immediate
		rabbitmq.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: mod,
		})
	return err
}

func (self *RabbitMQSession) confirmPushMsg(c chan rabbitmq.Confirmation) {
	LogInfo("%s waiting confirmPushMsg", self)
	for {
		confirmed, ok := <-c
		if !ok {
			LogInfo("%s break confirmPushMsg by closed", self)
			break
		}
		self.onMessageHandle.(RabbitMsgReceiver).OnConfirmMessage(confirmed.DeliveryTag, confirmed.Ack)
	}
}

func (self *RabbitMQSession) PushMsg(data []byte) bool {
	for {
		if self.IsClosed() {
			break
		}
		err := self.publish(data, true)
		if err != nil {
			LogError("%s PushMsg Error: %v", self, err)
			err := self.reconnect()
			if err != nil {
				LogError("Failed to push message: %s", string(data))
				break
			}
		}
		LogDebug("%s PushMsg: %s", self, string(data))
		return true
	}
	return false
}

func (self *RabbitMQSession) PushMsgNoSave(data []byte) bool {
	for {
		if self.IsClosed() {
			break
		}
		err := self.publish(data, false)
		if err != nil {
			LogError("pushMsg Error: %v", err)
			err := self.reconnect()
			if err != nil {
				LogError("Failed to push message: %s", string(data))
				break
			}
		}
		LogDebug("PushMsgNoSave: %s", string(data))
		return true
	}
	return false
}

func (self *RabbitMQSession) StartProducer() {
	self.connect()
}

func (self *RabbitMQSession) StartConsumer() {
	for {
		self.Lock()
		if self.IsClosed() {
			LogInfo("Consumer finished by closed")
			self.Unlock()
			break
		}
		self.Unlock()

		err := self.reconnect()
		if err != nil {
			LogInfo("Consumer finished by connect failure: %v", err)
			break
		}

		ex := self.Exchange
		if err = self.Channel.QueueBind(
			self.QueueName, // name of the queue
			"",             // bindingKey
			ex,             // sourceExchange
			false,          // noWait
			nil,            // arguments
		); err != nil {
			LogError("Consumer Failed To Queue Bind")
			continue
		}

		ch_msgs, err := self.Channel.Consume(
			self.Queue.Name, // queue
			self.Queue.Name, // consumerTag,
			false,           // auto-ack
			false,           // exclusive
			false,           // no-local
			false,           // no-wait
			nil,             // args
		)
		if err != nil {
			LogError("Consumer Failed To Create Consume")
			continue
		}
		LogDebug("%s >>>>>>>>>>> start consumer", self)
		for msg := range ch_msgs {
			LogDebug(
				"%s got %dB delivery: [%v] %q",
				self,
				len(msg.Body),
				msg.DeliveryTag,
				msg.Body,
			)
			ack, requeue := self.onMessageHandle.(RabbitMsgReceiver).OnReceiveMessage(msg.Body)
			if ack {
				// 确认收到本条消息, multiple必须为false
				msg.Ack(false)
			} else {
				msg.Nack(false, requeue)
			}
		}
		LogDebug("%s >>>>>>>>>>> finished consumer", self)
	}
}
