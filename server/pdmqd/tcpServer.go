/**
 * @Time : 2019-07-09 15:53
 * @Author : zhuangjingpeng
 * @File : tcpServer
 * @Desc : file function description
 */
package pdmqd

import (
	"fmt"
	"net"
)

type TCPHandler interface {
	Handle(conn net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) error {
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
			break
		}
		go handler.Handle(clientConn)
	}
	return nil
}
