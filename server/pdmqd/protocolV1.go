/**
 * @Time : 2019-07-15 12:02
 * @Author : zhuangjingpeng
 * @File : protocolV1
 * @Desc : file function description
 */
package pdmqd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cihub/seelog"
	"io"
	"net"
	"sync/atomic"
)

type protocolV1 struct {
	ctx *context
}

const (
	ProtocolCommonResponse  int32 = 1
	ProtocolErrorResponse   int32 = 2
	ProtocolMessageResponse int32 = 3
)

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

		fmt.Printf("params: %+v\n", string(params[0]))
		response, err := p.Exec(client, params)

		fmt.Printf("response: %+v\n", string(response))

		if err != nil || response == nil {
			seelog.Errorf("response is [%v],err is [%v]", response, err)
			continue
		}
		if response != nil {
			err = p.Send(client, ProtocolCommonResponse, response)
			if err != nil {
				seelog.Errorf("send common response is [%v],err is [%v]", response, err)
				break
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

	//msgTimeCount := client.MsgTimeout
	subEventChan := client.SubEventChan
	close(startChan)
	for {
		if subChannel != nil {
			memoryMsgChan = subChannel.memoryMsgChan

			fmt.Printf("receive memoryMsgChan data ", subChannel.messageCount)
		} else {
			fmt.Println("cannot receive memoryMsgChan data ")
		}
		select {
		case subChannel = <-subEventChan:
			subEventChan = nil
			fmt.Printf("subChannel [%+v]\n", subChannel)
			fmt.Printf("subChannel is topic name is [%+v] channel name is [%+v]\n", subChannel.topicName, subChannel.ChannelName)
		case msg := <-memoryMsgChan:
			fmt.Printf("memoryMsgChan is [%+v]\n", msg)
			msg.Attempts++
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
	if err != nil {
		seelog.Infof("msg write buf total [%d] error [%v]\n", total, err.Error())
	}

	seelog.Infof("msg write buf [%v]\n", string(buf.Bytes()))
	if err != nil {
		seelog.Errorf(" protocolV1 sendmessage error %v\n", err.Error())
		return err
	}
	err = p.Send(client, ProtocolMessageResponse, buf.Bytes())
	if err != nil {
		seelog.Errorf(" protocolV1 send error %v\n", err.Error())
		return err
	}
	return nil
}

func (p *protocolV1) Send(client *clientV1, protocolType int32, buf []byte) error {
	//client.Lock()

	fmt.Printf("send buf is [%+v]\n", string(buf))
	_, err := p.SendFramedResponse(client.Writer, protocolType, buf)
	if err != nil {
		return err
	}
	return err

}

func (p *protocolV1) Exec(client *clientV1, params [][]byte) ([]byte, error) {
	if bytes.Equal(params[0], []byte("IDENTIFY")) {
		return []byte("yes"), nil
	}
	switch {
	case bytes.Equal(params[0], []byte("SUB")):
		return p.SUB(client, params)
	}
	return []byte("error"), nil
}

func (p *protocolV1) SUB(client *clientV1, params [][]byte) ([]byte, error) {
	topicName := string(params[1])
	channelName := string(params[2])

	var channel *Channel
	for {
		topic := p.ctx.pdmqd.GetTopic(topicName)
		channel = topic.GetChannel(channelName)

		fmt.Println(topic, channel)
		if err := channel.AddClient(client.ID, client); err != nil {
			seelog.Errorf("channel consumers for %s:%s exceeds limit of %d", topicName, channelName, p.ctx.pdmqd.config.MaxChannelConsumers)
		}
		break
	}
	client.Channel = channel
	client.SubEventChan <- channel
	return []byte("ok"), nil
}

// SendFramedResponse is a server side utility function to prefix data with a length header
// and frame header and write to the supplied Writer
func (p *protocolV1) SendFramedResponse(w io.Writer, protocolType int32, data []byte) (int, error) {
	beBuf := make([]byte, 4)
	size := uint32(len(data)) + 4

	binary.BigEndian.PutUint32(beBuf, size)
	fmt.Println("fmt.Println(beBuf)", beBuf)
	n, err := w.Write(beBuf)
	if err != nil {
		return n, err
	}

	binary.BigEndian.PutUint32(beBuf, uint32(protocolType))
	fmt.Println("2fmt.Println(beBuf)", beBuf)
	n, err = w.Write(beBuf)
	if err != nil {
		return n + 4, err
	}

	n, err = w.Write(data)
	fmt.Println("3fmt.Println(beBuf)", string(data))

	return n + 8, err
}
