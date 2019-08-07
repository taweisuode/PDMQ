/**
 * @Time : 2019-07-12 15:25
 * @Author : zhuangjingpeng
 * @File : pdmqloopd
 * @Desc : file function description
 */
package pdmqloopd

import (
	"PDMQ/internal/util"
	"fmt"
	"github.com/cihub/seelog"
	"net"
	"sync"
	"time"
)

type PDMQLOOPD struct {
	sync.RWMutex
	config  *PDMQLOOPDConfig
	version string

	tcpListener  net.Listener
	httpListener net.Listener

	startTime time.Time
	exitChan  chan int
	waitGroup util.WaitGroupWrapper
	DB        *RegistrationDB
}

func New(config *PDMQLOOPDConfig) (*PDMQLOOPD, error) {
	var (
		err error
	)
	pdmqloopd := &PDMQLOOPD{
		config:    config,
		startTime: time.Now(),
		exitChan:  make(chan int),
		DB:        NewRegistrationDB(),
		version:   "1.0V1",
	}

	pdmqloopd.tcpListener, err = net.Listen("tcp", config.TCPAddress)
	if err != nil {
		return nil, seelog.Errorf("listen (%s) failed - %s", config.TCPAddress, err)
	}
	pdmqloopd.httpListener, err = net.Listen("tcp", config.HTTPAddress)
	if err != nil {
		return nil, seelog.Errorf("listen (%s) failed - %s", config.HTTPAddress, err)
	}

	return pdmqloopd, err
}

func (pdmqloopd *PDMQLOOPD) Main() error {
	ctx := &context{pdmqloopd}

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				seelog.Error("main start error is ", err)
			}
			exitCh <- err
		})
	}

	tcpServer := &tcpServer{ctx: ctx}
	pdmqloopd.waitGroup.Wrap(func() {
		exitFunc(TCPServer(pdmqloopd.tcpListener, tcpServer))
	})

	httpServer := newHTTPServer(ctx)
	pdmqloopd.waitGroup.Wrap(func() {
		exitFunc(HTTPServer(pdmqloopd.httpListener, httpServer, "http"))
	})

	fmt.Printf("[PDMQLOOPD] [%+v] TCP: listening on [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), pdmqloopd.config.TCPAddress)
	fmt.Printf("[PDMQLOOPD] [%+v] HTTP: listening on [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), pdmqloopd.config.HTTPAddress)
	err := <-exitCh
	return err
}

func (pdmqloopd *PDMQLOOPD) RealTCPAddr() *net.TCPAddr {
	return pdmqloopd.tcpListener.Addr().(*net.TCPAddr)
}

func (pdmqloopd *PDMQLOOPD) RealHTTPAddr() *net.TCPAddr {
	return pdmqloopd.httpListener.Addr().(*net.TCPAddr)
}

func (pdmqloopd *PDMQLOOPD) Exit() {
	if pdmqloopd.tcpListener != nil {
		pdmqloopd.tcpListener.Close()
	}

	if pdmqloopd.httpListener != nil {
		pdmqloopd.httpListener.Close()
	}
	pdmqloopd.waitGroup.Wait()
}
