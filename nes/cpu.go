package nes

import (
	"fmt"
	"log"
)

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

var InstructionNames = [256]string{
	"BRK", "ORA", "???", "???", "???", "ORA", "ASL", "???",
	"PHP", "ORA", "ASL", "???", "???", "ORA", "ASL", "???",
	"BPL", "ORA", "???", "???", "???", "ORA", "ASL", "???",
	"CLC", "ORA", "???", "???", "???", "ORA", "ASL", "???",
	"JSR", "AND", "???", "???", "BIT", "AND", "ROL", "???",
	"PLP", "AND", "ROL", "???", "BIT", "AND", "ROL", "???",
	"BMI", "AND", "???", "???", "???", "AND", "ROL", "???",
	"SEC", "AND", "???", "???", "???", "AND", "ROL", "???",
	"RTI", "EOR", "???", "???", "???", "EOR", "LSR", "???",
	"PHA", "EOR", "LSR", "???", "JMP", "EOR", "LSR", "???",
	"BVC", "EOR", "???", "???", "???", "EOR", "LSR", "???",
	"CLI", "EOR", "???", "???", "???", "EOR", "LSR", "???",
	"RTS", "ADC", "???", "???", "???", "ADC", "ROR", "???",
	"PLA", "ADC", "ROR", "???", "JMP", "ADC", "ROR", "???",
	"BVS", "ADC", "???", "???", "???", "ADC", "ROR", "???",
	"SEI", "ADC", "???", "???", "???", "ADC", "ROR", "???",
	"???", "STA", "???", "???", "STY", "STA", "STX", "???",
	"DEY", "???", "TXA", "???", "STY", "STA", "STX", "???",
	"BCC", "STA", "???", "???", "STY", "STA", "STX", "???",
	"TYA", "STA", "TXS", "???", "???", "STA", "???", "???",
	"LDY", "LDA", "LDX", "???", "LDY", "LDA", "LDX", "???",
	"TAY", "LDA", "TAX", "???", "LDY", "LDA", "LDX", "???",
	"BCS", "LDA", "???", "???", "LDY", "LDA", "LDX", "???",
	"CLV", "LDA", "TSX", "???", "LDY", "LDA", "LDX", "???",
	"CPY", "CMP", "???", "???", "CPY", "CMP", "DEC", "???",
	"INY", "CMP", "DEX", "???", "CPY", "CMP", "DEC", "???",
	"BNE", "CMP", "???", "???", "???", "CMP", "DEC", "???",
	"CLD", "CMP", "???", "???", "???", "CMP", "DEC", "???",
	"CPX", "SBC", "???", "???", "CPX", "SBC", "INC", "???",
	"INX", "SBC", "NOP", "???", "CPX", "SBC", "INC", "???",
	"BEQ", "SBC", "???", "???", "???", "SBC", "INC", "???",
	"SED", "SBC", "???", "???", "???", "SBC", "INC", "???",
}

