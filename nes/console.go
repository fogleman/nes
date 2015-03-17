package nes

import (
	"image"
	"image/color"
	"log"
)

type Console struct {
	CPU         *CPU
	APU         *APU
	PPU         *PPU
	Cartridge   *Cartridge
	Controller1 *Controller
	Controller2 *Controller
	Mapper      Mapper
	RAM         []byte
}

func NewConsole(path string) (*Console, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ram := make([]byte, 2048)
	controller1 := NewController()
	controller2 := NewController()
	console := Console{
		nil, nil, nil, cartridge, controller1, controller2, nil, ram}
	console.Mapper = NewMapper(&console)
	console.CPU = NewCPU(&console)
	console.APU = NewAPU()
	console.PPU = NewPPU(&console)
	return &console, nil
}

func (console *Console) Step() int {
	cpuCycles := console.CPU.Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		console.PPU.Step()
		console.Mapper.Step()
	}
	return cpuCycles
}

func (console *Console) StepFrame() {
	frame := console.PPU.Frame
	for frame == console.PPU.Frame {
		console.Step()
	}
}

func (console *Console) Buffer() *image.RGBA {
	return console.PPU.buffer
}

func (console *Console) BackgroundColor() color.RGBA {
	return palette[console.PPU.readPalette(0)%64]
}

func (console *Console) SetPressed(controller, button int, pressed bool) {
	switch controller {
	case 1:
		console.Controller1.SetPressed(button, pressed)
	case 2:
		console.Controller2.SetPressed(button, pressed)
	default:
		log.Fatalf("unhandled controller press: %d", controller)
	}
}
