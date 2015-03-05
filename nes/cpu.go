package nes

import (
	"fmt"
	"log"
)

// addressing modes
const (
	_ = iota
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

// instructionModes indicates the addressing mode for each instruction
var instructionModes = [256]byte{
	6, 7, 0, 0, 0, 11, 11, 0, 6, 5, 4, 0, 0, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
	1, 7, 0, 0, 11, 11, 11, 0, 6, 5, 4, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
	6, 7, 0, 0, 0, 11, 11, 0, 6, 5, 4, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
	6, 7, 0, 0, 0, 11, 11, 0, 6, 5, 4, 0, 8, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
	0, 7, 0, 0, 11, 11, 11, 0, 6, 0, 6, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 12, 12, 13, 0, 6, 3, 6, 0, 0, 2, 0, 0,
	5, 7, 5, 0, 11, 11, 11, 0, 6, 5, 6, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 12, 12, 13, 0, 6, 3, 6, 0, 2, 2, 3, 0,
	5, 7, 0, 0, 11, 11, 11, 0, 6, 5, 6, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
	5, 7, 0, 0, 11, 11, 11, 0, 6, 5, 6, 0, 1, 1, 1, 0,
	10, 9, 0, 0, 0, 12, 12, 0, 6, 3, 0, 0, 0, 2, 2, 0,
}

// instructionSizes indicates the size of each instruction in bytes
var instructionSizes = [256]byte{
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

// instructionCycles indicates the number of cycles used by each instruction,
// not including conditional cycles
var instructionCycles = [256]byte{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

// instructionPageCycles indicates the number of cycles used by each
// instruction when a page is crossed
var instructionPageCycles = [256]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0,
}

// instructionNames indicates the name of each instruction
var instructionNames = [256]string{
	"BRK", "ORA", "UNK", "UNK", "UNK", "ORA", "ASL", "UNK",
	"PHP", "ORA", "ASL", "UNK", "UNK", "ORA", "ASL", "UNK",
	"BPL", "ORA", "UNK", "UNK", "UNK", "ORA", "ASL", "UNK",
	"CLC", "ORA", "UNK", "UNK", "UNK", "ORA", "ASL", "UNK",
	"JSR", "AND", "UNK", "UNK", "BIT", "AND", "ROL", "UNK",
	"PLP", "AND", "ROL", "UNK", "BIT", "AND", "ROL", "UNK",
	"BMI", "AND", "UNK", "UNK", "UNK", "AND", "ROL", "UNK",
	"SEC", "AND", "UNK", "UNK", "UNK", "AND", "ROL", "UNK",
	"RTI", "EOR", "UNK", "UNK", "UNK", "EOR", "LSR", "UNK",
	"PHA", "EOR", "LSR", "UNK", "JMP", "EOR", "LSR", "UNK",
	"BVC", "EOR", "UNK", "UNK", "UNK", "EOR", "LSR", "UNK",
	"CLI", "EOR", "UNK", "UNK", "UNK", "EOR", "LSR", "UNK",
	"RTS", "ADC", "UNK", "UNK", "UNK", "ADC", "ROR", "UNK",
	"PLA", "ADC", "ROR", "UNK", "JMP", "ADC", "ROR", "UNK",
	"BVS", "ADC", "UNK", "UNK", "UNK", "ADC", "ROR", "UNK",
	"SEI", "ADC", "UNK", "UNK", "UNK", "ADC", "ROR", "UNK",
	"UNK", "STA", "UNK", "UNK", "STY", "STA", "STX", "UNK",
	"DEY", "UNK", "TXA", "UNK", "STY", "STA", "STX", "UNK",
	"BCC", "STA", "UNK", "UNK", "STY", "STA", "STX", "UNK",
	"TYA", "STA", "TXS", "UNK", "UNK", "STA", "UNK", "UNK",
	"LDY", "LDA", "LDX", "UNK", "LDY", "LDA", "LDX", "UNK",
	"TAY", "LDA", "TAX", "UNK", "LDY", "LDA", "LDX", "UNK",
	"BCS", "LDA", "UNK", "UNK", "LDY", "LDA", "LDX", "UNK",
	"CLV", "LDA", "TSX", "UNK", "LDY", "LDA", "LDX", "UNK",
	"CPY", "CMP", "UNK", "UNK", "CPY", "CMP", "DEC", "UNK",
	"INY", "CMP", "DEX", "UNK", "CPY", "CMP", "DEC", "UNK",
	"BNE", "CMP", "UNK", "UNK", "UNK", "CMP", "DEC", "UNK",
	"CLD", "CMP", "UNK", "UNK", "UNK", "CMP", "DEC", "UNK",
	"CPX", "SBC", "UNK", "UNK", "CPX", "SBC", "INC", "UNK",
	"INX", "SBC", "NOP", "UNK", "CPX", "SBC", "INC", "UNK",
	"BEQ", "SBC", "UNK", "UNK", "UNK", "SBC", "INC", "UNK",
	"SED", "SBC", "UNK", "UNK", "UNK", "SBC", "INC", "UNK",
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
	U      byte   // unused flag
	V      byte   // overflow flag
	N      byte   // negative flag
	table  [256]func(*stepInfo)
}

func NewCPU(memory Memory) *CPU {
	cpu := CPU{Memory: memory}
	cpu.createTable()
	cpu.Reset()
	return &cpu
}

// createTable builds a function table for each instruction
func (c *CPU) createTable() {
	c.table = [256]func(*stepInfo){
		c.brk, c.ora, c.unk, c.unk, c.unk, c.ora, c.asl, c.unk,
		c.php, c.ora, c.asl, c.unk, c.unk, c.ora, c.asl, c.unk,
		c.bpl, c.ora, c.unk, c.unk, c.unk, c.ora, c.asl, c.unk,
		c.clc, c.ora, c.unk, c.unk, c.unk, c.ora, c.asl, c.unk,
		c.jsr, c.and, c.unk, c.unk, c.bit, c.and, c.rol, c.unk,
		c.plp, c.and, c.rol, c.unk, c.bit, c.and, c.rol, c.unk,
		c.bmi, c.and, c.unk, c.unk, c.unk, c.and, c.rol, c.unk,
		c.sec, c.and, c.unk, c.unk, c.unk, c.and, c.rol, c.unk,
		c.rti, c.eor, c.unk, c.unk, c.unk, c.eor, c.lsr, c.unk,
		c.pha, c.eor, c.lsr, c.unk, c.jmp, c.eor, c.lsr, c.unk,
		c.bvc, c.eor, c.unk, c.unk, c.unk, c.eor, c.lsr, c.unk,
		c.cli, c.eor, c.unk, c.unk, c.unk, c.eor, c.lsr, c.unk,
		c.rts, c.adc, c.unk, c.unk, c.unk, c.adc, c.ror, c.unk,
		c.pla, c.adc, c.ror, c.unk, c.jmp, c.adc, c.ror, c.unk,
		c.bvs, c.adc, c.unk, c.unk, c.unk, c.adc, c.ror, c.unk,
		c.sei, c.adc, c.unk, c.unk, c.unk, c.adc, c.ror, c.unk,
		c.unk, c.sta, c.unk, c.unk, c.sty, c.sta, c.stx, c.unk,
		c.dey, c.unk, c.txa, c.unk, c.sty, c.sta, c.stx, c.unk,
		c.bcc, c.sta, c.unk, c.unk, c.sty, c.sta, c.stx, c.unk,
		c.tya, c.sta, c.txs, c.unk, c.unk, c.sta, c.unk, c.unk,
		c.ldy, c.lda, c.ldx, c.unk, c.ldy, c.lda, c.ldx, c.unk,
		c.tay, c.lda, c.tax, c.unk, c.ldy, c.lda, c.ldx, c.unk,
		c.bcs, c.lda, c.unk, c.unk, c.ldy, c.lda, c.ldx, c.unk,
		c.clv, c.lda, c.tsx, c.unk, c.ldy, c.lda, c.ldx, c.unk,
		c.cpy, c.cmp, c.unk, c.unk, c.cpy, c.cmp, c.dec, c.unk,
		c.iny, c.cmp, c.dex, c.unk, c.cpy, c.cmp, c.dec, c.unk,
		c.bne, c.cmp, c.unk, c.unk, c.unk, c.cmp, c.dec, c.unk,
		c.cld, c.cmp, c.unk, c.unk, c.unk, c.cmp, c.dec, c.unk,
		c.cpx, c.sbc, c.unk, c.unk, c.cpx, c.sbc, c.inc, c.unk,
		c.inx, c.sbc, c.nop, c.unk, c.cpx, c.sbc, c.inc, c.unk,
		c.beq, c.sbc, c.unk, c.unk, c.unk, c.sbc, c.inc, c.unk,
		c.sed, c.sbc, c.unk, c.unk, c.unk, c.sbc, c.inc, c.unk,
	}
}

// Reset resets the CPU to its initial powerup state
func (cpu *CPU) Reset() {
	cpu.Cycles = 0
	// cpu.PC = cpu.Read16(0xFFFC)
	cpu.PC = 0xC000
	cpu.SP = 0xFD
	cpu.SetFlags(0x24)
}

// pagesDiffer returns true if the two addresses are within different pages
func pagesDiffer(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// read16bug emulates a 6502 bug that caused the low byte to wrap without
// incrementing the high byte
func (cpu *CPU) read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := cpu.Read(a)
	hi := cpu.Read(b)
	return uint16(hi)<<8 | uint16(lo)
}

// printInstruction prints the current CPU state
func (cpu *CPU) printInstruction() {
	opcode := cpu.Read(cpu.PC)
	bytes := instructionSizes[opcode]
	name := instructionNames[opcode]
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
		"%4X  %s %s %s  %s %28s"+
			"A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%3d\n",
		cpu.PC, w0, w1, w2, name, "",
		cpu.A, cpu.X, cpu.Y, cpu.Flags(), cpu.SP, (cpu.Cycles*3)%341)
}

// push pushes a byte onto the stack
func (cpu *CPU) push(value byte) {
	cpu.Write(0x100|uint16(cpu.SP), value)
	cpu.SP--
}

// pull pops a byte from the stack
func (cpu *CPU) pull() byte {
	cpu.SP++
	return cpu.Read(0x100 | uint16(cpu.SP))
}

// push pushes two bytes onto the stack
func (cpu *CPU) push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	cpu.push(hi)
	cpu.push(lo)
}

