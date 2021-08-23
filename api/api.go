package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chiwon99881/one/blockchain"
	"github.com/chiwon99881/one/utils"
	"github.com/gorilla/mux"
)

func (u URL) MarshalText() (text []byte, err error) {
	marshalURL := fmt.Sprintf("http://localhost:%s%s", port, u)
	return []byte(marshalURL), nil
}

type URL string

type urlDescription struct {
	URL         URL    `json:"url"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Payload     string `json:"payload,omitempty"`
}

var port string

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
	}
	marshalToJSON, err := json.Marshal(url)
	utils.HandleErr(err)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(rw, "%s", marshalToJSON)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Something's wrong...")
		return
	}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	blocks := blockchain.Blocks(blockchain.BlockChain())
	resToJSON, err := json.Marshal(blocks)
	utils.HandleErr(err)
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(rw, "%s", resToJSON)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, "Something wrong...")
		return
	}
}

func Start(aPort string) {
	port = aPort
	fmt.Printf("Server listening on http://localhost:%s", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET")
	utils.HandleErr(http.ListenAndServe(":4000", router))
}
