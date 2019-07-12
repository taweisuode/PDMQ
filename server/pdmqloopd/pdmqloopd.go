/**
 * @Time : 2019-07-12 15:25
 * @Author : zhuangjingpeng
 * @File : pdmqloopd
 * @Desc : file function description
 */
package pdmqloopd

import (
	"PDMQ/internal/util"
	"github.com/cihub/seelog"
	"net"
	"sync"
	"time"
)

type PDMQLOOPD struct {
	sync.RWMutex
	config *PDMQLOOPDConfig

	tcpListener  net.Listener
	httpListener net.Listener

	exitChan  chan int
	waitGroup util.WaitGroupWrapper

	startTime time.Time
}

func New(config *PDMQLOOPDConfig) (*PDMQLOOPD, error) {
	var (
		err error
	)
	pdmqloopd := &PDMQLOOPD{
		config:    config,
		startTime: time.Now(),
		exitChan:  make(chan int),
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

func (loopd *PDMQLOOPD) Main() error {
	ctx := &context{loopd}

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
	loopd.waitGroup.Wrap(func() {
		exitFunc(TCPServer(loopd.tcpListener, tcpServer))
	})

	httpServer := newHTTPServer(ctx)
	loopd.waitGroup.Wrap(func() {
		exitFunc(HTTPServer(loopd.httpListener, httpServer, "http"))
	})

	err := <-exitCh
	return err
}

func (loopd *PDMQLOOPD) Exit() {
	if loopd.tcpListener != nil {
		loopd.tcpListener.Close()
	}

	if loopd.httpListener != nil {
		loopd.httpListener.Close()
	}
	loopd.waitGroup.Wait()
}