// pull pops two bytes from the stack
func (cpu *CPU) pull16() uint16 {
	lo := uint16(cpu.pull())
	hi := uint16(cpu.pull())
	return hi<<8 | lo
}

// Flags returns the processor status flags
func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.B << 4
	flags |= cpu.U << 5
	flags |= cpu.V << 6
	flags |= cpu.N << 7
	return flags
}

// SetFlags sets the processor status flags
func (cpu *CPU) SetFlags(flags byte) {
	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
}

// setZ sets the zero flag if the argument is zero
func (cpu *CPU) setZ(value byte) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func (cpu *CPU) setN(value byte) {
	if value&0x80 != 0 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
}

// stepInfo contains information that the instruction functions to use
type stepInfo struct {
	address uint16
	mode    byte
}

// Step executes a single CPU instruction
func (cpu *CPU) Step() {
	opcode := cpu.Read(cpu.PC)
	mode := instructionModes[opcode]

	var address uint16
	var pageCrossed bool
	switch mode {
	case modeAbsolute:
		address = cpu.Read16(cpu.PC + 1)
	case modeAbsoluteX:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.X)
		pageCrossed = pagesDiffer(address-uint16(cpu.X), address)
	case modeAbsoluteY:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
		pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	case modeAccumulator:
		address = 0
	case modeImmediate:
		address = cpu.PC + 1
	case modeImplied:
		address = 0
	case modeIndexedIndirect:
		address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1) + cpu.X))
	case modeIndirect:
		address = cpu.read16bug(cpu.Read16(cpu.PC + 1))
	case modeIndirectIndexed:
		address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1))) + uint16(cpu.Y)
		pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	case modeRelative:
		offset := uint16(cpu.Read(cpu.PC + 1))
		if offset < 0x80 {
			address = cpu.PC + 2 + offset
		} else {
			address = cpu.PC + 2 + offset - 0x100
		}
	case modeZeroPage:
		address = uint16(cpu.Read(cpu.PC + 1))
	case modeZeroPageX:
		address = uint16(cpu.Read(cpu.PC+1) + cpu.X)
	case modeZeroPageY:
		address = uint16(cpu.Read(cpu.PC+1) + cpu.Y)
	}

	cpu.PC += uint16(instructionSizes[opcode])
	cpu.Cycles += uint64(instructionCycles[opcode])
	if pageCrossed {
		cpu.Cycles += uint64(instructionPageCycles[opcode])
	}

	info := &stepInfo{address, mode}
	cpu.table[opcode](info)
}

