/**
 * @Time : 2019-07-15 12:02
 * @Author : zhuangjingpeng
 * @File : protocolV1
 * @Desc : file function description
 */
package pdmqd

import (
	"bytes"
	"fmt"
	"github.com/cihub/seelog"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type protocolV1 struct {
	ctx *context
}

func (p *protocolV1) IOLoop(connect net.Conn) error {
	var (
		err  error
		line []byte
	)
	clientID := atomic.AddInt64(&p.ctx.pdmqd.clientIDSequence, 1)
	client := newClientV1(clientID, connect, p.ctx)

	p.ctx.pdmqd.AddClient(clientID, client)

	messagePushStartedChan := make(chan bool)
	go p.messagePush(client, messagePushStartedChan)
	<-messagePushStartedChan

	for {
		line, err = client.Reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				seelog.Errorf("failed to read command - %s", err)
			}
			break
		}
		params := bytes.Split(line, []byte(" "))

		fmt.Printf("[PDMQD] [%+v] get consumer request params: %+v\n", time.Now().Format("2006-01-02 15:04:05"), string(params[0]))

		//这里跟nsq 不同的地方在于 我将所有信息都生成为Message，方便统一调动
		msg, err := p.Exec(client, params)
		if err != nil && msg == nil {
			seelog.Errorf("response is [%v],err is [%v]", msg, err)
			continue
		}
		if params[0] != nil && msg != nil {
			err = p.SendMessage(client, msg)
			if err != nil {
				seelog.Errorf("response is [%v],err is [%v]", msg, err)
				continue
			}
		}
	}
	connect.Close()
	close(client.ExitChan)
	return err
}

func (p *protocolV1) messagePush(client *clientV1, startChan chan bool) {
	var (
		err           error
		memoryMsgChan chan *Message
		subChannel    *Channel
	)

	subEventChan := client.SubEventChan
	close(startChan)
	for {
		if subChannel != nil {
			memoryMsgChan = subChannel.memoryMsgChan
		} else {
			seelog.Error("cannot receive memoryMsgChan data ")
		}
		select {
		case subChannel = <-subEventChan:
			//消费端注册到pdmqd 会往subEventChan中投递channel消息，这块用来接收这个消息
			subEventChan = nil
			seelog.Infof("subChannel is topic name is [%+v] channel name is [%+v]\n", subChannel.topicName, subChannel.ChannelName)
		case msg := <-memoryMsgChan:
			//pub消息最终会落在这个case中
			msg.Attempts++
			msg.ProtocolType = ProtocolMessageResponse
			err = p.SendMessage(client, msg)
			if err != nil {
				seelog.Errorf("protocolV1 send message error %v\n", err.Error())
				goto exit
			}
		case <-client.ExitChan:
			//fmt.Printf("receive client [%+v] exitchan\n", subChannel.ChannelName)
			goto exit
		}
	}
exit:
	seelog.Infof("protocolV1: [%s] exiting messagePump", client)
	if err != nil {
		seelog.Errorf("protocolV1 sendmessage error %v\n", err.Error())
	}
}

func (p *protocolV1) SendMessage(client *clientV1, msg *Message) error {
	buf := &bytes.Buffer{}

	total, err := msg.WriteTo(buf)

	fmt.Printf("[PDMQD] [%+v] msg ID is [%+v],msg body is [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), string(msg.ID[:]), string(msg.Body))
	if err != nil {
		seelog.Infof("msg write buf total [%d] error [%v]\n", total, err.Error())
	}

	seelog.Infof("msg write buf [%v]\n", string(buf.Bytes()))
	if err != nil {
		seelog.Errorf(" protocolV1 sendmessage error %v\n", err.Error())
		return err
	}
	err = p.Send(client, buf.Bytes())
	if err != nil {
		seelog.Errorf(" protocolV1 send error %v\n", err.Error())
		return err
	}
	return nil
}

func (p *protocolV1) Send(client *clientV1, buf []byte) error {

	_, err := p.SendProtocolResponse(client, buf)
	if err != nil {
		client.Unlock()
	}
	return err

}

func (p *protocolV1) Exec(client *clientV1, params [][]byte) (*Message, error) {
	var (
		err error
	)
	topicName := string(params[1])
	topic := p.ctx.pdmqd.GetTopic(topicName)
	msg := CreateMessage(topic.GenerateID(), params[0])

	switch {
	case bytes.Equal(params[0], []byte("FIN")):
		return nil, nil
	case bytes.Equal(params[0], []byte("IDENTIFY")):
		return nil, nil
	case bytes.Equal(params[0], []byte("SUB")):
		if err = p.SUB(client, params); err != nil {
			return nil, err
		}
	case bytes.Equal(params[0], []byte("RDY")):
		if err = p.RDY(client, params); err != nil {

		}
	}
	return msg, nil
}

//客户端注册时 调用sub 请求
//往SubEventChan中投递，而这个chan 在pdmqd中的messagePush 中 接收
func (p *protocolV1) SUB(client *clientV1, params [][]byte) error {
	topicName := string(params[1])
	channelName := string(params[2])
	var channel *Channel
	for {
		topic := p.ctx.pdmqd.GetTopic(topicName)
		channel = topic.GetChannel(channelName)
		if err := channel.AddClient(client.ID, client); err != nil {
			seelog.Errorf("channel consumers for %s:%s exceeds limit of %d", topicName, channelName, p.ctx.pdmqd.config.MaxChannelConsumers)
		}
		break
	}
	client.Channel = channel
	client.SubEventChan <- channel
	return nil
}

//Todo：这一版不做处理，需要为消费端限流做反馈
func (p *protocolV1) RDY(client *clientV1, params [][]byte) error {
	return nil
}

//发送该协议统一的返回信息
func (p *protocolV1) SendProtocolResponse(w io.Writer, data []byte) (int, error) {
	n, err := w.Write(data)
	seelog.Infof("write to client data is [%+v], len is [%d]\n", string(data), n)

	return n, err
}
