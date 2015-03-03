package main

import (
	"fmt"
	"log"

	"github.com/fogleman/nes/nes"
)

func main() {
	cartridge, err := nes.LoadNESFile("roms/nestest.nes")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cartridge)
}
