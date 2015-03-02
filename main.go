package main

import (
	"log"

	"github.com/fogleman/nes/nes"
)

func main() {
	if err := nes.LoadNESFile("roms/nestest.nes"); err != nil {
		log.Fatalln(err)
	}
}
