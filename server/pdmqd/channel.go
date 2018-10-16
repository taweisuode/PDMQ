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
