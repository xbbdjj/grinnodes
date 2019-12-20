package message

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	Version   uint16
	Height    uint64
	Timestamp int64
}

//NewHeader ... from p2p message playload byte to decode header message
func NewHeader(playload []byte) (Header, error) {
	var header Header
	r := bytes.NewReader(playload)

	var version uint16
	err := binary.Read(r, binary.BigEndian, &version)
	if err != nil {
		return header, err
	}

	var height uint64
	err = binary.Read(r, binary.BigEndian, &height)
	if err != nil {
		return header, err
	}

	var timestamp int64
	err = binary.Read(r, binary.BigEndian, &timestamp)
	if err != nil {
		return header, err
	}

	header = Header{
		Version:   version,
		Height:    height,
		Timestamp: timestamp,
	}
	return header, nil
}
