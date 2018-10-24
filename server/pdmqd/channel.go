package pdmqd

type Consumer interface {
	UnPause()
	Pause()
	Close() error
	TimedOutMessage()
	Empty()
}
type Channel struct {
	topicName   string
	ChannelName string
	ctx         *context

	memoryMsgChan chan *Message
	clients       map[int64]*Consumer
}

func CreateChannel(topicName string, ChannelName string, ctx *context) *Channel {
	return &Channel{
		topicName:     topicName,
		ChannelName:   ChannelName,
		ctx:           ctx,
		memoryMsgChan: make(chan *Message, ctx.pdmqd.config.MsgChanSize),
		clients:       make(map[int64]*Consumer),
	}
}
