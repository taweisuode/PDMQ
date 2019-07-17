package pdmqd

import (
	"PDMQ/server/waitGroup"
	"fmt"
	"github.com/cihub/seelog"
	"sync"
	"sync/atomic"
	"time"
)

type Topic struct {

	//msg数量跟大小累加
	messageCount uint64
	messageBytes uint64

	topicName     string
	channelMap    map[string]*Channel
	memoryMsgChan chan *Message
	waitGroup     waitGroup.WaitGroupWrapper
	ctx           *context

	startChan chan int
	exitChan  chan int

	idFactory *guidFactory

	paused    int
	pauseChan chan int

	sync.RWMutex
}

func CreateTopic(topicName string, ctx *context) *Topic {
	t := &Topic{
		topicName:     topicName,
		channelMap:    make(map[string]*Channel),
		memoryMsgChan: make(chan *Message, ctx.pdmqd.config.MsgChanSize),
		ctx:           ctx,
		paused:        0,
		startChan:     make(chan int, 10),
		pauseChan:     make(chan int),
		idFactory:     NewGUIDFactory(ctx.pdmqd.config.ID),
	}
	t.waitGroup.Wrap(t.msgOutput)
	return t
}

/**
 *  @desc:  这里就是处理topic 的方法
 *  @input: data String
 *  @resp:  err resp
 *
**/
func (topic *Topic) msgOutput() {
	var msg *Message
	//var buf []byte

	var chans []*Channel
	var memoryMsgChan chan *Message
	for {
		select {
		case <-topic.startChan:
		}
		break
	}
	topic.RLock()

	for _, c := range topic.channelMap {
		chans = append(chans, c)
	}
	memoryMsgChan = topic.memoryMsgChan
	for {
		select {
		case msg = <-memoryMsgChan:
		case <-topic.exitChan:
			goto exit
		}
		//为每个channel 投递消息
		for _, channel := range chans {

			fmt.Printf("topic is [%+v], channel is [%+v],msg is [%+v]\n", channel.topicName, channel.ChannelName, string(msg.Body))
			err := channel.PutMessage(msg)
			if err != nil {
				seelog.Infof("TOPIC(%s) ERROR: failed to put msg(%s) to channel(%s) - %s", topic.topicName, msg.ID, channel.ChannelName, err)
			}
		}
	}
exit:
	seelog.Infof("topic(%s): closing ... messageOutput", topic.topicName)
}

func (pdmqd *PDMQD) GetTopic(topicName string) *Topic {
	pdmqd.RLock()
	topic, ok := pdmqd.topicMap[topicName]

	pdmqd.RUnlock()
	if ok {
		return topic
	}
	pdmqd.Lock()
	topic = CreateTopic(topicName, &context{pdmqd: pdmqd})
	pdmqd.topicMap[topicName] = topic

	pdmqd.Unlock()

	topic.Start()
	return topic
}

func (topic *Topic) PutMessage(msg *Message) error {
	topic.RLock()
	defer topic.RUnlock()
	err := topic.put(msg)
	if err != nil {
		seelog.Errorf("topic(%s) put message(%s) error(%v)", topic.topicName, msg.ID, err.Error())
		return err
	}
	atomic.AddUint64(&topic.messageCount, 1)
	atomic.AddUint64(&topic.messageBytes, uint64(len(msg.Body)))
	return nil
}

func (topic *Topic) put(msg *Message) error {
	select {
	case topic.memoryMsgChan <- msg:
	default:
	}
	return nil
}

/**
 *  @desc:  now that all channels are added, start topic msgOutput
 *  @input: data String
 *  @resp:  err resp
 *
**/
func (topic *Topic) Start() {
	topic.startChan <- 1
	/*select {
	case
	default:
	}*/
}
func (topic *Topic) GetChannel(channelName string) *Channel {
	topic.Lock()
	channel, ok := topic.channelMap[channelName]
	if !ok {
		channel = CreateChannel(topic.topicName, channelName, topic.ctx)
	}
	topic.channelMap[channelName] = channel

	topic.Unlock()

	return channel
}

func (t *Topic) GenerateID() MessageID {
retry:
	id, err := t.idFactory.NewGUID()
	if err != nil {
		time.Sleep(time.Millisecond)
		goto retry
	}
	return id.Hex()
}
