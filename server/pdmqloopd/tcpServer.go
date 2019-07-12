/**
 * @Time : 2019-07-09 15:53
 * @Author : zhuangjingpeng
 * @File : tcpServer
 * @Desc : file function description
 */
package pdmqloopd

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

func Handle(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])

		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		str := string(buf[:n])
		fmt.Printf("receive from client, data: %v\n", str)
	}
}
