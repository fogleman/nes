package nes

import "fmt"

type CPU struct {
	Memory        // memory map
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
	cpu.PC = cpu.Read16(0xFFFC)
}

func (cpu *CPU) ReadPC() byte {
	result := cpu.Read(cpu.PC)
	cpu.PC += 1
	return result
}

func (cpu *CPU) ReadPC16() uint16 {
	result := cpu.Read16(cpu.PC)
	cpu.PC += 2
	return result
}

func (cpu *CPU) Step() {
	fmt.Println(cpu)
	cpu.ExecuteInstruction()
}

func (cpu *CPU) ExecuteInstruction() {
	opcode := cpu.ReadPC()
	switch opcode {
	case 0x69:
		cpu.ADC(cpu.Immediate())
	case 0x78:
		cpu.SEI()
	case 0x9A:
		cpu.TXS()
	case 0xA2:
		cpu.LDX(cpu.Immediate())
	case 0xAD:
		cpu.LDA(cpu.Absolute())
	case 0xD8:
		cpu.CLD()
	default:
		fmt.Printf("Unrecognized opcode: 0x%02x\n", opcode)
	}
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

// Addressing Modes

func (cpu *CPU) Immediate() uint16 {
	result := cpu.PC
	cpu.PC++
	return result
}

func (cpu *CPU) ZeroPage() uint16 {
	return uint16(cpu.ReadPC())
}

func (cpu *CPU) ZeroPageX() uint16 {
	return uint16(cpu.ReadPC() + cpu.X)
}

func (cpu *CPU) ZeroPageY() uint16 {
	return uint16(cpu.ReadPC() + cpu.Y)
}

func (cpu *CPU) Relative() uint16 {
	offset := uint16(cpu.ReadPC())
	if offset < 0x80 {
		return cpu.PC + offset
	} else {
		return cpu.PC + offset - 0x100
	}
}

func (cpu *CPU) Absolute() uint16 {
	return cpu.ReadPC16()
}

func (cpu *CPU) AbsoluteX() uint16 {
	return cpu.ReadPC16() + uint16(cpu.X)
}

func (cpu *CPU) AbsoluteY() uint16 {
	return cpu.ReadPC16() + uint16(cpu.Y)
}

func (cpu *CPU) Indirect() uint16 {
	a := cpu.ReadPC16()
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := cpu.Read(a)
	hi := cpu.Read(b)
	return uint16(hi)<<8 | uint16(lo)
}

func (cpu *CPU) IndexedIndirect() uint16 {
	return cpu.Read16(uint16(cpu.ReadPC() + cpu.X))
}

func (cpu *CPU) IndirectIndexed() uint16 {
	return cpu.Read16(uint16(cpu.ReadPC())) + uint16(cpu.Y)
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
