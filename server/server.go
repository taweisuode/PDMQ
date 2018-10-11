package main

import (
	argv "pmq/server/argv"
	topic "pmq/server/topic"
)

func main() {

	//载入配置信息
	argv.ParseFlag()

	//进行tcp 监听
	topic.TcpListen()
}
