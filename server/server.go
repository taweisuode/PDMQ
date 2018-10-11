package main

import (
	argv "pmq/server/argv"
	"pmq/server/pmqd"
)

func main() {

	//载入配置信息
	argv.ParseFlag()

	//进行tcp 监听
	pmqd.TcpListen()
}
