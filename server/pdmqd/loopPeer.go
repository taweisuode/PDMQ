/**
 * @Time : 2019-08-01 17:15
 * @Author : zhuangjingpeng
 * @File : loopPeer
 * @Desc : file function description
 */
package pdmqd

import (
	pdmq "GO_PDMQ"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	stateInit         = iota
	stateDisconnected = 1
	stateConnected    = 2
	stateSubscribed   = 3
	stateClosing      = 4
)

type loopPeer struct {
	addr            string
	conn            net.Conn
	state           int32
	connectCallback func(*loopPeer)
	MsgMaxSize      int
	Info            peerInfo
}

// peerInfo contains metadata for a lookupPeer instance (and is JSON marshalable)
type peerInfo struct {
	TCPPort          int    `json:"tcp_port"`
	HTTPPort         int    `json:"http_port"`
	Version          string `json:"version"`
	BroadcastAddress string `json:"broadcast_address"`
}

func newLookupPeer(addr string, msgMaxSize int, connectCallback func(*loopPeer)) *loopPeer {
	return &loopPeer{
		addr:            addr,
		state:           stateDisconnected,
		MsgMaxSize:      msgMaxSize,
		connectCallback: connectCallback,
	}
}

// Connect will Dial the specified address, with timeouts
func (lp *loopPeer) Connect() error {
	conn, err := net.DialTimeout("tcp", lp.addr, time.Second)
	if err != nil {
		return err
	}
	lp.conn = conn
	return nil
}

// String returns the specified address
func (lp *loopPeer) String() string {
	return lp.addr
}

// Read implements the io.Reader interface, adding deadlines
func (lp *loopPeer) Read(data []byte) (int, error) {
	lp.conn.SetReadDeadline(time.Now().Add(time.Second))
	return lp.conn.Read(data)
}

// Write implements the io.Writer interface, adding deadlines
func (lp *loopPeer) Write(data []byte) (int, error) {
	lp.conn.SetWriteDeadline(time.Now().Add(time.Second))
	return lp.conn.Write(data)
}

// Close implements the io.Closer interface
func (lp *loopPeer) Close() error {
	lp.state = stateDisconnected
	if lp.conn != nil {
		return lp.conn.Close()
	}
	return nil
}

//loopPeer 执行各个命令
func (lp *loopPeer) Command(cmd *pdmq.Command) ([]byte, error) {
	initialState := lp.state
	if lp.state != stateConnected {
		err := lp.Connect()
		if err != nil {
			return nil, err
		}
		lp.state = stateConnected
		_, err = lp.Write([]byte("V1"))
		if err != nil {
			lp.Close()
			return nil, err
		}
		if initialState == stateDisconnected {
			lp.connectCallback(lp)
		}
		if lp.state != stateConnected {
			return nil, fmt.Errorf("lookupPeer connectCallback() failed")
		}
	}
	if cmd == nil {
		return nil, nil
	}
	n, err := cmd.WriteTo(lp)
	fmt.Println(123, n, err)
	if err != nil {
		lp.Close()
		return nil, err
	}
	fmt.Println(lp, lp.MsgMaxSize)
	resp, err := readResponseBounded(lp, lp.MsgMaxSize)
	fmt.Printf("cmd is [%+v]\n", cmd.String())
	fmt.Printf("get loopPeer response is [%+v],error is [%+v]\n", resp, err)
	if err != nil {
		lp.Close()
		return nil, err
	}
	return resp, nil
}

//获取消息的响应
func readResponseBounded(r io.Reader, limit int) ([]byte, error) {
	var msgSize int32

	//message size
	err := binary.Read(r, binary.BigEndian, &msgSize)

	fmt.Println(123, err, msgSize)
	if err != nil {
		return nil, err
	}

	if int(msgSize) > limit {
		return nil, fmt.Errorf("response body size (%d) is greater than limit (%d)", msgSize, limit)
	}
	// message binary data
	buf := make([]byte, 100)
	_, err = io.ReadFull(r, buf)
	fmt.Printf("read debug r is [%+v],buf is [%+v],err is [%+v]\n", r, string(buf), err)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
