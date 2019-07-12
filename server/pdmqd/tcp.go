/**
 * @Time : 2019-07-08 19:24
 * @Author : zhuangjingpeng
 * @File : tcp
 * @Desc : file function description
 */
package pdmqd

import (
	"fmt"
	"net"
)

type tcpServer struct {
	ctx *context
}

func (tcp *tcpServer) Handle(clientConn net.Conn) {
	/*	buf := make([]byte, 4)
		_, err := io.ReadFull(clientConn, buf)
		if err != nil {
			clientConn.Close()
			return
		}
		var prot Protocol

		err = prot.IOLoop(clientConn)
		if err != nil {
			return
		}*/

	defer clientConn.Close()
	var buf [512]byte
	for {
		n, err := clientConn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := clientConn.RemoteAddr()
		fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		_, err2 := clientConn.Write([]byte("Welcome client!"))
		if err2 != nil {
			return
		}
	}
}

type Protocol interface {
	IOLoop(conn net.Conn) error
}
