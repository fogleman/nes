package nes

type CPU struct {
	PC                  uint16 // program counter
	SP                  byte   // stack pointer
	A, X, Y             byte   // registers
	C, Z, I, D, B, V, N bool   // flags (TODO: bytes? combine?)
}
