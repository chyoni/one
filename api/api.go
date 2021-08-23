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

type errResponse struct {
	ErrMessage string `json:"errMessage"`
}

type AddBlockPayload struct {
	Data string `json:"data"`
}

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
	defer r.Body.Close()
	bodyData := &AddBlockPayload{}
	err := json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", errResponse{ErrMessage: err.Error()})
		return
	}
	blockchain.AddBlock(blockchain.BlockChain(), bodyData.Data)
	rw.WriteHeader(http.StatusCreated)
}

func JSONHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Start(aPort string) {
	port = aPort
	fmt.Printf("Server listening on http://localhost:%s", port)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(JSONHeaderMiddleware)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET")
	router.HandleFunc("/block/{block_hash}", block).Methods("GET")
	router.HandleFunc("/block", addBlock).Methods("POST")
	utils.HandleErr(http.ListenAndServe(":4000", router))
}
