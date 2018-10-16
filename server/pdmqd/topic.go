package pdmqd

import (
	"PDMQ/server/waitGroup"
)

type Topic struct {
	topicName     string
	channelMap    map[string]*Channel
	memoryMsgChan chan *Message
	waitGroup     waitGroup.WaitGroupWrapper
}

func CreateTopic(topicName string, ctx *context) *Topic {
	t := &Topic{
		topicName:     topicName,
		channelMap:    make(map[string]*Channel),
		memoryMsgChan: make(chan *Message, ctx.pdmqd.config.MsgChanSize),
	}

	t.waitGroup.Wrap(t.msgOutput)
	return t
}

func (*Topic) msgOutput() {

}
