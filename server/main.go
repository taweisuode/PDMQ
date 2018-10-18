package main

import (
	argv "PDMQ/server/argv"
	"PDMQ/server/pdmqd"
)

/**
* @desc pdmqd 是一个分布式的消息队列，代码灵感来源于nsq,
   希望能够通过编写pdmqd，达到学习go语言的作用
*/
func main() {

	//载入配置信息
	config := argv.ParseFlag()
	//开启pmqd
	pdmqd.Start(config)
}
