/**
 * @Time : 2019-08-01 17:15
 * @Author : zhuangjingpeng
 * @File : loopPeer
 * @Desc : file function description
 */
package pdmqd

import (
	pdmq "GO_PDMQ"
	"fmt"
	"io"
	"io/ioutil"
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
	fmt.Printf("lp is [%+v],%+v,%+v\n", lp, lp.conn, lp.MsgMaxSize)
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
	/*var msgSize int32
	err := binary.Read(r, binary.BigEndian, &msgSize)
	fmt.Printf("read response length is [%+v],err is [%+v]\n", msgSize, err)
	if err != nil {
		return nil, err
	}*/
	result, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("ReadAll error: ", err.Error())
	}
	fmt.Println("result = ", string(result))
	/*buf := make([]byte, 1024)
	//接收pdmqd 的 ping 消息 以及 consumer 的注册消息
	for {
		//服务器端返回的数据写入空buf
		cnt, err := r.Read(buf)

		if err != nil {
			fmt.Printf("客户端读取数据失败 %s\n", err)
			break
		}

		return buf[0:cnt], nil
	}*/
	return result, nil
	/*	_, err := io.ReadFull(r, buf)
		fmt.Printf("read debug r is [%+v],buf is [%+v],err is [%+v]\n", r, string(buf), err)
		return buf, nil*/
}
