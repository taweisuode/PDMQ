/**
 * @Time : 2019-07-15 12:02
 * @Author : zhuangjingpeng
 * @File : protocolV1
 * @Desc : file function description
 */
package pdmqd

import (
	"fmt"
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
	fmt.Println(conn, err)
}
