package p2p

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/xbbdjj/grinnodes/client"
	"github.com/xbbdjj/grinnodes/log"
	"github.com/xbbdjj/grinnodes/p2p/message"
	"github.com/xbbdjj/grinnodes/storage"
	"go.uber.org/zap"
)

type node struct {
	addr     net.TCPAddr
	lastSeen int64
}

var peerAll sync.Map
var peerConnected sync.Map

//Start ... start peer connect
func Start() {
	res, err := storage.P2PToConnect()
	if err == nil {
		for _, p := range res {
			n := node{addr: p, lastSeen: time.Now().Unix()}
			peerAll.Store(fmt.Sprintf("%s:%d", n.addr.IP, n.addr.Port), n)
		}
	}

	for {
		peerAll.Range(func(key, value interface{}) bool {
			time.Sleep(time.Second)
			n, ok := value.(node)
			if ok {
				if n.lastSeen < time.Now().Unix()-86400 {
					peerAll.Delete(fmt.Sprintf("%s:%d", n.addr.IP, n.addr.Port))
				} else {
					go work(n.addr)
				}
			}
			return true
		})
	}
}

//one goroutine to handle one p2p conenct
func work(addr net.TCPAddr) {
	_, ok := peerConnected.LoadOrStore(addr.IP.String(), time.Now().Unix)
	if ok {
		return
	}
	defer func() {
		peerConnected.Delete(addr.IP.String())
	}()

	conn, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		//msg := err.Error()
		//"operation timed out"
		//"connection refused"
		//"no route to host"
		//"network is unreachable"
		//"can't assign requested address"
		return
	}
	defer conn.Close()

	storage.P2PConnected(addr)

	log.Logger.Info("p2p connected", zap.String("ip", addr.IP.String()))

	msg := message.NewMainnetMessage(message.MsgTypeHand)
	hand := message.Hand{
		Version:         client.Status.Version,
		Capabilities:    15,
		Nonce:           uint64(time.Now().Unix()),
		TotalDifficulty: client.Status.TotalDifficulty,
		UserAgent:       client.Status.UserAgent,
		Genesis:         "40adad0aec27797b48840aa9e00472015c21baea118ce7a2ff1a82c0f8f5bf82",
	}
	msg.Payload = hand.Payload()

	_, err = conn.Write(msg.Bytes())
	if err != nil {
		log.Logger.Error("p2p write hand message", zap.String("ip", addr.IP.String()), zap.Error(err))
		return
	}

	tPing := time.NewTicker(time.Second * 30)
	defer tPing.Stop()

	tPeer := time.NewTicker(time.Minute * 10)
	defer tPeer.Stop()

	for {
		time.Sleep(time.Second * 10)
		select {
		case _ = <-tPing.C:
			ping := message.Ping{
				TotalDifficulty: client.Status.TotalDifficulty,
				Height:          client.Status.Height,
			}
			m := message.NewMainnetMessage(message.MsgTypePing)
			m.Payload = ping.Playload()

			_, err = conn.Write(m.Bytes())
			if err != nil {
				log.Logger.Error("p2p write ping message", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}
		case _ = <-tPeer.C:
			mask := [3]uint32{0x01, 0x02, 0x04}
			r := rand.New(rand.NewSource(time.Now().Unix()))
			get := message.GetPeerAddrs{
				Capabilities: mask[r.Intn(len(mask))],
			}
			msg := message.NewMainnetMessage(message.MsgTypeGetPeerAddrs)
			msg.Payload = get.Playload()

			_, err := conn.Write(msg.Bytes())
			if err != nil {
				log.Logger.Error("p2p write getPeerAddrs message", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}
		default:
			var magic1 uint8
			var magic2 uint8
			err := binary.Read(conn, binary.BigEndian, &magic1)
			if err != nil {
				log.Logger.Error("p2p read magic1", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}
			if magic1 != 97 {
				continue
			}
			err = binary.Read(conn, binary.BigEndian, &magic2)
			if err != nil {
				log.Logger.Error("p2p read magic2", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}
			if magic2 != 61 {
				continue
			}

			var msg uint8
			err = binary.Read(conn, binary.BigEndian, &msg)
			if err != nil {
				log.Logger.Error("p2p read msgType", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}

			var length uint64
			err = binary.Read(conn, binary.BigEndian, &length)
			if err != nil {
				log.Logger.Error("p2p read length", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}

			b := make([]byte, length, length)
			err = binary.Read(conn, binary.BigEndian, &b)
			if err != nil {
				log.Logger.Error("p2p read palyload", zap.String("ip", addr.IP.String()), zap.Error(err))
				return
			}

			storage.P2PConnecting(addr)

			if msg == 2 {
				shake, err := message.NewShake(b)
				// fmt.Printf("%#v\n", shake)

				if err != nil {
					log.Logger.Error("p2p decode shake", zap.String("ip", addr.IP.String()), zap.Error(err))
				} else {
					log.Logger.Debug("p2p receive shake", zap.String("ip", addr.IP.String()))
					storage.ReceiveShake(addr, shake)
				}
			} else if msg == 4 {
				pong, err := message.NewPong(b)
				// fmt.Printf("%#v\n", pong)

				if err != nil {
					log.Logger.Error("p2p decode pong", zap.String("ip", addr.IP.String()), zap.Error(err))
				} else {
					log.Logger.Debug("p2p receive pong", zap.String("ip", addr.IP.String()))
					storage.ReceivePong(addr, pong)
				}
			} else if msg == 6 {
				peers, err := message.NewPeers(b)
				// fmt.Printf("%#v\n", p)

				if err != nil {
					log.Logger.Error("p2p decode peers", zap.String("ip", addr.IP.String()), zap.Error(err))
				} else {
					log.Logger.Debug("p2p receive peers", zap.String("ip", addr.IP.String()))

					for _, p := range peers {
						if storage.IsPublicIP(p.IP) && p.Port > 0 {
							n := node{addr: p, lastSeen: time.Now().Unix()}
							peerAll.Store(fmt.Sprintf("%s:%d", p.IP, p.Port), n)
							storage.AddPeer(p)
						}
					}
				}
			} else if msg == 8 {
				header, err := message.NewHeader(b)
				// fmt.Printf("%#v\n", header)

				if err != nil {
					log.Logger.Error("p2p decode header", zap.String("ip", addr.IP.String()), zap.Error(err))
				} else {
					log.Logger.Debug("p2p receive header", zap.String("ip", addr.IP.String()))
					storage.ReceiveHeader(addr, header)
				}
			}
		}
	}
}
