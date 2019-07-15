/**
 * @Time : 2019-07-09 14:40
 * @Author : zhuangjingpeng
 * @File : options
 * @Desc : file function description
 */
package pdmqd

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PDMQDConfig struct {
	ID int64

	LogPath string

	TCPAddress               string        `flag:"tcp-address"`
	HTTPAddress              string        `flag:"http-address"`
	HTTPClientConnectTimeout time.Duration `flag:"http-client-connect-timeout" cfg:"http_client_connect_timeout"`
	HTTPClientRequestTimeout time.Duration `flag:"http-client-request-timeout" cfg:"http_client_request_timeout"`

	MsgChanSize int64
	MsgMaxSize  int

	MaxChannelConsumers int

	MsgTimeout time.Duration
}

func InitConfig() *PDMQDConfig {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	current := getCurrentPath()
	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)

	initConf := &PDMQDConfig{
		ID: defaultID,

		LogPath:                  current + "/../../../internal/log/log.xml",
		TCPAddress:               "0.0.0.0:9400",
		HTTPAddress:              "0.0.0.0:9401",
		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 5 * time.Second,

		MsgChanSize: 9999,
		MsgMaxSize:  1024 * 1024 * 2,  //最大消息体为2Mb
		MsgTimeout:  60 * time.Second, //消息超时时间为1分钟
	}
	return initConf
}
func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
