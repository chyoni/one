package p2p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	P map[string]*peer `json:"peers"`
	m sync.Mutex
}

type peer struct {
	conn *websocket.Conn
	Addr string `json:"addr"`
	Port string `json:"port"`
}

var Peers *peers = &peers{
	P: make(map[string]*peer),
}

func AllPeers(rw http.ResponseWriter) error {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	err := json.NewEncoder(rw).Encode(Peers.P)
	return err
}

func initPeer(conn *websocket.Conn, addr, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", addr, port)
	peer := &peer{
		conn: conn,
		Addr: addr,
		Port: port,
	}
	Peers.P[key] = peer
	return peer
}
