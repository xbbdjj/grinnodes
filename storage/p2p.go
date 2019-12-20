package storage

import (
	"github.com/xbbdjj/grinnodes/p2p/message"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//AllPeersToConnect ... get all peer to connect
func PeersToConnect() ([]net.TCPAddr, error) {
	ts := time.Now().Unix() - 3600*3
	peers := []net.TCPAddr{}
	rows, err := db.Query("SELECT ip, port FROM peer "+
		"WHERE p2p_first_connected > 0 OR p2p_failed < ? "+
		"ORDER BY p2p_last_connected DESC ",
		ts,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		var port int
		err := rows.Scan(&ip, &port)
		if err != nil {
			return nil, err
		}
		addr := net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		}
		peers = append(peers, addr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return peers, nil
}

//P2PConnected ... update connected time when receive p2p message
func P2PConnected(addr net.TCPAddr) error {
	ts := time.Now().Unix()
	if _, err := db.Exec(
		"UPDATE peer SET p2p_first_connected = ? WHERE ip = ?  AND p2p_first_connected = 0",
		ts,
		addr.IP.String(),
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		"UPDATE peer SET p2p_last_connected = ? WHERE ip = ? ",
		ts,
		addr.IP.String(),
	); err != nil {
		return err
	}
	return nil
}

//P2PFailed ... update fail time when cannot to connect
func P2PFailed(addr net.TCPAddr) error {
	if _, err := db.Exec(
		"UPDATE peer SET p2p_failed = ? WHERE ip = ? ",
		time.Now().Unix(),
		addr.IP.String(),
	); err != nil {
		return err
	}
	return nil
}

//ReceivePong ... update peer info by p2p pong message
func ReceivePong(addr net.TCPAddr, pong message.Pong) error {
	if _, err := db.Exec(
		"UPDATE peer SET node_height = ? WHERE ip = ? AND node_height < ?",
		pong.Height,
		addr.IP.String(),
		pong.Height,
	); err != nil {
		return err
	}
	if _, err := db.Exec(
		"UPDATE peer SET node_total_difficulty = ? WHERE ip = ? AND node_total_difficulty < ?",
		pong.TotalDifficulty,
		addr.IP.String(),
		pong.TotalDifficulty,
	); err != nil {
		return err
	}
	return nil
}

//ReceiveShake ... update peer info by p2p shake message
func ReceiveShake(addr net.TCPAddr, shake message.Shake) error {
	if _, err := db.Exec(
		"UPDATE peer SET node_user_agent = ?, node_protocol_version = ?, node_capabilities = ? WHERE ip = ? ",
		shake.UserAgent,
		shake.Version,
		shake.Capabilities,
		addr.IP.String(),
	); err != nil {
		return err
	}
	return nil
}

//ReceivePeers ... update peer info by p2p PeerAddrs message
func ReceivePeers(peers []net.TCPAddr) error {
	for _, p := range peers {
		if !IsPublicIP(p.IP) || p.Port == 0 {
			continue
		}
		db.Exec(
			"INSERT INTO peer (ip, port, p2p_first_seen) VALUES (?, ?, ?)",
			p.IP.String(),
			p.Port,
			time.Now().Unix(),
		)

		if _, err := db.Exec(
			"UPDATE peer SET port = ?, p2p_last_seen = ? WHERE ip = ?",
			p.Port,
			time.Now().Unix(),
			p.IP.String(),
		); err != nil {
			return err
		}
	}
	return nil
}

//ReceiveHeader ... update peer info by Header message
func ReceiveHeader(addr net.TCPAddr, header message.Header) error {
	if _, err := db.Exec(
		"UPDATE peer SET node_height = ?  WHERE ip = ? AND node_height < ? ",
		header.Height,
		addr.IP.String(),
		header.Height,
	); err != nil {
		return err
	}
	return nil
}
