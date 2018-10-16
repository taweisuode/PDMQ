package pdmqd

type Consumer interface {
	UnPause()
	Pause()
	Close()
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
