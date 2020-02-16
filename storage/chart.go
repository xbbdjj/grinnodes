package storage

import (
	"github.com/xbbdjj/grinnodes/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AddChart(total int, publicTotal int) error {
	if _, err := db.Exec(
		"INSERT INTO chart (time, peer_total, peer_public_total) VALUES (?, ?, ?)",
		time.Now().Unix(),
		total,
		publicTotal,
	); err != nil {
		return err
	}
	return nil
}

type chart struct {
	UnixTime    int64
	Total       int
	PublicTotal int
}

func GetChart() ([]chart, error) {
	arr := []chart{}
	rows, err := db.Query("SELECT time, peer_total, peer_public_total FROM chart")
	if err != nil {
		return arr, err
	}
	defer rows.Close()
	for rows.Next() {
		var unix int64
		var total int
		var public int
		err := rows.Scan(&unix, &total, &public)
		if err != nil {
			return arr, err
		}
		c := chart{UnixTime: unix * 1000, Total: total, PublicTotal: public}
		arr = append(arr, c)
	}
	err = rows.Err()
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func ClearOldPeer() {
	ts := time.Now().Unix() - config.NewConfig().PeerClearDuration
	db.Exec(
		"DELETE FROM peer WHERE p2p_last_connected < ? AND p2p_last_seen < ? ",
		ts,
		ts,
	)
}
