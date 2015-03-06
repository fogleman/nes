package main

import (
	"log"
	"os"

	"github.com/fogleman/nes/nes"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("Usage: go run main.go rom_file.nes")
	}
	nes, err := nes.NewNES(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	for {
		nes.CPU.PrintInstruction()
		nes.Step()
	}
}
