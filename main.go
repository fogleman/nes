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
	for i := 0; i < 10; i++ {
		nes.CPU.Step()
	}
}
