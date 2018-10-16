package main

import (
	argv "PDMQ/server/argv"
	"PDMQ/server/pdmqd"
)

func main() {

	//载入配置信息
	config := argv.ParseFlag()
	//开启pmqd
	pdmqd.Start(config)
}
