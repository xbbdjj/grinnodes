package message

import (
	"bytes"
	"encoding/binary"
)

type MsgType uint8

const (
	MsgTypeError MsgType = iota
	MsgTypeHand
	MsgTypeShake
	MsgTypePing
	MsgTypePong
	MsgTypeGetPeerAddrs
	MsgTypePeerAddrs
	MsgTypeGetHeaders
	MsgTypeHeader
	MsgTypeHeaders
	MsgTypeGetBlock
	MsgTypeBlock
	MsgTypeGetCompactBlock
	MsgTypeCompactBlock
	MsgTypeStemTransaction
	MsgTypeTransaction
	MsgTypeTxHashSetRequest
	MsgTypeTxHashSetArchive
	MsgTypeBanReason
)

//Message ... p2p message struct
type Message struct {
	Magic   [2]byte
	Type    MsgType
	Length  uint64
	Payload []byte
}

func NewMainnetMessage(msg MsgType) Message {
	return Message{
		Magic: [2]byte{97, 61},
		Type:  msg,
	}
}

//Bytes ... encode message to bytes
func (msg *Message) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, msg.Magic)
	binary.Write(buf, binary.BigEndian, msg.Type)
	binary.Write(buf, binary.BigEndian, uint64(len(msg.Payload)))
	binary.Write(buf, binary.BigEndian, msg.Payload)
	return buf.Bytes()
}
