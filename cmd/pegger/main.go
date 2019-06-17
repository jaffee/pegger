package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/jaffee/commandeer"
	"github.com/jaffee/pegger"
)

func main() {
	if err := commandeer.Run(pegger.NewPegger()); err != nil {
		log.Fatal(err)
	}
}
