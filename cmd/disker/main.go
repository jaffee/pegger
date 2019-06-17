package main

import (
	"log"

	"github.com/jaffee/commandeer"
	"github.com/jaffee/pegger"
)

func main() {
	if err := commandeer.Run(pegger.NewDisker()); err != nil {
		log.Fatalf("Error Running: %v\n", err)
	}
}
