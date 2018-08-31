package main

import (
	"log"

	idex "github.com/mathieugilbert/go-idex/pack"
)

func main() {
	ix := idex.New()
	ticker, err := ix.Ticker("ETH_SAN")
	if err != nil {
		log.Printf("%+v\n", err)
	}
	log.Printf("%+v\n", ticker)
}
