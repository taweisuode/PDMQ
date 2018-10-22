package main

import (
	"fmt"
	"go-nsq"
	"time"
)

// ConsumerHandler 消费者处理者
type ConsumerHandler struct{}

// HandleMessage 处理消息
func (*ConsumerHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Println(string(msg.Body))
	return nil
}

// ConsumerA 消费者
func ConsumerA() {
	consumer, err := nsq.NewConsumer("pdmqd", "1", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewConsumer", err)
		panic(err)
	}

	consumer.AddHandler(&ConsumerHandler{})

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println("ConnectToNSQLookupd", err)
		panic(err)
	}
}

// ConsumerB 消费者
func ConsumerB() {
	consumer, err := nsq.NewConsumer("pdmqd", "2", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewConsumer", err)
		panic(err)
	}

	consumer.AddHandler(&ConsumerHandler{})

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		fmt.Println("ConnectToNSQLookupd", err)
		panic(err)
	}
}
func main() {
	for i := 0; i < 100; i++ {
		ConsumerA()
		ConsumerB()
		time.Sleep(time.Second * 1)
	}
}
