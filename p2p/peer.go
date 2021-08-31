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
	conn  *websocket.Conn
	Addr  string `json:"addr"`
	Port  string `json:"port"`
	inbox chan interface{}
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

func (p *peer) read() {
	for {
		message := &Message{}
		err := p.conn.ReadJSON(message)
		if err != nil {
			fmt.Println(err.Error())
			defer p.conn.Close()
			// delete peer
			return
		}
		//BroadcastMessage(message.MessageKind, message.Payload)
	}
}

func (p *peer) write() {
	for {
		m, ok := <-p.inbox
		if !ok {
			defer p.conn.Close()
			// delete peer
			return
		}
		err := p.conn.WriteJSON(m)
		if err != nil {
			defer p.conn.Close()
			// delete peer
			return
		}
	}
}

func initPeer(conn *websocket.Conn, addr, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", addr, port)
	peer := &peer{
		conn:  conn,
		Addr:  addr,
		Port:  port,
		inbox: make(chan interface{}),
	}
	Peers.P[key] = peer
	go peer.read()
	go peer.write()
	return peer
}
