package message

import (
	"bytes"
	"encoding/binary"
)

type Ping struct {
	TotalDifficulty uint64
	Height          uint64
}

func (p *Ping) Playload() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.TotalDifficulty)
	binary.Write(buf, binary.BigEndian, p.Height)
	return buf.Bytes()
}
