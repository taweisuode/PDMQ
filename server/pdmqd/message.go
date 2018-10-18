package pdmqd

import (
	"fmt"
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

/**
 * @desc revert []byte into Message struct
 * @param (query []byte
 * @return (msg *Message)
 */
func MakeMessage(buf []byte) *Message {
	var msg *Message
	fmt.Println(time.Nanosecond.Nanoseconds())

	msg.Timestamp = time.Nanosecond.Nanoseconds()
	return msg
}
