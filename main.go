package main

import (
	"log"

	idex "github.com/mathieugilbert/go-idex/pack"
)

func main() {
	ix := idex.New()

	v, err := ix.Volume24()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v\n", v)

	}

}
