package api

import (
	"fmt"
	"net/http"

	"github.com/chiwon99881/one/utils"
)

var port string

func Start(aPort string) {
	port = aPort
	fmt.Printf("Server listening on http://localhost:%s", port)
	utils.HandleErr(http.ListenAndServe(":4000", nil))
}
