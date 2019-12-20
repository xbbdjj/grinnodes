package message

import (
	"bytes"
	"encoding/binary"
)

type GetPeerAddrs struct {
	Capabilities uint32
}

func (g *GetPeerAddrs) Playload() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, g.Capabilities)
	return buf.Bytes()
}
