package main

import (
	"fmt"

	"github.com/chiwon99881/one/wallet"
)

func main() {
	// cli.Start()
	// defer db.Close()
	w := wallet.Wallet()
	fmt.Println(w.Address)
}
