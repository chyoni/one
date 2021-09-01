package p2p

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chiwon99881/one/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	remotePort := r.URL.Query().Get("remotePort")
	remoteAddr := r.RemoteAddr
	addrSlice := strings.Split(remoteAddr, ":")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		if addrSlice[0] == "" || remotePort == "" {
			return false
		}
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		utils.HandleErr(err)
	}
	fmt.Printf("from %s:%s want to upgrade\n", addrSlice[0], remotePort)
	initPeer(conn, addrSlice[0], remotePort)
}

func ConnectPeer(addr, port string, remotePort int) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?remotePort=%d", addr, port, remotePort), nil)
	if err != nil {
		utils.HandleErr(err)
	}
	fmt.Printf("Request to %s:%s for upgrade\n", addr, port)
	peer := initPeer(conn, addr, port)
	fmt.Printf("Sending my all blocks to %s:%s", addr, port)
	peer.sendNewestBlock()
}
