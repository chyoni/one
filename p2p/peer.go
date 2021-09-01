package p2p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/chiwon99881/one/utils"
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
	Key   string `json:"key"`
	inbox chan []byte
}

type allPeersResponse struct {
	Peers []string `json:"peers"`
}

var Peers *peers = &peers{
	P: make(map[string]*peer),
}

func AllPeers(rw http.ResponseWriter) error {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	allPeersResponse := &allPeersResponse{}
	for key := range Peers.P {
		allPeersResponse.Peers = append(allPeersResponse.Peers, key)
	}
	err := json.NewEncoder(rw).Encode(allPeersResponse)
	return err
}

func (p *peer) read() {
	defer Peers.m.Unlock()
	for {
		message := &Message{}
		err := p.conn.ReadJSON(message)
		if err != nil {
			Peers.m.Lock()
			defer p.conn.Close()
			delete(Peers.P, p.Key)
			return
		}
		BroadcastMessage(message.MessageKind, message.Payload, p)
	}
}

func (p *peer) write() {
	defer Peers.m.Unlock()
	for {
		m, ok := <-p.inbox
		if !ok {
			Peers.m.Lock()
			defer p.conn.Close()
			delete(Peers.P, p.Key)
			return
		}
		message := &Message{}
		utils.FromBytes(message, m)
		err := p.conn.WriteJSON(message)
		if err != nil {
			Peers.m.Lock()
			defer p.conn.Close()
			delete(Peers.P, p.Key)
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
		Key:   key,
		inbox: make(chan []byte),
	}
	Peers.P[key] = peer
	go peer.read()
	go peer.write()
	return peer
}
