package message

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type Shake struct {
	Version         uint32
	Capabilities    uint32
	TotalDifficulty uint64
	UserAgent       string
	Genesis         string
}

func NewShake(playload []byte) (Shake, error) {
	var shake Shake
	r := bytes.NewReader(playload)

	var version uint32
	err := binary.Read(r, binary.BigEndian, &version)
	if err != nil {
		return shake, err
	}

	var capabilities uint32
	err = binary.Read(r, binary.BigEndian, &capabilities)
	if err != nil {
		return shake, err
	}

	var diffculty uint64
	err = binary.Read(r, binary.BigEndian, &diffculty)
	if err != nil {
		return shake, err
	}

	var length uint64
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return shake, err
	}

	b := make([]byte, length, length)
	err = binary.Read(r, binary.BigEndian, &b)
	if err != nil {
		return shake, err
	}
	ua := string(b)

	hash := make([]byte, 32, 32)
	err = binary.Read(r, binary.BigEndian, &hash)
	if err != nil {
		return shake, err
	}
	h := hex.EncodeToString(hash)

	shake = Shake{
		Version:         version,
		Capabilities:    capabilities,
		TotalDifficulty: diffculty,
		UserAgent:       ua,
		Genesis:         h,
	}
	return shake, nil
}
