package test

import (
	"PDMQ/server/argv"
	"PDMQ/server/pdmqd"
	"strconv"
	"testing"
	"time"
)

func PutMessageTest(t testing.T) {
	config := flag.ParseFlag()
	pdmqd := pdmqd.New(config)

	topicName := "test_put_message" + strconv.Itoa(int(time.Now().Unix()))
	topic := pdmqd.GetTopic(topicName)
	channel1 := topic.GetChannel("ch")

	var id MessageID
	msg := NewMessage(id, []byte("test"))
	topic.PutMessage(msg)

	outputMsg := <-channel1.memoryMsgChan
	test.Equal(t, msg.ID, outputMsg.ID)
	test.Equal(t, msg.Body, outputMsg.Body)
}
