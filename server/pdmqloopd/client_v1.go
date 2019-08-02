/**
 * @Time : 2019-07-15 14:15
 * @Author : zhuangjingpeng
 * @File : client_v1
 * @Desc : file function description
 */
package pdmqloopd

import (
	"net"
)

type clientV1 struct {
	net.Conn
	peerInfo *PeerInfo
}

func NewClientV1(conn net.Conn) *clientV1 {
	return &clientV1{
		Conn: conn,
	}
}
