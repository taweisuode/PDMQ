package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":9999", 2*time.Second)
	if err != nil {
		fmt.Println("tcp connect error ", err.Error())
		return
	}
	defer conn.Close()
	for {
		var buf = make([]byte, 32)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			break
		} else {
			if string(buf[:n]) == "exit" {
				fmt.Println("connect exit")
				break
			}
			fmt.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
		}
	}
}
