package pdmqd

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Consumer interface {
	UnPause()
	Pause()
	Close() error
	TimedOutMessage()
	Empty()
}
type Channel struct {
	messageCount uint64 //消息数

	topicName   string
	ChannelName string
	ctx         *context

	sync.RWMutex

	memoryMsgChan chan *Message
	clients       map[int64]Consumer
}

func CreateChannel(topicName string, ChannelName string, ctx *context) *Channel {
	return &Channel{
		topicName:     topicName,
		ChannelName:   ChannelName,
		ctx:           ctx,
		memoryMsgChan: make(chan *Message, ctx.pdmqd.config.MsgChanSize),
		clients:       make(map[int64]Consumer),
	}
}

func (c *Channel) PutMessage(msg *Message) error {
	c.RLock()
	defer c.RUnlock()
	if err := c.put(msg); err != nil {
		return err
	}
	atomic.AddUint64(&c.messageCount, 1)
	return nil
}

func (c *Channel) put(msg *Message) error {
	select {
	case c.memoryMsgChan <- msg:
	}
	return nil
}

// AddClient adds a client to the Channel's client list
func (c *Channel) AddClient(clientID int64, client Consumer) error {
	c.Lock()
	defer c.Unlock()

	_, ok := c.clients[clientID]
	if ok {
		return nil
	}

	maxChannelConsumers := c.ctx.pdmqd.config.MaxChannelConsumers
	if maxChannelConsumers != 0 && len(c.clients) >= maxChannelConsumers {
		return errors.New("E_TOO_MANY_CHANNEL_CONSUMERS")
	}

	c.clients[clientID] = client
	return nil
}
