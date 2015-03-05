package nes

import (
	"fmt"
	"log"
	"reflect"
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

var InstructionModes = [256]byte{
	0x6, 0x7, 0x0, 0x0, 0x0, 0xB, 0xB, 0x0, 0x6, 0x5, 0x4, 0x0, 0x0, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x1, 0x7, 0x0, 0x0, 0xB, 0xB, 0xB, 0x0, 0x6, 0x5, 0x4, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x6, 0x7, 0x0, 0x0, 0x0, 0xB, 0xB, 0x0, 0x6, 0x5, 0x4, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x6, 0x7, 0x0, 0x0, 0x0, 0xB, 0xB, 0x0, 0x6, 0x5, 0x4, 0x0, 0x8, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x0, 0x7, 0x0, 0x0, 0xB, 0xB, 0xB, 0x0, 0x6, 0x0, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0xC, 0xC, 0xD, 0x0, 0x6, 0x3, 0x6, 0x0, 0x0, 0x2, 0x0, 0x0,
	0x5, 0x7, 0x5, 0x0, 0xB, 0xB, 0xB, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0xC, 0xC, 0xD, 0x0, 0x6, 0x3, 0x6, 0x0, 0x2, 0x2, 0x3, 0x0,
	0x5, 0x7, 0x0, 0x0, 0xB, 0xB, 0xB, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
	0x5, 0x7, 0x0, 0x0, 0xB, 0xB, 0xB, 0x0, 0x6, 0x5, 0x6, 0x0, 0x1, 0x1, 0x1, 0x0,
	0xA, 0x9, 0x0, 0x0, 0x0, 0xC, 0xC, 0x0, 0x6, 0x3, 0x0, 0x0, 0x0, 0x2, 0x2, 0x0,
}

var InstructionSizes = [256]byte{
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

var InstructionNames = [256]string{
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

func PagesDiffer(a, b uint16) bool {
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
	U      byte   // unused flag
	V      byte   // overflow flag
	N      byte   // negative flag
}

func NewCPU(memory Memory) *CPU {
	cpu := CPU{Memory: memory}
	for i := 0; i < 256; i++ {
		x := reflect.ValueOf(&cpu).MethodByName(InstructionNames[i]).Interface()
		if method, ok := x.(func(uint16)); ok {
			cpu.Table[i] = method
		}
	}
	cpu.Reset()
	return &cpu
}

func (cpu *CPU) Reset() {
	cpu.Cycles = 0
	// cpu.PC = cpu.Read16(0xFFFC)
	cpu.PC = 0xC000
	cpu.SP = 0xFD
	cpu.SetFlags(0x24)
}

// Stack Functions

func (cpu *CPU) Push(value byte) {
	cpu.Write(0x100|uint16(cpu.SP), value)
	cpu.SP--
}

func (cpu *CPU) Pull() byte {
	cpu.SP++
	return cpu.Read(0x100 | uint16(cpu.SP))
}

func (cpu *CPU) Push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	cpu.Push(hi)
	cpu.Push(lo)
}

func (cpu *CPU) Pull16() uint16 {
	lo := uint16(cpu.Pull())
	hi := uint16(cpu.Pull())
	return hi<<8 | lo
}

// Flag Functions

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
	bytes := InstructionSizes[opcode]
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
		"%4X  %s %s %s  %s %28s"+
			"A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%3d\n",
		cpu.PC, w0, w1, w2, name, "",
		cpu.A, cpu.X, cpu.Y, cpu.Flags(), cpu.SP, (cpu.Cycles*3)%341)
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
		// pageCrossed = PagesDiffer(address-uint16(cpu.X), address)
	case ABSOLUTE_Y:
		address = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
		// pageCrossed = PagesDiffer(address-uint16(cpu.Y), address)
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
		// pageCrossed = PagesDiffer(address-uint16(cpu.Y), address)
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

	cpu.PC += uint16(InstructionSizes[opcode])
	cpu.Cycles += uint64(InstructionCycles[opcode])

	cpu.Table[opcode](address)
}

// Instructions
func (cpu *CPU) ADC(address uint16) {
	log.Fatalln("Unimplemented instruction: ADC")
}

func (cpu *CPU) AND(address uint16) {
	cpu.A = cpu.A & cpu.Read(address)
	cpu.SetZ(cpu.A)
	cpu.SetN(cpu.A)
}

func (cpu *CPU) ASL(address uint16) {
	log.Fatalln("Unimplemented instruction: ASL")
}

func (cpu *CPU) BCC(address uint16) {
	if cpu.C == 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BCS(address uint16) {
	if cpu.C != 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BEQ(address uint16) {
	if cpu.Z != 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BIT(address uint16) {
	value := cpu.Read(address)
	cpu.V = (value >> 6) & 1
	cpu.SetZ(value & cpu.A)
	cpu.SetN(value)
}

func (cpu *CPU) BMI(address uint16) {
	if cpu.N != 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BNE(address uint16) {
	if cpu.Z == 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BPL(address uint16) {
	if cpu.N == 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BRK(address uint16) {
	log.Fatalln("Unimplemented instruction: BRK")
}

func (cpu *CPU) BVC(address uint16) {
	if cpu.V == 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) BVS(address uint16) {
	if cpu.V != 0 {
		cpu.PC = address
		cpu.Cycles++
	}
}

func (cpu *CPU) CLC(address uint16) {
	cpu.C = 0
}

func (cpu *CPU) CLD(address uint16) {
	cpu.D = 0
}

func (cpu *CPU) CLI(address uint16) {
	cpu.I = 0
}

func (cpu *CPU) CLV(address uint16) {
	cpu.V = 0
}

func (cpu *CPU) CMP(address uint16) {
	value := cpu.A - cpu.Read(address)
	if value >= 0 {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	cpu.SetZ(value)
	cpu.SetN(value)
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
	cpu.A = cpu.A ^ cpu.Read(address)
	cpu.SetZ(cpu.A)
	cpu.SetN(cpu.A)
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
	cpu.Push16(cpu.PC - 1)
	cpu.PC = address
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
}

func (cpu *CPU) ORA(address uint16) {
	cpu.A = cpu.A | cpu.Read(address)
	cpu.SetZ(cpu.A)
	cpu.SetN(cpu.A)
}

func (cpu *CPU) PHA(address uint16) {
	cpu.Push(cpu.A)
}

func (cpu *CPU) PHP(address uint16) {
	cpu.Push(cpu.Flags() | 0x10)
}

func (cpu *CPU) PLA(address uint16) {
	cpu.A = cpu.Pull()
	cpu.SetZ(cpu.A)
	cpu.SetN(cpu.A)
}

func (cpu *CPU) PLP(address uint16) {
	cpu.SetFlags(cpu.Pull()&0xEF | 0x20)
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
	cpu.PC = cpu.Pull16() + 1
}

func (cpu *CPU) SBC(address uint16) {
	log.Fatalln("Unimplemented instruction: SBC")
}

func (cpu *CPU) SEC(address uint16) {
	cpu.C = 1
}

func (cpu *CPU) SED(address uint16) {
	cpu.D = 1
}

func (cpu *CPU) SEI(address uint16) {
	cpu.I = 1
}

func (cpu *CPU) STA(address uint16) {
	cpu.Write(address, cpu.A)
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
