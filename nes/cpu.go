package nes

import "fmt"

// Addressing Modes
const (
	_ = iota
	ABSOLUTE
	ABSOLUTE_X
	ABSOLUTE_Y
	ACCUMULATOR
	IMMEDIATE
	IMPLIED
	INDEXED_INDIRECT
	INDIRECT
	INDIRECT_INDEXED
	RELATIVE
	ZERO_PAGE
	ZERO_PAGE_X
	ZERO_PAGE_Y
)

var InstructionMode = [256]byte{
	0x6, 0x7, 0x0, 0x0, 0x0, 0xb, 0xb, 0x0, 0x6, 0x5, 0x4, 0x0, 0x0, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x1, 0x7, 0x0, 0x0, 0xb, 0xb, 0xb, 0x0, 0x6, 0x5, 0x4, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x6, 0x7, 0x0, 0x0, 0x0, 0xb, 0xb, 0x0, 0x6, 0x5, 0x4, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x6, 0x7, 0x0, 0x0, 0x0, 0xb, 0xb, 0x0, 0x6, 0x5, 0x4, 0x0, 0x8, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x0, 0x7, 0x0, 0x0, 0xb, 0xb, 0xb, 0x0, 0x6, 0x0, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0xc, 0xc, 0xd, 0x0, 0x6, 0x3, 0x6, 0x0, 0x0, 0x2, 0x0, 0x0,
	0x5, 0x7, 0x5, 0x0, 0xb, 0xb, 0xb, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0xc, 0xc, 0xd, 0x0, 0x6, 0x3, 0x6, 0x0, 0x2, 0x2, 0x3, 0x0,
	0x5, 0x7, 0x0, 0x0, 0xb, 0xb, 0xb, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x5, 0x7, 0x0, 0x0, 0xb, 0xb, 0xb, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xa, 0x9, 0x0, 0x0, 0x0, 0xc, 0xc, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
}

var InstructionBytes = [256]byte{
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 0, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	0, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
}

var InstructionCycles = [256]byte{
	7, 6, 0, 0, 0, 3, 5, 0, 3, 2, 2, 0, 0, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 3, 3, 5, 0, 4, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 0, 3, 5, 0, 3, 2, 2, 0, 3, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 0, 3, 5, 0, 4, 2, 2, 0, 5, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	0, 6, 0, 0, 3, 3, 3, 0, 2, 0, 2, 0, 4, 4, 4, 0,
	2, 6, 0, 0, 4, 4, 4, 0, 2, 5, 2, 0, 0, 5, 0, 0,
	2, 6, 2, 0, 3, 3, 3, 0, 2, 2, 2, 0, 4, 4, 4, 0,
	2, 5, 0, 0, 4, 4, 4, 0, 2, 4, 2, 0, 4, 4, 4, 0,
	2, 6, 0, 0, 3, 3, 5, 0, 2, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	2, 6, 0, 0, 3, 3, 5, 0, 2, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
}

func PageCrossed(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

type CPU struct {
	Memory        // memory interface
	Cycles uint64 // number of cycles
	PC     uint16 // program counter
	SP     byte   // stack pointer
	A      byte   // accumulator
	X      byte   // x register
	Y      byte   // y register
	C      byte   // carry flag
	Z      byte   // zero flag
	I      byte   // interrupt disable flag
	D      byte   // decimal mode flag
	B      byte   // break command flag
	V      byte   // overflow flag
	N      byte   // negative flag
}

func NewCPU(memory Memory) *CPU {
	cpu := CPU{Memory: memory}
	cpu.Reset()
	return &cpu
}

func (cpu *CPU) Reset() {
	cpu.Cycles = 0
	cpu.PC = cpu.Read16(0xFFFC)
}

// Flag Functions

func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.B << 4
	flags |= cpu.V << 6
	flags |= cpu.N << 7
	return flags
}

func (cpu *CPU) SetFlags(flags byte) {
	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
}

func (cpu *CPU) SetZ(value byte) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

func (cpu *CPU) SetN(value byte) {
	if value&0x80 != 0 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
}

// Instructions

func (cpu *CPU) ADC(address uint16) {
}

func (cpu *CPU) CLD() {
	cpu.D = 0
}

func (cpu *CPU) LDA(address uint16) {
	cpu.A = cpu.Read(address)
	cpu.SetZ(cpu.A)
	cpu.SetN(cpu.A)
}

func (cpu *CPU) LDX(address uint16) {
	cpu.X = cpu.Read(address)
	cpu.SetZ(cpu.X)
	cpu.SetN(cpu.X)
}

func (cpu *CPU) SEI() {
	cpu.I = 1
}

func (cpu *CPU) TXS() {
	cpu.SP = cpu.X
}

// Step

func (cpu *CPU) Step() {
	fmt.Println(cpu)
	opcode := cpu.Read(cpu.PC)
	mode := InstructionMode[opcode]

	var address uint16
	// var pageCrossed bool
	switch mode {
	case ABSOLUTE:
		address = cpu.Read16(cpu.PC + 1)
	case ABSOLUTE_X:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.X)
		// pageCrossed = PageCrossed(address-uint16(cpu.X), address)
	case ABSOLUTE_Y:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
		// pageCrossed = PageCrossed(address-uint16(cpu.Y), address)
	case ACCUMULATOR:
		break
	case IMMEDIATE:
		address = cpu.PC + 1
	case IMPLIED:
		break
	case INDEXED_INDIRECT:
		address = cpu.Read16(uint16(cpu.Read(cpu.PC+1) + cpu.X))
	case INDIRECT:
		a := cpu.Read16(cpu.PC + 1)
		b := (a & 0xFF00) | uint16(byte(a)+1)
		lo := cpu.Read(a)
		hi := cpu.Read(b)
		address = uint16(hi)<<8 | uint16(lo)
	case INDIRECT_INDEXED:
		address = cpu.Read16(uint16(cpu.Read(cpu.PC+1))) + uint16(cpu.Y)
		// pageCrossed = PageCrossed(address-uint16(cpu.Y), address)
	case RELATIVE:
		offset := uint16(cpu.Read(cpu.PC + 1))
		if offset < 0x80 {
			address = cpu.PC + 2 + offset
		} else {
			address = cpu.PC + 2 + offset - 0x100
		}
	case ZERO_PAGE:
		address = uint16(cpu.Read(cpu.PC + 1))
	case ZERO_PAGE_X:
		address = uint16(cpu.Read(cpu.PC+1) + cpu.X)
	case ZERO_PAGE_Y:
		address = uint16(cpu.Read(cpu.PC+1) + cpu.Y)
	}

	cpu.PC += uint16(InstructionBytes[opcode])

	switch opcode {
	case 0x69:
		cpu.ADC(address)
	case 0x78:
		cpu.SEI()
	case 0x9A:
		cpu.TXS()
	case 0xA2:
		cpu.LDX(address)
	case 0xAD:
		cpu.LDA(address)
	case 0xD8:
		cpu.CLD()
	default:
		fmt.Printf("Unrecognized opcode: 0x%02x\n", opcode)
	}
}
