package main

import (
	"PDMQ/server/pdmqd"
	"github.com/cihub/seelog"
	"github.com/judwhite/go-svc/svc"
	"log"
	"os"
	"sync"
	"syscall"
)

/**
* @desc pdmqd 是一个分布式的消息队列，代码灵感来源于nsq,
   希望能够通过编写pdmqd，达到学习go语言的作用
*/

type program struct {
	once  sync.Once
	pdmqd *pdmqd.PDMQD
}

func main() {
	program := &program{}
	if err := svc.Run(program, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatalf("%v", err.Error())
	}
}

func (p *program) Init(env svc.Environment) error {
	return nil
}

func (p *program) Start() (err error) {
	initConf := pdmqd.InitConfig()
	config := pdmqd.ParseFlag(initConf)
	p.pdmqd, err = pdmqd.New(config)
	if err != nil {
		seelog.Error("new pdmqd error is ", err)
	}
	if err != nil {
		seelog.Errorf("failed to persist metadata - %s ", err)
	}

	go func() {
		//开启pmqd
		err := p.pdmqd.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)

		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.pdmqd.Exit()
	})
	return nil
}
