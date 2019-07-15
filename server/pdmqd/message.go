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
)

type MessageID [MsgIDLength]byte

type Message struct {
	ID        MessageID
	Body      []byte
	Timestamp int64
	Attempts  uint16

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
		ID:        id,
		Body:      body,
		Timestamp: time.Now().UnixNano(),
	}
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
	var buf [10]byte
	var total int64

	binary.BigEndian.PutUint64(buf[:8], uint64(m.Timestamp))
	binary.BigEndian.PutUint16(buf[8:10], uint16(m.Attempts))

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
