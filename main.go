package main

import (
	"log"

	"github.com/fogleman/nes/nes"
)

func main() {
	nes, err := nes.NewNES("roms/nestest.nes")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		nes.CPU.Step()
	}
}
