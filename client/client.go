package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	conn, err := net.DialTimeout("tcp", ":9998", 2*time.Second)
	if err != nil {
		fmt.Println("tcp connect error ", err.Error())
		return
	}
	defer conn.Close()
	// 下面进行读写
	var wg sync.WaitGroup
	wg.Add(2)
	go handleWrite(conn, &wg)
	//go handleRead(conn, &wg)
	wg.Wait()
	/*for {
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
	}*/
	/*for {
		for {
			conn, err := topicObject.tcpListener.Accept()
			if err != nil {
				fmt.Println("tcp accept fail", err.Error())
			}
			id++
			if tconn, ok := conn.(*net.TCPConn); ok {
				go HandleConn(tconn, id)
			}
		}
		fmt.Println("connect close")
		defer conn.Close()
	}*/
}

func handleWrite(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("what would you do ?")
	work := ""
	for {
		fmt.Scanf("%s", &work)
		switch work {
		case "create_topic":
			fmt.Println("send your message")
			topic := ""
			message := ""
			fmt.Scanf("%s,%s", &topic, &message)

			_, err := conn.Write(([]byte(topic + "")))
			if err != nil {
				fmt.Println("write data error", err)
			}
			continue
		case "exit":
			os.Exit(1)
			break
		default:
			continue
		}
	}
	fmt.Println("connect close")
	defer conn.Close()
}