var InstructionModes = [256]byte{
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
	Memory // memory interface
	Table  [256]func(uint16)
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
	cpu.Table = [256]func(uint16){
		cpu.BRK, cpu.ORA, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ORA, cpu.ASL, cpu.UNK,
		cpu.PHP, cpu.ORA, cpu.ASL, cpu.UNK, cpu.UNK, cpu.ORA, cpu.ASL, cpu.UNK,
		cpu.BPL, cpu.ORA, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ORA, cpu.ASL, cpu.UNK,
		cpu.CLC, cpu.ORA, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ORA, cpu.ASL, cpu.UNK,
		cpu.JSR, cpu.AND, cpu.UNK, cpu.UNK, cpu.BIT, cpu.AND, cpu.ROL, cpu.UNK,
		cpu.PLP, cpu.AND, cpu.ROL, cpu.UNK, cpu.BIT, cpu.AND, cpu.ROL, cpu.UNK,
		cpu.BMI, cpu.AND, cpu.UNK, cpu.UNK, cpu.UNK, cpu.AND, cpu.ROL, cpu.UNK,
		cpu.SEC, cpu.AND, cpu.UNK, cpu.UNK, cpu.UNK, cpu.AND, cpu.ROL, cpu.UNK,
		cpu.RTI, cpu.EOR, cpu.UNK, cpu.UNK, cpu.UNK, cpu.EOR, cpu.LSR, cpu.UNK,
		cpu.PHA, cpu.EOR, cpu.LSR, cpu.UNK, cpu.JMP, cpu.EOR, cpu.LSR, cpu.UNK,
		cpu.BVC, cpu.EOR, cpu.UNK, cpu.UNK, cpu.UNK, cpu.EOR, cpu.LSR, cpu.UNK,
		cpu.CLI, cpu.EOR, cpu.UNK, cpu.UNK, cpu.UNK, cpu.EOR, cpu.LSR, cpu.UNK,
		cpu.RTS, cpu.ADC, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ADC, cpu.ROR, cpu.UNK,
		cpu.PLA, cpu.ADC, cpu.ROR, cpu.UNK, cpu.JMP, cpu.ADC, cpu.ROR, cpu.UNK,
		cpu.BVS, cpu.ADC, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ADC, cpu.ROR, cpu.UNK,
		cpu.SEI, cpu.ADC, cpu.UNK, cpu.UNK, cpu.UNK, cpu.ADC, cpu.ROR, cpu.UNK,
		cpu.UNK, cpu.STA, cpu.UNK, cpu.UNK, cpu.STY, cpu.STA, cpu.STX, cpu.UNK,
		cpu.DEY, cpu.UNK, cpu.TXA, cpu.UNK, cpu.STY, cpu.STA, cpu.STX, cpu.UNK,
		cpu.BCC, cpu.STA, cpu.UNK, cpu.UNK, cpu.STY, cpu.STA, cpu.STX, cpu.UNK,
		cpu.TYA, cpu.STA, cpu.TXS, cpu.UNK, cpu.UNK, cpu.STA, cpu.UNK, cpu.UNK,
		cpu.LDY, cpu.LDA, cpu.LDX, cpu.UNK, cpu.LDY, cpu.LDA, cpu.LDX, cpu.UNK,
		cpu.TAY, cpu.LDA, cpu.TAX, cpu.UNK, cpu.LDY, cpu.LDA, cpu.LDX, cpu.UNK,
		cpu.BCS, cpu.LDA, cpu.UNK, cpu.UNK, cpu.LDY, cpu.LDA, cpu.LDX, cpu.UNK,
		cpu.CLV, cpu.LDA, cpu.TSX, cpu.UNK, cpu.LDY, cpu.LDA, cpu.LDX, cpu.UNK,
		cpu.CPY, cpu.CMP, cpu.UNK, cpu.UNK, cpu.CPY, cpu.CMP, cpu.DEC, cpu.UNK,
		cpu.INY, cpu.CMP, cpu.DEX, cpu.UNK, cpu.CPY, cpu.CMP, cpu.DEC, cpu.UNK,
		cpu.BNE, cpu.CMP, cpu.UNK, cpu.UNK, cpu.UNK, cpu.CMP, cpu.DEC, cpu.UNK,
		cpu.CLD, cpu.CMP, cpu.UNK, cpu.UNK, cpu.UNK, cpu.CMP, cpu.DEC, cpu.UNK,
		cpu.CPX, cpu.SBC, cpu.UNK, cpu.UNK, cpu.CPX, cpu.SBC, cpu.INC, cpu.UNK,
		cpu.INX, cpu.SBC, cpu.NOP, cpu.UNK, cpu.CPX, cpu.SBC, cpu.INC, cpu.UNK,
		cpu.BEQ, cpu.SBC, cpu.UNK, cpu.UNK, cpu.UNK, cpu.SBC, cpu.INC, cpu.UNK,
		cpu.SED, cpu.SBC, cpu.UNK, cpu.UNK, cpu.UNK, cpu.SBC, cpu.INC, cpu.UNK,
	}
	cpu.Reset()
	return &cpu
}

