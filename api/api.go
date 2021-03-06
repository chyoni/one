package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chiwon99881/one/blockchain"
	"github.com/chiwon99881/one/p2p"
	"github.com/chiwon99881/one/utils"
	"github.com/gorilla/mux"
)

func (u URL) MarshalText() (text []byte, err error) {
	marshalURL := fmt.Sprintf("http://localhost:%d%s", port, u)
	return []byte(marshalURL), nil
}

type URL string

type errResponse struct {
	ErrMessage string `json:"errMessage"`
}

type addTransactionPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type BalanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addPeerPayload struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

type urlDescription struct {
	URL         URL    `json:"url"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Payload     string `json:"payload,omitempty"`
}

var port int

func home(rw http.ResponseWriter, r *http.Request) {
	url := []urlDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "GET",
			Description: "See all blocks in one coin's blockchain",
		},
		{
			URL:         URL("/block"),
			Method:      "POST",
			Description: "Add a block to one coin's blockchain",
			Payload:     "data:string",
		},
		{
			URL:         URL("/block/{block_hash}"),
			Method:      "GET",
			Description: "See a block in one coin's blockchain",
		},
		{
			URL:         URL("/blockchain"),
			Method:      "GET",
			Description: "See coin's blockchain status",
		},
		{
			URL:         URL("/mempool"),
			Method:      "GET",
			Description: "See all transactions in mempool",
		},
		{
			URL:         URL("/transaction/add"),
			Method:      "POST",
			Description: "Add a transaction to mempool",
		},
		{
			URL:         URL("/balance/{address}?total=true"),
			Method:      "GET",
			Description: "See who's balance. If you give total querystring, total balance return.",
		},
	}
	marshalToJSON, err := json.Marshal(url)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(rw, "%s", marshalToJSON)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	blocks := blockchain.Blocks(blockchain.BlockChain())
	resToJSON, err := json.Marshal(blocks)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(rw, "%s", resToJSON)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	paramsMap := mux.Vars(r)
	hash := paramsMap["block_hash"]
	block := blockchain.FindBlock(hash)
	err := json.NewEncoder(rw).Encode(block)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func addBlock(rw http.ResponseWriter, r *http.Request) {
	newBlock := blockchain.AddBlock(blockchain.BlockChain())
	for _, peer := range p2p.Peers.P {
		peer.NewBlock(newBlock)
	}
	rw.WriteHeader(http.StatusCreated)
}

func chainStatus(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	err := blockchain.Status(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	err := json.NewEncoder(rw).Encode(blockchain.Mempool().Txs)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address, exist := vars["address"]
	if !exist {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: "no address has been given."})
		return
	}
	switch r.URL.Query().Get("total") {
	case "true":
		balance := blockchain.GetBalanceByAddress(address)
		rw.WriteHeader(http.StatusOK)
		res, err := json.Marshal(BalanceResponse{Address: address, Balance: balance})
		utils.HandleErr(err)
		fmt.Fprintf(rw, "%s", res)
	case "":
		txOuts := blockchain.GetUTxOutsByAddress(address)
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(txOuts)
	}
}

func addTx(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	addTxPayload := &addTransactionPayload{}
	err := json.NewDecoder(r.Body).Decode(addTxPayload)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	tx, err := blockchain.Mempool().AddTx(addTxPayload.To, addTxPayload.Amount)
	for _, peer := range p2p.Peers.P {
		peer.NewTx(tx)
	}
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func addPeer(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var addPeerPayload addPeerPayload
	err := json.NewDecoder(r.Body).Decode(&addPeerPayload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusOK)
	p2p.ConnectPeer(addPeerPayload.Addr, addPeerPayload.Port, port, true)
}

func myPeer(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	err := p2p.AllPeers(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
}

func JSONHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}

func Start(aPort int) {
	port = aPort
	fmt.Printf("Server listening on http://localhost:%d\n", port)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(JSONHeaderMiddleware)
	router.Use(loggerMiddleware)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET")
	router.HandleFunc("/block/{block_hash}", block).Methods("GET")
	router.HandleFunc("/block", addBlock).Methods("POST")
	router.HandleFunc("/blockchain", chainStatus).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/transaction/add", addTx).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peer", addPeer).Methods("POST")
	router.HandleFunc("/myPeer", myPeer).Methods("GET")
	utils.HandleErr(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
