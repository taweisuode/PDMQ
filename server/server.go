package main

import (
	argv "PDMQ/server/argv"
	"PDMQ/server/pdmqd"
)

func main() {

	//载入配置信息
	argv.ParseFlag()

	//进行tcp 监听
	pmqd.TcpListen()
	pmqd.HttpListen()
}
