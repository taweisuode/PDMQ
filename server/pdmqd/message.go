package pdmqd

import (
	"encoding/binary"
	"fmt"
	"github.com/timespacegroup/go-utils"
	"io"
	"time"
)

const (
	MsgIDLength       = 16
	minValidMsgLength = MsgIDLength + 8 + 2 // Timestamp + Attempts

	ProtocolCommonResponse  int16 = 1
	ProtocolErrorResponse   int16 = 2
	ProtocolMessageResponse int16 = 3
)

type MessageID [MsgIDLength]byte

type Message struct {
	ID           [MsgIDLength]byte
	Body         []byte
	Timestamp    int64
	Attempts     uint16
	ProtocolType int16

	// for in-flight handling
	deliveryTS time.Time
	clientID   int64
	pri        int64
	index      int
	deferred   time.Duration
}

func (msg *Message) CreateMessageId() MessageID {
	var buf MessageID
	guid := []byte(tsgutils.GUID())
	for index, value := range guid {
		if index < MsgIDLength {
			buf[index] = value
		}
	}
	return buf
}
func CreateMessage(id MessageID, body []byte) *Message {
	return &Message{
		ID:           id,
		Body:         body,
		ProtocolType: ProtocolCommonResponse, //默认消息类型是普通消息
		Timestamp:    time.Now().UnixNano(),
	}
}

func decodeMessage(b []byte) (*Message, error) {
	var msg Message

	if len(b) < minValidMsgLength {
		return nil, fmt.Errorf("invalid message buffer size (%d)", len(b))
	}

	msg.Timestamp = int64(binary.BigEndian.Uint64(b[:8]))
	msg.Attempts = binary.BigEndian.Uint16(b[8:10])
	copy(msg.ID[:], b[10:10+MsgIDLength])
	msg.Body = b[10+MsgIDLength:]

	return &msg, nil
}

/**
 * @desc revert []byte into Message struct
 * @param (query []byte
 * @return (msg *Message)
 */
func RevertMessage(buf []byte) *Message {
	var msg *Message
	fmt.Println(time.Now().Nanosecond())

	msg.Timestamp = time.Now().UnixNano()
	return msg
}

func (m *Message) WriteTo(w io.Writer) (int64, error) {
	var buf [12]byte
	var total int64

	binary.BigEndian.PutUint64(buf[:8], uint64(m.Timestamp))
	binary.BigEndian.PutUint16(buf[8:10], uint16(m.Attempts))      //消息叠加次数可过百
	binary.BigEndian.PutUint16(buf[10:12], uint16(m.ProtocolType)) //这里将消息类型也加入其中

	n, err := w.Write(buf[:])
	total += int64(n)
	if err != nil {
		return total, err
	}
	n, err = w.Write(m.ID[:])
	total += int64(n)
	if err != nil {
		return total, err
	}

	n, err = w.Write(m.Body)
	total += int64(n)
	if err != nil {
		return total, err
	}

	return total, nil
}
