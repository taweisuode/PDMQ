/**
 * @Time : 2019-07-08 19:24
 * @Author : zhuangjingpeng
 * @File : tcp
 * @Desc : file function description
 */
package pdmqd

import (
	"github.com/cihub/seelog"
	"io"
	"net"
)

type tcpServer struct {
	ctx *context
}

func (tcp *tcpServer) Handle(clientConn net.Conn) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		seelog.Errorf("read buf from clientConn err, %v", err)
		clientConn.Close()
		return
	}
	var protocol Protocol
	protocolVal := string(buf)
	switch protocolVal {
	case "V1":
		protocol = &protocolV1{ctx: tcp.ctx}
	}

	err = protocol.IOLoop(clientConn)
	if err != nil {
		seelog.Errorf("client(%s) - %s", clientConn.RemoteAddr(), err)
		return
	}
}

type Protocol interface {
	IOLoop(conn net.Conn) error
}
