package main

import (
	"github.com/chiwon99881/one/cli"
	"github.com/chiwon99881/one/db"
)

func main() {
	cli.Start()
	defer db.Close()
}
