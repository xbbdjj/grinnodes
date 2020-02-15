package storage

import (
	"github.com/xbbdjj/grinnodes/p2p/message"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//P2PToConnect ... get public peer to connect
func P2PToConnect() ([]net.TCPAddr, error) {
	peers := []net.TCPAddr{}
	rows, err := db.Query("SELECT ip, port FROM peer where p2p_last_connected > 0 ")
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

//P2PConnected ... update connected info when connected
func P2PConnected(addr net.TCPAddr) {
	db.Exec(
		"UPDATE peer SET port = ? WHERE ip = ? ",
		addr.Port,
		addr.IP.String(),
	)

	ts := time.Now().Unix()
	db.Exec(
		"UPDATE peer SET p2p_first_connected = ? WHERE ip = ?  AND p2p_first_connected = 0",
		ts,
		addr.IP.String(),
	)
}

//P2PConnecting ... update connected info when receive p2p message
func P2PConnecting(addr net.TCPAddr) {
	ts := time.Now().Unix()
	db.Exec(
		"UPDATE peer SET p2p_last_connected = ? WHERE ip = ? ",
		ts,
		addr.IP.String(),
	)
}

//ReceivePong ... update peer info by p2p pong message
func ReceivePong(addr net.TCPAddr, pong message.Pong) {
	db.Exec(
		"UPDATE peer SET node_height = ? WHERE ip = ? AND node_height < ?",
		pong.Height,
		addr.IP.String(),
		pong.Height,
	)

	db.Exec(
		"UPDATE peer SET node_total_difficulty = ? WHERE ip = ? AND node_total_difficulty < ?",
		pong.TotalDifficulty,
		addr.IP.String(),
		pong.TotalDifficulty,
	)
}

//ReceiveShake ... update peer info by p2p shake message
func ReceiveShake(addr net.TCPAddr, shake message.Shake) {
	db.Exec(
		"UPDATE peer SET node_user_agent = ?, node_protocol_version = ?, node_capabilities = ? WHERE ip = ? ",
		shake.UserAgent,
		shake.Version,
		shake.Capabilities,
		addr.IP.String(),
	)
}

//AddPeer ... add one peer info by p2p PeerAddrs message
func AddPeer(peer net.TCPAddr) {
	if !IsPublicIP(peer.IP) || peer.Port == 0 {
		return
	}
	ts := time.Now().Unix()
	db.Exec(
		"INSERT INTO peer (ip, p2p_first_seen) VALUES (?, ?)",
		peer.IP.String(),
		ts,
	)
	db.Exec(
		"UPDATE peer SET p2p_last_seen = ?  WHERE ip = ? AND p2p_last_seen < ? ",
		ts,
		peer.IP.String(),
		ts,
	)
}

//ReceiveHeader ... update peer info by Header message
func ReceiveHeader(addr net.TCPAddr, header message.Header) {
	db.Exec(
		"UPDATE peer SET node_height = ?  WHERE ip = ? AND node_height < ? ",
		header.Height,
		addr.IP.String(),
		header.Height,
	)
}
