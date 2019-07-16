/**
 * @Time : 2019-07-15 12:02
 * @Author : zhuangjingpeng
 * @File : protocolV1
 * @Desc : file function description
 */
package pdmqd

import (
	"PDMQ/internal/test"
	"fmt"
	"github.com/nsqio/go-nsq"
	"io"
	"net"
	"testing"
	"time"
)

func startPDMQD(config *PDMQDConfig) (*net.TCPAddr, *net.TCPAddr, *PDMQD) {
	config.TCPAddress = "127.0.0.1:9400"
	config.HTTPAddress = "127.0.0.1:9401"

	nsqd, err := New(config)
	if err != nil {
		panic(err)
	}
	go func() {
		err := nsqd.Main()
		if err != nil {
			panic(err)
		}
	}()
	return nsqd.RealTCPAddr(), nsqd.RealHTTPAddr(), nsqd
}
func mustConnectPDMQD(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	conn.Write([]byte("V1"))
	return conn, nil
}
func TestProtocolV1_Base(t *testing.T) {
	config := InitConfig()
	tcpAddr, httpAddr, pdmqd := startPDMQD(config)
	fmt.Println(tcpAddr, httpAddr)
	fmt.Printf("pdmqd %+v", pdmqd)

	topic := pdmqd.GetTopic("hello")
	msg := CreateMessage(topic.GenerateID(), []byte("world"))
	topic.PutMessage(msg)

	conn, err := mustConnectPDMQD(tcpAddr)
	sub(t, conn, "world", "ch")

	resp, err := nsq.ReadResponse(conn)
	test.Nil(t, err)
	frameType, data, err := nsq.UnpackResponse(resp)
	msgOut, _ := decodeMessage(data)
	test.Equal(t, frameTypeMessage, frameType)
	test.Equal(t, msg.ID, msgOut.ID)
	test.Equal(t, msg.Body, msgOut.Body)
	test.Equal(t, uint16(1), msgOut.Attempts)
	fmt.Println(conn, err)
}
func sub(t *testing.T, conn io.ReadWriter, topicName string, channelName string) {
	total, err := nsq.Subscribe(topicName, channelName).WriteTo(conn)
	fmt.Println("111", total, err)
	//readValidate(t, conn, frameTypeResponse, "OK")
}
