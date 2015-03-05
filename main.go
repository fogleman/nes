package main

import "github.com/fogleman/nes/ui"

func main() {
	// nes, err := nes.NewNES("roms/nestest.nes")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for {
	// 	nes.CPU.Step()
	// }
	ui.Run()
}
