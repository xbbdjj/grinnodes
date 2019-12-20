package client

import (
	"encoding/json"
	"fmt"
	"github.com/xbbdjj/grinnodes/config"
	"github.com/xbbdjj/grinnodes/log"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type status struct {
	UserAgent       string
	Version         uint32
	Height          uint64
	TotalDifficulty uint64
}

type tip struct {
	Height          uint64 `json:"height"`
	TotalDifficulty uint64 `json:"total_difficulty"`
}
type v1status struct {
	Version   uint32 `json:"protocol_version"`
	UserAgent string `json:"user_agent"`
	Tip       tip    `json:"tip"`
}

var Status status

func Sync() {
	for {
		time.Sleep(time.Minute)
		var s v1status
		url := fmt.Sprintf("http://%s/v1/status", config.Conf.ClientHttpAddr)
		resp, err := http.Get(url)
		if err != nil {
			log.Logger.Error("api status", zap.String("addr", config.Conf.ClientHttpAddr), zap.Error(err))
			continue
		}

		json.NewDecoder(resp.Body).Decode(&s)
		resp.Body.Close()

		Status.UserAgent = s.UserAgent
		Status.Version = s.Version
		Status.Height = s.Tip.Height
		Status.TotalDifficulty = s.Tip.TotalDifficulty
	}
}
