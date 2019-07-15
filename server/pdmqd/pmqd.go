package pdmqd

import (
	"PDMQ/internal/util"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type Client interface {
	//Stats() ClientStats
	IsProducer() bool
}

type PDMQD struct {

	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64

	config   *PDMQDConfig
	topicMap map[string]*Topic

	clientLock sync.RWMutex
	clients    map[int64]Client

	tcpListener  net.Listener
	httpListener net.Listener

	exitChan  chan int
	waitGroup util.WaitGroupWrapper

	startTime time.Time

	sync.RWMutex
}

//
func New(config *PDMQDConfig) (*PDMQD, error) {
	var (
		err error
	)
	pdmqd := &PDMQD{
		config:    config,
		startTime: time.Now(),
		topicMap:  make(map[string]*Topic),
		exitChan:  make(chan int),
	}

	pdmqd.tcpListener, err = net.Listen("tcp", config.TCPAddress)
	if err != nil {
		return nil, fmt.Errorf("listen (%s) failed - %s", config.TCPAddress, err)
	}
	pdmqd.httpListener, err = net.Listen("tcp", config.HTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("listen (%s) failed - %s", config.HTTPAddress, err)
	}

	return pdmqd, err
}

/**
 * @desc pmqd 主进程 开启tcp 跟http 监听
 * @param
 * @return
 */
func (pdmqd *PDMQD) Main() error {
	ctx := &context{pdmqd: pdmqd}
	exitChan := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		//去掉这个once 的话 exitFunc 会执行多次
		once.Do(func() {
			if err != nil {
			}
			fmt.Println(err)
			exitChan <- err
		})
	}
	tcpServer := &tcpServer{ctx: ctx}
	util.PrintJson("tcpServer:", tcpServer)

	//这里捕获退出的方法
	pdmqd.waitGroup.Wrap(func() {
		exitFunc(TCPServer(pdmqd.tcpListener, tcpServer))
	})

	httpServer := newHTTPServer(ctx)

	pdmqd.waitGroup.Wrap(func() {
		exitFunc(HTTPServer(pdmqd.httpListener, httpServer, "http"))
	})
	err := <-exitChan
	return err

	//pdmqd.TcpListen()

	//pdmqd.HttpListen()
}
func testA() error {
	//a := errors.New("hello error")
	return nil
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
	var connect Connect
	for {
		conn, err := topicObject.tcpListener.Accept()
		fmt.Println(conn)
		if err != nil {
			fmt.Println("tcp accept fail", err.Error())
		}
		if tconn, ok := conn.(*net.TCPConn); ok {
			fmt.Println(tconn, pdmqd)
			os.Exit(1)
			go connect.AcceptConnect(pdmqd, tconn)
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

func (pdmqd *PDMQD) AddClient(clientID int64, client Client) {
	pdmqd.clientLock.Lock()
	pdmqd.clients[clientID] = client
	pdmqd.clientLock.Unlock()
}

func (pdmqd *PDMQD) Exit() {
	if pdmqd.tcpListener != nil {
		pdmqd.tcpListener.Close()
	}
	if pdmqd.httpListener != nil {
		pdmqd.httpListener.Close()
	}
	//todo 这里处理异常关闭的消息，存入硬盘
}
