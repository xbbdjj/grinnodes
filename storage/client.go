package storage

import (
	"errors"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PeerConnected struct {
	IP              string
	Port            int
	Capabilities    uint8
	UserAgent       string
	Version         uint32
	Height          uint64
	TotalDifficulty uint64
}

//PublicIP ... peer have public ip can connected
func PublicIP() ([]string, error) {
	ips := []string{}
	rows, err := db.Query("SELECT ip FROM peer WHERE p2p_first_connected > 0 ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ips, nil
}

//ClientPeersConnected ... request node api /v1/peers/connected
func ClientPeersConnected(p PeerConnected) error {
	if !IsPublicIP(net.ParseIP(p.IP)) || p.Port == 0 {
		return errors.New("ip private or port zero")
	}
	ts := time.Now().Unix()
	db.Exec("INSERT INTO peer (ip, port, api_first_seen) VALUES (?, ?, ?) ", p.IP, p.Port, ts)
	if _, err := db.Exec(
		"UPDATE peer SET api_first_seen = ?  WHERE ip = ? AND api_first_seen = 0 ",
		ts,
		p.IP,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		"UPDATE peer SET api_last_seen = ? WHERE ip = ? ",
		ts,
		p.IP,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		"UPDATE peer SET port = ?, node_user_agent = ?, node_protocol_version = ?, "+
			"node_capabilities = ?, node_height = ?, node_total_difficulty = ? "+
			"WHERE ip = ? AND p2p_first_connected = 0 ",
		p.Port,
		p.UserAgent,
		p.Version,
		p.Capabilities,
		p.Height,
		p.TotalDifficulty,
		p.IP,
	); err != nil {
		return err
	}
	return nil
}
