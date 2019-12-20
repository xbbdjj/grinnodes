package message

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
)

type Hand struct {
	Version         uint32
	Capabilities    uint32
	Nonce           uint64
	TotalDifficulty uint64
	SenderAddress   net.TCPAddr
	ReceiverAddress net.TCPAddr
	UserAgent       string
	Genesis         string
}

func (h *Hand) Payload() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, h.Version)
	binary.Write(buf, binary.BigEndian, h.Capabilities)
	binary.Write(buf, binary.BigEndian, h.Nonce)
	binary.Write(buf, binary.BigEndian, h.TotalDifficulty)

	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint16(0))

	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint16(0))

	binary.Write(buf, binary.BigEndian, uint64(len(h.UserAgent)))
	buf.WriteString(h.UserAgent)

	genesis, err := hex.DecodeString(h.Genesis)
	if err == nil {
		binary.Write(buf, binary.BigEndian, genesis)
	}
	return buf.Bytes()
}