func (cpu *CPU) Reset() {
	fmt.Printf("%T\n", cpu.ADC)
	cpu.Cycles = 0
	// cpu.PC = cpu.Read16(0xFFFC)
	cpu.PC = 0xC000
	cpu.SP = 0xFD
	cpu.SetFlags(0x24)
}

// Flag Functions

func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.B << 4
	flags |= 1 << 5
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

// Step

func (cpu *CPU) PrintInstruction() {
	opcode := cpu.Read(cpu.PC)
	bytes := InstructionBytes[opcode]
	name := InstructionNames[opcode]
	w0 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+0))
	w1 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+1))
	w2 := fmt.Sprintf("%02X", cpu.Read(cpu.PC+2))
	if bytes < 2 {
		w1 = "  "
	}
	if bytes < 3 {
		w2 = "  "
	}
	fmt.Printf(
		"%4X  %s %s %s  %s %27s"+
			"A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%3d\n",
		cpu.PC, w0, w1, w2, name, "",
		cpu.A, cpu.X, cpu.Y, cpu.Flags(), cpu.SP, cpu.Cycles*3)
}

func (cpu *CPU) Step() {
	cpu.PrintInstruction()

	opcode := cpu.Read(cpu.PC)
	mode := InstructionModes[opcode]

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
	cpu.Cycles += uint64(InstructionCycles[opcode])

	cpu.Table[opcode](address)
}

// Instructions
func (cpu *CPU) ADC(address uint16) {
	log.Fatalln("Unimplemented instruction: ADC")
}

func (cpu *CPU) AND(address uint16) {
	log.Fatalln("Unimplemented instruction: AND")
}

func (cpu *CPU) ASL(address uint16) {
	log.Fatalln("Unimplemented instruction: ASL")
}

func (cpu *CPU) BCC(address uint16) {
	log.Fatalln("Unimplemented instruction: BCC")
}

func (cpu *CPU) BCS(address uint16) {
	log.Fatalln("Unimplemented instruction: BCS")
}

func (cpu *CPU) BEQ(address uint16) {
	log.Fatalln("Unimplemented instruction: BEQ")
}

func (cpu *CPU) BIT(address uint16) {
	log.Fatalln("Unimplemented instruction: BIT")
}

func (cpu *CPU) BMI(address uint16) {
	log.Fatalln("Unimplemented instruction: BMI")
}

func (cpu *CPU) BNE(address uint16) {
	log.Fatalln("Unimplemented instruction: BNE")
}

func (cpu *CPU) BPL(address uint16) {
	log.Fatalln("Unimplemented instruction: BPL")
}

func (cpu *CPU) BRK(address uint16) {
	log.Fatalln("Unimplemented instruction: BRK")
}

func (cpu *CPU) BVC(address uint16) {
	log.Fatalln("Unimplemented instruction: BVC")
}

func (cpu *CPU) BVS(address uint16) {
	log.Fatalln("Unimplemented instruction: BVS")
}

func (cpu *CPU) CLC(address uint16) {
	log.Fatalln("Unimplemented instruction: CLC")
}

func (cpu *CPU) CLD(address uint16) {
	cpu.D = 0
}

func (cpu *CPU) CLI(address uint16) {
	log.Fatalln("Unimplemented instruction: CLI")
}

func (cpu *CPU) CLV(address uint16) {
	log.Fatalln("Unimplemented instruction: CLV")
}

func (cpu *CPU) CMP(address uint16) {
	log.Fatalln("Unimplemented instruction: CMP")
}

func (cpu *CPU) CPX(address uint16) {
	log.Fatalln("Unimplemented instruction: CPX")
}

func (cpu *CPU) CPY(address uint16) {
	log.Fatalln("Unimplemented instruction: CPY")
}

func (cpu *CPU) DEC(address uint16) {
	log.Fatalln("Unimplemented instruction: DEC")
}

