package nes

import (
	"fmt"
	"strings"
	"testing"
)

func readString(cpu *CPU, address uint16) string {
	var bytes []byte
	for {
		c := cpu.Read(address)
		if c == 0 {
			break
		}
		bytes = append(bytes, c)
		address++
	}
	return string(bytes)
}

// instr_test: http://wiki.nesdev.com/w/index.php/Emulator_tests
func TestOfficialInstructions(t *testing.T) {
	nes, err := NewNES("../roms/official_instructions.nes")
	if err != nil {
		t.Fatal(err)
	}
	cpu := nes.CPU
	cpu.Write(0x6000, 0xFF)
	for {
		for i := 0; i < 65536; i++ {
			cpu.Step()
		}
		if cpu.Read(0x6000) < 0x80 {
			fmt.Println(cpu.Cycles)
			break
		}
		message := strings.TrimSpace(readString(cpu, 0x6004))
		if len(message) > 0 {
			fmt.Println(message)
		}
	}
}
