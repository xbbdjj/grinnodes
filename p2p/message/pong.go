package message

import (
	"bytes"
	"encoding/binary"
)

type Pong struct {
	TotalDifficulty uint64
	Height          uint64
}

func NewPong(playload []byte) (Pong, error) {
	var pong Pong
	r := bytes.NewReader(playload)

	var diffculty uint64
	err := binary.Read(r, binary.BigEndian, &diffculty)
	if err != nil {
		return pong, err
	}

	var height uint64
	err = binary.Read(r, binary.BigEndian, &height)
	if err != nil {
		return pong, err
	}

	pong = Pong{TotalDifficulty: diffculty, Height: height}
	return pong, nil
}
