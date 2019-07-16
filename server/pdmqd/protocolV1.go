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
)

type protocolV1 struct {
	ctx *context
}

const (
	frameTypeResponse int32 = 0
	frameTypeError    int32 = 1
	frameTypeMessage  int32 = 2
)

func (p *protocolV1) IOLoop(connect net.Conn) error {
	var (
		err  error
		line []byte
	)
	clientID := atomic.AddInt64(&p.ctx.pdmqd.clientIDSequence, 1)
	client := newClientV1(clientID, connect, p.ctx)

	p.ctx.pdmqd.AddClient(clientID, client)

	fmt.Printf("p.ctx.pdmqd [%+v]\n", p.ctx.pdmqd.clients)

	messagePushStartedChan := make(chan bool)
	go p.messagePush(client, messagePushStartedChan)
	<-messagePushStartedChan

	//var buf = make([]byte, 1024)
	//n, err := connect.Read(buf)
	//fmt.Println("123", n, err, string(buf))
	for {
		line, err = client.Reader.ReadSlice('\n')

		fmt.Printf("line [%+v] error [%+v]", string(line), err)
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
			err = p.Send(client, response)
			if err != nil {
				break
			}
		}
	}
	connect.Close()
	close(client.ExitChan)
	return nil
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
		select {
		case subChannel = <-subEventChan:
			subEventChan = nil
			fmt.Println(subChannel)
		case msg := <-memoryMsgChan:
			msg.Attempts++
			err = p.SendMessage(client, msg)
			if err != nil {
				seelog.Errorf("protocolV1 sendmessage error %v\n", err.Error())
				goto exit
			}
		case <-client.ExitChan:
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
	seelog.Infof("msg write buf total [%d] error [%v]\n", total, err.Error())

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
	client.Lock()

	len, err := client.Write(buf)

	fmt.Println(len, err)
	if err != nil {
		client.Unlock()
	}
	return err

}

func (p *protocolV1) Exec(client *clientV1, params [][]byte) ([]byte, error) {
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
		channel := topic.GetChannel(channelName)
		if err := channel.AddClient(client.ID, client); err != nil {
			seelog.Errorf("channel consumers for %s:%s exceeds limit of %d", topicName, channelName, p.ctx.pdmqd.config.MaxChannelConsumers)
		}
		break
	}
	client.Channel = channel
	client.SubEventChan <- channel
	return []byte("ok"), nil
}
