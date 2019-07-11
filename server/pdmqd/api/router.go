/**
 * @Time : 2019-07-09 17:24
 * @Author : zhuangjingpeng
 * @File : router
 * @Desc : file function description
 */
package api

import (
	"fmt"
	"net/http"
	"time"
)

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*HTTPContext)

// HandlersChain defines a HandlerFunc array.

type HTTPContext http.HandlerFunc

func Ping(c *HTTPContext) {
	//fmt.Println("hello")
	fmt.Fprintf(c, "ok") //这个写入到w的是输出到客户端的
	//w.Write([]byte("bye bye ,this is v2 httpServer"))
}

func SendResult(w http.ResponseWriter, errCode int, msg string, data interface{}) {
	result := map[string]interface{}{
		"code":        errCode,
		"message":     msg,
		"currentTime": time.Now().Unix(),
		"data":        data,
	}
	c.JSON(200, result)
}
