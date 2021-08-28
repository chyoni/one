package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/chiwon99881/one/api"
)

func Start() {

	if len(os.Args) == 1 {
		fmt.Printf("Please give the 'port' flag like this '-port=4000'\n\n")
		os.Exit(2)
	}

	var port int

	flag.IntVar(&port, "port", 4000, "Your node start this port.")
	flag.Parse()

	if flag.Parsed() {
		api.Start(port)
	}
}
