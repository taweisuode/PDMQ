package pdmqd

import (
	"net"
)

type Connect interface {
	AcceptConnect(pdmqd *PDMQD, conn *net.TCPConn)
}
