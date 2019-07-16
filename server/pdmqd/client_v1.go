/**
 * @Time : 2019-07-15 14:15
 * @Author : zhuangjingpeng
 * @File : client_v1
 * @Desc : file function description
 */
package pdmqd

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const defaultBufferSize = 16 * 1024

type clientV1 struct {
	ReadyCount    int64
	InFlightCount int64
	MessageCount  int64
	FinishCount   int64
	RequeueCount  int64

	ID  int64
	ctx *context
	net.Conn

	sync.RWMutex

	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer

	MsgTimeout time.Duration

	State          int32
	ConnectTime    time.Time
	Channel        *Channel
	SubEventChan   chan *Channel
	ReadyStateChan chan int
	ExitChan       chan int

	lenBuf   [2]byte
	lenSlice []byte

	ClientHost string
	ClientPort string

	pubCount map[string]int64
}

func newClientV1(id int64, conn net.Conn, ctx *context) *clientV1 {
	host, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	client := &clientV1{
		ID:     id,
		ctx:    ctx,
		Conn:   conn,
		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),

		MsgTimeout: ctx.pdmqd.config.MsgTimeout,

		ReadyStateChan: make(chan int, 1),
		SubEventChan:   make(chan *Channel, 1),
		ExitChan:       make(chan int),
		ClientHost:     host,
		ClientPort:     port,

		pubCount: make(map[string]int64),
	}
	client.lenSlice = client.lenBuf[:]
	return client
}

func (c *clientV1) IsProducer() bool {
	c.RLock()
	count := len(c.pubCount) > 0
	c.RUnlock()
	return count
}

func (c *clientV1) Pause() {
	c.tryUpdateReadyState()
}

func (c *clientV1) UnPause() {
	c.tryUpdateReadyState()
}
func (c *clientV1) Close() error {
	return nil
}

func (c *clientV1) TimedOutMessage() {
	atomic.AddInt64(&c.InFlightCount, -1)
	c.tryUpdateReadyState()
}

func (c *clientV1) RequeuedMessage() {
	atomic.AddInt64(&c.RequeueCount, 1)
	atomic.AddInt64(&c.InFlightCount, -1)
	c.tryUpdateReadyState()
}
func (c *clientV1) Empty() {
	atomic.StoreInt64(&c.InFlightCount, 0)
	c.tryUpdateReadyState()
}

func (c *clientV1) tryUpdateReadyState() {
	// you can always *try* to write to ReadyStateChan because in the cases
	// where you cannot the message pump loop would have iterated anyway.
	// the atomic integer operations guarantee correctness of the value.
	select {
	case c.ReadyStateChan <- 1:
	default:
	}
}
