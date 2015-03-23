package main

import (
	"log"
	"os"

	"github.com/fogleman/nes/nes"
	"github.com/fogleman/nes/ui"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("Usage: nes rom_file.nes")
	}
	console, err := nes.NewConsole(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	ui.Run(console)
}
