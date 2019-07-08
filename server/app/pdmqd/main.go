package main

import (
	argv "PDMQ/server/argv"
	"PDMQ/server/pdmqd"
	"github.com/judwhite/go-svc/svc"
	"log"
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
		log.Fatal("%v", err.Error())
	}
}

func (p *program) Init(env svc.Environment) error {
	return nil
}

func (p *program) Start() error {
	config := argv.ParseFlag()

	p.pdmqd = pdmqd.New(config)
	go func() {
		//开启pmqd

		p.pdmqd.Main()
	}()
	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.pdmqd.Exit()
	})
	return nil
}
