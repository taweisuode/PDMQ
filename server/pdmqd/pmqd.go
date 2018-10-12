package pmqd

import (
	"PDMQ/server/argv"
	"fmt"
	"net"
)

type Topic struct {
	config       flag.Config
	tcpListener  net.Listener
	httpListener net.Listener
}

func TcpListen() {
	config := flag.Construct()
	fmt.Println(config.TCPAddress, config.HTTPAddress)
	tcpListener, err := net.Listen("tcp", config.TCPAddress)
	topicObject := &Topic{
		tcpListener: tcpListener,
	}
	if err != nil {
		fmt.Println("tcp connect fail", err.Error())
		return
	}
	id := 0
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
	fmt.Println("hello world")
}

func HttpListen() {
	config := flag.Construct()
	httpListener, err := net.Listen("tcp", config.HTTPAddress)
	topicObject := &Topic{
		httpListener: httpListener,
	}
	if err != nil {
		fmt.Println("http connect fail", err.Error())
		return
	}
	id := 0
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
}
func HandleConn(conn *net.TCPConn, id int) {
	fmt.Println("send your message")
	message := ""
	for {
		fmt.Scanf("%s", &message)
		switch message {
		case "create_topic":
			topic := ""
			channel := ""
			fmt.Printf("please input topic and message", &topic, &channel)
		}
		_, err := conn.Write(([]byte(message)))
		if err != nil {
			fmt.Println("write data error", err)
		}
		if message == "exit" {
			break
		}
	}
	fmt.Println("connect close")
	defer conn.Close()
}