func (cpu *CPU) DEX(address uint16) {
	log.Fatalln("Unimplemented instruction: DEX")
}

func (cpu *CPU) DEY(address uint16) {
	log.Fatalln("Unimplemented instruction: DEY")
}

func (cpu *CPU) EOR(address uint16) {
	log.Fatalln("Unimplemented instruction: EOR")
}

func (cpu *CPU) INC(address uint16) {
	log.Fatalln("Unimplemented instruction: INC")
}

func (cpu *CPU) INX(address uint16) {
	log.Fatalln("Unimplemented instruction: INX")
}

func (cpu *CPU) INY(address uint16) {
	log.Fatalln("Unimplemented instruction: INY")
}

func (cpu *CPU) JMP(address uint16) {
	cpu.PC = address
}

func (cpu *CPU) JSR(address uint16) {
	log.Fatalln("Unimplemented instruction: JSR")
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

func (cpu *CPU) LDY(address uint16) {
	cpu.Y = cpu.Read(address)
	cpu.SetZ(cpu.Y)
	cpu.SetN(cpu.Y)
}

func (cpu *CPU) LSR(address uint16) {
	log.Fatalln("Unimplemented instruction: LSR")
}

func (cpu *CPU) NOP(address uint16) {
	log.Fatalln("Unimplemented instruction: NOP")
}

func (cpu *CPU) ORA(address uint16) {
	log.Fatalln("Unimplemented instruction: ORA")
}

func (cpu *CPU) PHA(address uint16) {
	log.Fatalln("Unimplemented instruction: PHA")
}

func (cpu *CPU) PHP(address uint16) {
	log.Fatalln("Unimplemented instruction: PHP")
}

func (cpu *CPU) PLA(address uint16) {
	log.Fatalln("Unimplemented instruction: PLA")
}

func (cpu *CPU) PLP(address uint16) {
	log.Fatalln("Unimplemented instruction: PLP")
}

func (cpu *CPU) ROL(address uint16) {
	log.Fatalln("Unimplemented instruction: ROL")
}

func (cpu *CPU) ROR(address uint16) {
	log.Fatalln("Unimplemented instruction: ROR")
}

func (cpu *CPU) RTI(address uint16) {
	log.Fatalln("Unimplemented instruction: RTI")
}

func (cpu *CPU) RTS(address uint16) {
	log.Fatalln("Unimplemented instruction: RTS")
}

func (cpu *CPU) SBC(address uint16) {
	log.Fatalln("Unimplemented instruction: SBC")
}

func (cpu *CPU) SEC(address uint16) {
	log.Fatalln("Unimplemented instruction: SEC")
}

func (cpu *CPU) SED(address uint16) {
	log.Fatalln("Unimplemented instruction: SED")
}

func (cpu *CPU) SEI(address uint16) {
	cpu.I = 1
}

func (cpu *CPU) STA(address uint16) {
	log.Fatalln("Unimplemented instruction: STA")
}

func (cpu *CPU) STX(address uint16) {
	cpu.Write(address, cpu.X)
}

func (cpu *CPU) STY(address uint16) {
	log.Fatalln("Unimplemented instruction: STY")
}

func (cpu *CPU) TAX(address uint16) {
	log.Fatalln("Unimplemented instruction: TAX")
}

func (cpu *CPU) TAY(address uint16) {
	log.Fatalln("Unimplemented instruction: TAY")
}

func (cpu *CPU) TSX(address uint16) {
	log.Fatalln("Unimplemented instruction: TSX")
}

func (cpu *CPU) TXA(address uint16) {
	log.Fatalln("Unimplemented instruction: TXA")
}

func (cpu *CPU) TXS(address uint16) {
	cpu.SP = cpu.X
}

func (cpu *CPU) TYA(address uint16) {
	log.Fatalln("Unimplemented instruction: TYA")
}

func (cpu *CPU) UNK(address uint16) {
	log.Fatalln("Unimplemented instruction: UNK")
}
