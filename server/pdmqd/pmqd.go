package pdmqd

import (
	"PDMQ/internal/util"
	"fmt"
	"github.com/cihub/seelog"
	"io"
	"net"
	"sync"
	"time"
)

type Client interface {
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

//pdmq 初始化配置项
func New(config *PDMQDConfig) (*PDMQD, error) {
	var (
		err error
	)
	pdmqd := &PDMQD{
		config:    config,
		startTime: time.Now(),
		topicMap:  make(map[string]*Topic),
		clients:   make(map[int64]Client),
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

//pmqd 主进程 开启tcp 跟http 监听
func (pdmqd *PDMQD) Main() error {
	ctx := &context{pdmqd: pdmqd}
	exitChan := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		//去掉这个once 的话 exitFunc 会执行多次
		once.Do(func() {
			if err != nil {
				seelog.Errorf("[PDMQD] [%+v],pdmq occur error [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
			}
			exitChan <- err
		})
	}
	tcpServer := &tcpServer{ctx: ctx}
	//这里捕获退出的方法
	pdmqd.waitGroup.Wrap(func() {
		exitFunc(TCPServer(pdmqd.tcpListener, tcpServer))
	})

	httpServer := newHTTPServer(ctx)

	pdmqd.waitGroup.Wrap(func() {
		exitFunc(HTTPServer(pdmqd.httpListener, httpServer, "http"))
	})

	fmt.Printf("[PDMQD] [%+v] TCP: listening on [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), pdmqd.config.TCPAddress)
	fmt.Printf("[PDMQD] [%+v] HTTP: listening on [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), pdmqd.config.HTTPAddress)

	pdmqd.waitGroup.Wrap(pdmqd.loop)
	err := <-exitChan
	return err
}
func (pdmqd *PDMQD) TcpListen() {
	fmt.Println(pdmqd.config)
	tcpListener, err := net.Listen("tcp", pdmqd.config.TCPAddress)
	topicObject := &PDMQD{
		tcpListener: tcpListener,
	}
	if err != nil {
		seelog.Errorf("[PDMQD] [%+v],tcp connect fail [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	var connect Connect

	for {
		conn, err := topicObject.tcpListener.Accept()
		if err != nil {
			seelog.Errorf("[PDMQD] [%+v],tcp accept fail [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		}
		if tconn, ok := conn.(*net.TCPConn); ok {
			go connect.AcceptConnect(pdmqd, tconn)
		}
	}
}

func (pdmqd *PDMQD) HttpListen() {
	httpListener, err := net.Listen("tcp", pdmqd.config.HTTPAddress)
	topicObject := &PDMQD{
		httpListener: httpListener,
	}
	if err != nil {
		seelog.Errorf("[PDMQD] [%+v],http connect fail [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	for {
		conn, err := topicObject.tcpListener.Accept()
		if err != nil {
			seelog.Errorf("[PDMQD] [%+v],tcp accept fail [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		}
		var buf = make([]byte, 32)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			break
		} else {
			if string(buf[:n]) == "exit" {
				break
			}
		}
	}
}

func (pdmqd *PDMQD) RealTCPAddr() *net.TCPAddr {
	return pdmqd.tcpListener.Addr().(*net.TCPAddr)
}

func (pdmqd *PDMQD) RealHTTPAddr() *net.TCPAddr {
	return pdmqd.httpListener.Addr().(*net.TCPAddr)
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