// ADC - Add with Carry
func (cpu *CPU) adc(info *stepInfo) {
	a := cpu.A
	b := cpu.Read(info.address)
	c := cpu.C
	cpu.A = a + b + c
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
	if int(a)+int(b)+int(c) > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^b)&0x80 == 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

// AND - Logical AND
func (cpu *CPU) and(info *stepInfo) {
	cpu.A = cpu.A & cpu.Read(info.address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// ASL - Arithmetic Shift Left
func (cpu *CPU) asl(info *stepInfo) {
	if info.mode == modeAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A <<= 1
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		value := cpu.Read(info.address)
		cpu.C = (value >> 7) & 1
		value <<= 1
		cpu.Write(info.address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// BCC - Branch if Carry Clear
func (cpu *CPU) bcc(info *stepInfo) {
	if cpu.C == 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BCS - Branch if Carry Set
func (cpu *CPU) bcs(info *stepInfo) {
	if cpu.C != 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BEQ - Branch if Equal
func (cpu *CPU) beq(info *stepInfo) {
	if cpu.Z != 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BIT - Bit Test
func (cpu *CPU) bit(info *stepInfo) {
	value := cpu.Read(info.address)
	cpu.V = (value >> 6) & 1
	cpu.setZ(value & cpu.A)
	cpu.setN(value)
}

// BMI - Branch if Minus
func (cpu *CPU) bmi(info *stepInfo) {
	if cpu.N != 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BNE - Branch if Not Equal
func (cpu *CPU) bne(info *stepInfo) {
	if cpu.Z == 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BPL - Branch if Positive
func (cpu *CPU) bpl(info *stepInfo) {
	if cpu.N == 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BRK - Force Interrupt
func (cpu *CPU) brk(info *stepInfo) {
	log.Fatalln("Unimplemented instruction: BRK")
}

// BVC - Branch if Overflow Clear
func (cpu *CPU) bvc(info *stepInfo) {
	if cpu.V == 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// BVS - Branch if Overflow Set
func (cpu *CPU) bvs(info *stepInfo) {
	if cpu.V != 0 {
		cpu.PC = info.address
		cpu.Cycles++
	}
}

// CLC - Clear Carry Flag
func (cpu *CPU) clc(info *stepInfo) {
	cpu.C = 0
}

// CLD - Clear Decimal Mode
func (cpu *CPU) cld(info *stepInfo) {
	cpu.D = 0
}

// CLI - Clear Interrupt Disable
func (cpu *CPU) cli(info *stepInfo) {
	cpu.I = 0
}

// CLV - Clear Overflow Flag
func (cpu *CPU) clv(info *stepInfo) {
	cpu.V = 0
}

// CMP - Compare
func (cpu *CPU) cmp(info *stepInfo) {
	M := cpu.Read(info.address)
	value := cpu.A - M
	cpu.setZ(value)
	cpu.setN(value)
	if cpu.A >= M {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// CPX - Compare X Register
func (cpu *CPU) cpx(info *stepInfo) {
	M := cpu.Read(info.address)
	value := cpu.X - M
	cpu.setZ(value)
	cpu.setN(value)
	if cpu.X >= M {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// CPY - Compare Y Register
func (cpu *CPU) cpy(info *stepInfo) {
	M := cpu.Read(info.address)
	value := cpu.Y - M
	cpu.setZ(value)
	cpu.setN(value)
	if cpu.Y >= M {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// DEC - Decrement Memory
func (cpu *CPU) dec(info *stepInfo) {
	value := cpu.Read(info.address) - 1
	cpu.Write(info.address, value)
	cpu.setZ(value)
	cpu.setN(value)
}

// DEX - Decrement X Register
func (cpu *CPU) dex(info *stepInfo) {
	cpu.X--
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// DEY - Decrement Y Register
func (cpu *CPU) dey(info *stepInfo) {
	cpu.Y--
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// EOR - Exclusive OR
func (cpu *CPU) eor(info *stepInfo) {
	cpu.A = cpu.A ^ cpu.Read(info.address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// INC - Increment Memory
func (cpu *CPU) inc(info *stepInfo) {
	value := cpu.Read(info.address) + 1
	cpu.Write(info.address, value)
	cpu.setZ(value)
	cpu.setN(value)
}

// INX - Increment X Register
func (cpu *CPU) inx(info *stepInfo) {
	cpu.X++
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// INY - Increment Y Register
func (cpu *CPU) iny(info *stepInfo) {
	cpu.Y++
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// JMP - Jump
func (cpu *CPU) jmp(info *stepInfo) {
	cpu.PC = info.address
}

// JSR - Jump to Subroutine
func (cpu *CPU) jsr(info *stepInfo) {
	cpu.push16(cpu.PC - 1)
	cpu.PC = info.address
}

// LDA - Load Accumulator
func (cpu *CPU) lda(info *stepInfo) {
	cpu.A = cpu.Read(info.address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// LDX - Load X Register
func (cpu *CPU) ldx(info *stepInfo) {
	cpu.X = cpu.Read(info.address)
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// LDY - Load Y Register
func (cpu *CPU) ldy(info *stepInfo) {
	cpu.Y = cpu.Read(info.address)
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// LSR - Logical Shift Right
func (cpu *CPU) lsr(info *stepInfo) {
	if info.mode == modeAccumulator {
		cpu.C = cpu.A & 1
		cpu.A >>= 1
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		value := cpu.Read(info.address)
		cpu.C = value & 1
		value >>= 1
		cpu.Write(info.address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// NOP - No Operation
func (cpu *CPU) nop(info *stepInfo) {
}

// ORA - Logical Inclusive OR
func (cpu *CPU) ora(info *stepInfo) {
	cpu.A = cpu.A | cpu.Read(info.address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// PHA - Push Accumulator
func (cpu *CPU) pha(info *stepInfo) {
	cpu.push(cpu.A)
}

// PHP - Push Processor Status
func (cpu *CPU) php(info *stepInfo) {
	cpu.push(cpu.Flags() | 0x10)
}

// PLA - Pull Accumulator
func (cpu *CPU) pla(info *stepInfo) {
	cpu.A = cpu.pull()
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// PLP - Pull Processor Status
func (cpu *CPU) plp(info *stepInfo) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
}

// ROL - Rotate Left
func (cpu *CPU) rol(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := cpu.C
		cpu.C = (cpu.A >> 7) & 1
		cpu.A = (cpu.A << 1) | c
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		c := cpu.C
		value := cpu.Read(info.address)
		cpu.C = (value >> 7) & 1
		value = (value << 1) | c
		cpu.Write(info.address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// ROR - Rotate Right
func (cpu *CPU) ror(info *stepInfo) {
	if info.mode == modeAccumulator {
		c := cpu.C
		cpu.C = cpu.A & 1
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		c := cpu.C
		value := cpu.Read(info.address)
		cpu.C = value & 1
		value = (value >> 1) | (c << 7)
		cpu.Write(info.address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// RTI - Return from Interrupt
func (cpu *CPU) rti(info *stepInfo) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
	cpu.PC = cpu.pull16()
}

// RTS - Return from Subroutine
func (cpu *CPU) rts(info *stepInfo) {
	cpu.PC = cpu.pull16() + 1
}

// SBC - Subtract with Carry
func (cpu *CPU) sbc(info *stepInfo) {
	a := cpu.A
	b := cpu.Read(info.address)
	c := cpu.C
	cpu.A = a - b - (1 - c)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
	if int(a)-int(b)-int(1-c) >= 0 {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^b)&0x80 != 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

// SEC - Set Carry Flag
func (cpu *CPU) sec(info *stepInfo) {
	cpu.C = 1
}

// SED - Set Decimal Flag
func (cpu *CPU) sed(info *stepInfo) {
	cpu.D = 1
}

// SEI - Set Interrupt Disable
func (cpu *CPU) sei(info *stepInfo) {
	cpu.I = 1
}

// STA - Store Accumulator
func (cpu *CPU) sta(info *stepInfo) {
	cpu.Write(info.address, cpu.A)
}

// STX - Store X Register
func (cpu *CPU) stx(info *stepInfo) {
	cpu.Write(info.address, cpu.X)
}

// STY - Store Y Register
func (cpu *CPU) sty(info *stepInfo) {
	cpu.Write(info.address, cpu.Y)
}

// TAX - Transfer Accumulator to X
func (cpu *CPU) tax(info *stepInfo) {
	cpu.X = cpu.A
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// TAY - Transfer Accumulator to Y
func (cpu *CPU) tay(info *stepInfo) {
	cpu.Y = cpu.A
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// TSX - Transfer Stack Pointer to X
func (cpu *CPU) tsx(info *stepInfo) {
	cpu.X = cpu.SP
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// TXA - Transfer X to Accumulator
func (cpu *CPU) txa(info *stepInfo) {
	cpu.A = cpu.X
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// TXS - Transfer X to Stack Pointer
func (cpu *CPU) txs(info *stepInfo) {
	cpu.SP = cpu.X
}

// TYA - Transfer Y to Accumulator
func (cpu *CPU) tya(info *stepInfo) {
	cpu.A = cpu.Y
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// UNK - Unknown Opcode
func (cpu *CPU) unk(info *stepInfo) {
}
