package message

import (
	"bytes"
	"encoding/binary"
	"net"
)

func NewPeers(playload []byte) ([]net.TCPAddr, error) {
	r := bytes.NewReader(playload)

	var size uint32
	err := binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	peers := []net.TCPAddr{}

	var i uint32 = 0
	for i < size {
		i++
		var family uint8
		err := binary.Read(r, binary.BigEndian, &family)
		if err != nil {
			return nil, err
		}
		if family == 0 {
			ipv4 := make([]uint8, 4, 4)
			err := binary.Read(r, binary.BigEndian, &ipv4)
			if err != nil {
				return nil, err
			}

			var port uint16
			err = binary.Read(r, binary.BigEndian, &port)
			if err != nil {
				return nil, err
			}

			addr := net.TCPAddr{
				IP:   net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3]),
				Port: int(port),
			}
			peers = append(peers, addr)
		}
		if family == 1 {
			ipv6 := make([]uint8, 16, 16)
			err := binary.Read(r, binary.BigEndian, &ipv6)
			if err != nil {
				return nil, err
			}
			var port uint16
			err = binary.Read(r, binary.BigEndian, &port)
			if err != nil {
				return nil, err
			}
		}
	}
	return peers, nil
}
