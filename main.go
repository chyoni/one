package main

import (
	"github.com/chiwon99881/one/api"
	"github.com/chiwon99881/one/db"
)

func main() {
	api.Start("4000")
	defer db.Close()
}
