package chart

import (
	"github.com/xbbdjj/grinnodes/storage"
	"time"
)

func Start() {
	year, month, day := time.Now().Date()
	hour := time.Now().Hour() + 2
	s := time.Date(year, month, day, hour, 0, 0, 0, time.Now().Location())
	time.Sleep(time.Until(s))
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			storage.ClearOldPeer()
			total, err := storage.NodeTotal()
			if err != nil {
				continue
			}

			publictotal, e := storage.NodePublicCount()
			if e != nil {
				continue
			}
			storage.AddChart(total, publictotal)
		}
	}
}
