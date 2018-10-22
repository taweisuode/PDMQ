package main

import (
	"fmt"
	"time"

	"go-nsq"
)

// Producer 生产者
func Producer() {
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Println("NewProducer", err)
		panic(err)
	}

	for i := 0; i < 100; i++ {
		if err := producer.Publish("pdmqd", []byte(fmt.Sprintf("Hello World %d", i))); err != nil {
			fmt.Println("Publish", err)
			panic(err)
		}
		time.Sleep(time.Second * 5)
		i++
	}
}

func main() {
	Producer()
}
