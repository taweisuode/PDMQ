package pdmqd

import (
	"fmt"
	"github.com/timespacegroup/go-utils"
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
