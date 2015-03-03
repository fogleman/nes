package nes

type CPU struct {
	PC                  uint16 // program counter
	SP                  byte   // stack pointer
	A, X, Y             byte   // registers
	C, Z, I, D, B, V, N byte   // flags
}

func (cpu *CPU) GetFlags() byte {
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
