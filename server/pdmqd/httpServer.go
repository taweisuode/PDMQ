/**
 * @Time : 2019-07-09 15:53
 * @Author : zhuangjingpeng
 * @File : tcpServer
 * @Desc : file function description
 */
package pdmqd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
)

type httpServer struct {
	ctx    *context
	router http.Handler
}

func HTTPServer(listener net.Listener, handler http.Handler, proto string) error {
	server := &http.Server{
		Handler: handler,
	}
	err := server.Serve(listener)
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		return fmt.Errorf("http.Serve() error - %s", err)
	}
	fmt.Printf("%s: closing %s", proto, listener.Addr())
	return nil
}

//httpServer 实现了 http.Handler中的 ServeHTTP方法
func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

//用gin 作为http 服务框架
func newHTTPServer(ctx *context) *httpServer {
	ginApi := gin.New()
	gin.SetMode(gin.DebugMode)
	ginApi.Use(AddTraceId())
	server := &httpServer{
		ctx:    ctx,
		router: ginApi,
	}
	//探活
	ginApi.GET("/ping", server.Ping)

	//发布内容
	ginApi.POST("/pub", server.Pub)

	return server

}
