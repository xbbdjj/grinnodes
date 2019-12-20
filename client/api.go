package client

import (
	"encoding/json"
	"fmt"
	"github.com/xbbdjj/grinnodes/storage"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type capabilities struct {
	Bits uint8
}

type peerConnected struct {
	Capabilities    capabilities `json:"capabilities"`
	Addr            string       `json:"addr"`
	UserAgent       string       `json:"user_agent"`
	Version         uint32       `json:"version"`
	Height          uint64       `json:"height"`
	TotalDifficulty uint64       `json:"total_difficulty"`
}

func Start() {
	for {
		time.Sleep(time.Minute)
		peers, err := storage.PublicIP()
		if err != nil {
			continue
		}
		for _, ip := range peers {
			url := fmt.Sprintf("http://%s:3413/v1/peers/connected", ip)
			go connected(url)
			time.Sleep(time.Minute)
		}
	}
}

func connected(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	var peers []peerConnected
	json.NewDecoder(resp.Body).Decode(&peers)
	defer resp.Body.Close()

	for _, p := range peers {
		addr := strings.Split(p.Addr, ":")
		if len(addr) != 2 {
			continue
		}
		port, err := strconv.Atoi(addr[1])
		if err != nil {
			continue
		}
		v := storage.PeerConnected{
			IP:              addr[0],
			Port:            port,
			UserAgent:       p.UserAgent,
			Version:         p.Version,
			Capabilities:    p.Capabilities.Bits,
			Height:          p.Height,
			TotalDifficulty: p.TotalDifficulty,
		}
		storage.ClientPeersConnected(v)
	}
}
