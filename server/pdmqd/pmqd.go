package pdmqd

import (
	"PDMQ/server/argv"
	"fmt"
	"io"
	"net"
)

type PDMQD struct {
	config   flag.Config
	topicMap map[string]*Topic

	tcpListener  net.Listener
	httpListener net.Listener
}

func Start(config *flag.ArgvConfig) {
	pdmqd := New(config)
	pdmqd.Main()
}

//
func New(config *flag.ArgvConfig) (pdmqd *PDMQD) {
	currentConfig := flag.Config{TCPAddress: config.TcpListen, HTTPAddress: config.HttpListen}
	return &PDMQD{config: currentConfig}
}

/**
 * @desc pmqd 主进程 开启tcp 跟http 监听
 * @param
 * @return
 */
func (pdmqd *PDMQD) Main() {

	pdmqd.TcpListen()

	pdmqd.HttpListen()
}
func (pdmqd *PDMQD) TcpListen() {
	fmt.Println(pdmqd.config)
	tcpListener, err := net.Listen("tcp", pdmqd.config.TCPAddress)
	topicObject := &PDMQD{
		tcpListener: tcpListener,
	}
	if err != nil {
		fmt.Println("tcp connect fail", err.Error())
		return
	}
	for {
		conn, err := topicObject.tcpListener.Accept()
		fmt.Println(conn)
		if err != nil {
			fmt.Println("tcp accept fail", err.Error())
		}
		if tconn, ok := conn.(*net.TCPConn); ok {
			go AcceptConn(pdmqd, tconn)
		}
	}
	fmt.Println("hello world")
}

func (pdmqd *PDMQD) HttpListen() {
	httpListener, err := net.Listen("tcp", pdmqd.config.HTTPAddress)
	topicObject := &PDMQD{
		httpListener: httpListener,
	}
	if err != nil {
		fmt.Println("http connect fail", err.Error())
		return
	}
	//id := 0

	for {
		conn, err := topicObject.tcpListener.Accept()
		if err != nil {
			fmt.Println("tcp accept fail", err.Error())
		}
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
	/*for {
		conn, err := topicObject.tcpListener.Accept()
		if err != nil {
			fmt.Println("tcp accept fail", err.Error())
		}
		id++
		if tconn, ok := conn.(*net.TCPConn); ok {
			go HandleConn(tconn, id)
		}

	}*/
}
func HandleConn(conn *net.TCPConn) {
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
func AcceptConn(pdmqd *PDMQD, conn *net.TCPConn) {
	for {
		var buf = make([]byte, 32)
		n, err := conn.Read(buf)
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
