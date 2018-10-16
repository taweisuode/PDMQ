package pdmqd

type Topic struct {
	topicName     string
	channelMap    map[string]*channel
	memoryMsgChan chan *Message
}

func CreateTopic() {

}
