package pdmqd

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Connect interface {
	AcceptConnect(pdmqd *PDMQD, conn *net.TCPConn)
}

func AcceptConnect(pdmqd *PDMQD, conn *net.TCPConn) {
	for {
		var buf = make([]byte, 32)
		n, err := conn.Read(buf)
		fmt.Println(buf)
		os.Exit(1)
		CreateTopic(string(buf[:]), &context{pdmqd})
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			break
		} else {
			if string(buf[:n]) == "exit" {
				fmt.Println("connect exit")
				break
			}
			if n != 0 {
				fmt.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
			}
		}
	}

}
